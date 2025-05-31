# Scheduler Package - Hướng dẫn Sử dụng

## Cài đặt

```bash
go get go.fork.vn/scheduler@v0.1.0
```

## Sử dụng Cơ bản

### 1. Import Package
```go
import (
    "go.fork.vn/scheduler"
    "go.fork.vn/redis"
    "go.fork.vn/config"
)
```

### 2. Tạo Manager
```go
// Tạo scheduler manager
manager := scheduler.NewManager()

// Đặt tên cho scheduler (optional)
manager = manager.Name("MyAppScheduler")
```

### 3. Lập lịch Jobs Cơ bản

#### Interval-based Jobs
```go
// Chạy mỗi 5 phút
job := manager.Every(5).Minutes().Do(func() {
    fmt.Println("Task chạy mỗi 5 phút")
})

// Chạy mỗi 30 giây
job := manager.Every(30).Seconds().Do(func() {
    fmt.Println("Task chạy mỗi 30 giây")
})

// Chạy mỗi giờ
job := manager.Every(1).Hour().Do(func() {
    fmt.Println("Task chạy mỗi giờ")
})
```

#### Time-specific Jobs
```go
// Chạy lúc 10:30 hàng ngày
job := manager.At("10:30").Every(1).Day().Do(func() {
    fmt.Println("Task chạy lúc 10:30 mỗi ngày")
})

// Chạy nhiều thời điểm trong ngày
job := manager.At("09:00", "12:00", "18:00").Every(1).Day().Do(func() {
    fmt.Println("Task chạy 3 lần mỗi ngày")
})
```

#### Cron Expressions
```go
// Cron expression (mỗi 6 giờ)
job, err := manager.Cron("0 */6 * * *").Do(func() {
    fmt.Println("Task chạy mỗi 6 giờ")
})
if err != nil {
    log.Fatal("Lỗi cron expression:", err)
}

// Cron phức tạp (Thứ 2 đến Thứ 6, 9:00 AM)
job, err := manager.Cron("0 9 * * 1-5").Do(func() {
    fmt.Println("Task chạy lúc 9:00 AM từ T2-T6")
})
```

### 4. Khởi động và Dừng Scheduler
```go
// Khởi động scheduler
err := manager.Start()
if err != nil {
    log.Fatal("Không thể khởi động scheduler:", err)
}

// Dừng scheduler (thường trong defer hoặc signal handler)
defer manager.Stop()

// Hoặc dừng tất cả jobs và xóa lịch
defer manager.Clear()
```

## Tính năng Nâng cao

### 1. Job Tagging và Management
```go
// Thêm tags cho jobs
manager.Every(5).Minutes().
    Tag("backup").
    Tag("database").
    Do(func() {
        performDatabaseBackup()
    })

manager.Every(10).Minutes().
    Tag("cleanup").
    Tag("temporary").
    Do(func() {
        cleanupTempFiles()
    })

// Xóa jobs theo tag
err := manager.RemoveByTag("temporary")
if err != nil {
    log.Error("Lỗi khi xóa jobs:", err)
}

// Liệt kê tất cả jobs
jobs := manager.Jobs()
for _, job := range jobs {
    fmt.Printf("Job: %s, Next run: %s\n", job.Tags(), job.NextRun())
}
```

### 2. Singleton Mode
```go
// Đảm bảo chỉ một instance của job chạy
manager.Every(1).Hour().
    SingletonMode().
    Do(func() {
        // Task quan trọng chỉ nên chạy một lần
        performCriticalMaintenance()
    })
```

### 3. Distributed Locking với Redis

#### Setup Redis Locker
```go
import (
    "go.fork.vn/redis"
    "go.fork.vn/scheduler"
)

func setupDistributedScheduler() scheduler.Manager {
    // Khởi tạo Redis client
    redisManager := redis.NewManager()
    client := redisManager.GetConnection("default")
    
    // Cấu hình locker options
    options := &scheduler.RedisLockerOptions{
        LockDuration: 60 * time.Second,   // Lock tồn tại 60 giây
        MaxRetries:   5,                  // Thử lại tối đa 5 lần
        RetryDelay:   500 * time.Millisecond, // Delay giữa các lần thử
        KeyPrefix:    "myapp:scheduler:", // Prefix cho Redis keys
    }
    
    // Tạo distributed locker
    locker := scheduler.NewRedisLocker(client, options)
    
    // Áp dụng locker cho scheduler
    manager := scheduler.NewManager().
        Name("DistributedScheduler").
        WithDistributedLocker(locker)
    
    return manager
}
```

#### Sử dụng Distributed Scheduler
```go
func main() {
    manager := setupDistributedScheduler()
    
    // Jobs sẽ chỉ chạy trên một instance
    manager.Every(5).Minutes().
        Tag("distributed").
        Do(func() {
            fmt.Println("Task này chỉ chạy trên một server")
            processGlobalTask()
        })
    
    manager.Start()
    defer manager.Stop()
    
    // Keep application running
    select {}
}
```

## Tích hợp với Dependency Injection

### 1. Sử dụng ServiceProvider
```go
import (
    "go.fork.vn/scheduler"
    "go.fork.vn/di"
)

func main() {
    container := di.NewContainer()
    
    // Đăng ký ServiceProvider
    provider := &scheduler.ServiceProvider{}
    provider.Register(container)
    provider.Boot(container)
    
    // Sử dụng từ container
    container.Call(func(manager scheduler.Manager) {
        manager.Every(10).Seconds().Do(func() {
            fmt.Println("Task từ DI container")
        })
        
        manager.Start()
        defer manager.Stop()
    })
}
```

### 2. Custom Configuration với DI
```go
type App struct {
    container *di.Container
    basePath  string
}

func (app *App) Container() *di.Container {
    return app.container
}

func (app *App) BasePath() string {
    return app.basePath
}

func main() {
    app := &App{
        container: di.NewContainer(),
        basePath:  "/app",
    }
    
    // Đăng ký các providers cần thiết
    configProvider := &config.ServiceProvider{}
    redisProvider := &redis.ServiceProvider{}
    schedulerProvider := &scheduler.ServiceProvider{}
    
    configProvider.Register(app)
    redisProvider.Register(app)
    schedulerProvider.Register(app)
    
    configProvider.Boot(app)
    redisProvider.Boot(app)
    schedulerProvider.Boot(app)
}
```

## Patterns và Best Practices

### 1. Job Functions với Parameters
```go
// Function với parameters
func processUserData(userID int, action string) {
    fmt.Printf("Processing user %d with action %s\n", userID, action)
}

// Schedule với anonymous function
manager.Every(5).Minutes().Do(func() {
    processUserData(123, "cleanup")
})

// Hoặc sử dụng closure
func scheduleUserProcessing(userID int, action string) {
    manager.Every(10).Minutes().Do(func() {
        processUserData(userID, action)
    })
}
```

### 2. Error Handling trong Jobs
```go
// Pattern 1: Error handling trong job function
manager.Every(5).Minutes().Do(func() {
    if err := performTask(); err != nil {
        log.Error("Task failed: %v", err)
        // Có thể gửi alert hoặc notification
        sendErrorAlert(err)
    }
})

// Pattern 2: Retry mechanism
func performTaskWithRetry(maxRetries int) {
    for i := 0; i < maxRetries; i++ {
        if err := performTask(); err != nil {
            log.Warning("Task failed (attempt %d/%d): %v", i+1, maxRetries, err)
            if i == maxRetries-1 {
                log.Error("Task failed after %d attempts", maxRetries)
                return
            }
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        log.Info("Task completed successfully")
        return
    }
}

manager.Every(10).Minutes().Do(func() {
    performTaskWithRetry(3)
})
```

### 3. Context-aware Jobs
```go
import (
    "context"
    "time"
)

// Job với timeout
manager.Every(30).Minutes().Do(func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    
    if err := performLongRunningTask(ctx); err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            log.Warning("Task timed out after 5 minutes")
        } else {
            log.Error("Task failed: %v", err)
        }
    }
})

func performLongRunningTask(ctx context.Context) error {
    // Simulate long running task with context checking
    for i := 0; i < 100; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Do work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

### 4. Dynamic Job Management
```go
type JobManager struct {
    scheduler scheduler.Manager
    jobs      map[string]*gocron.Job
    mutex     sync.RWMutex
}

func NewJobManager(scheduler scheduler.Manager) *JobManager {
    return &JobManager{
        scheduler: scheduler,
        jobs:      make(map[string]*gocron.Job),
    }
}

func (jm *JobManager) AddPeriodicJob(name string, intervalMinutes int, task func()) error {
    jm.mutex.Lock()
    defer jm.mutex.Unlock()
    
    // Remove existing job if exists
    if existingJob, exists := jm.jobs[name]; exists {
        jm.scheduler.RemoveByTag(name)
    }
    
    // Add new job
    job := jm.scheduler.Every(uint(intervalMinutes)).Minutes().
        Tag(name).
        Do(task)
    
    jm.jobs[name] = job
    return nil
}

func (jm *JobManager) RemoveJob(name string) error {
    jm.mutex.Lock()
    defer jm.mutex.Unlock()
    
    if _, exists := jm.jobs[name]; !exists {
        return fmt.Errorf("job %s not found", name)
    }
    
    err := jm.scheduler.RemoveByTag(name)
    if err == nil {
        delete(jm.jobs, name)
    }
    
    return err
}

// Sử dụng
func main() {
    manager := scheduler.NewManager()
    jobManager := NewJobManager(manager)
    
    // Add jobs dynamically
    jobManager.AddPeriodicJob("cleanup", 30, func() {
        fmt.Println("Cleanup task")
    })
    
    jobManager.AddPeriodicJob("backup", 60, func() {
        fmt.Println("Backup task")
    })
    
    manager.Start()
    defer manager.Stop()
}
```

### 5. Event Listeners và Monitoring
```go
import (
    "go.fork.vn/log"
)

type SchedulerEventListener struct {
    logger log.Manager
}

func (l *SchedulerEventListener) OnJobRun(jobID, jobName string) {
    l.logger.Info("Job started: %s (%s)", jobName, jobID)
}

func (l *SchedulerEventListener) OnJobComplete(jobID, jobName string, duration time.Duration) {
    l.logger.Info("Job completed: %s (%s) in %v", jobName, jobID, duration)
}

func (l *SchedulerEventListener) OnJobError(jobID, jobName string, err error) {
    l.logger.Error("Job failed: %s (%s) - %v", jobName, jobID, err)
}

// Setup với event listener
func setupSchedulerWithMonitoring(logManager log.Manager) scheduler.Manager {
    eventListener := &SchedulerEventListener{logger: logManager}
    
    return scheduler.NewManager().
        Name("MonitoredScheduler").
        RegisterEventListeners(eventListener)
}
```

## Configuration Examples

### 1. Environment-based Configuration
```yaml
# config/app.yaml
scheduler:
  enabled: true
  max_concurrent: 10
  distributed:
    enabled: true
    redis:
      connection: "default"
    lock:
      duration: "60s"
      max_retries: 5
      retry_delay: "500ms"
      key_prefix: "myapp:scheduler:"
```

```go
// Sử dụng config
func setupSchedulerFromConfig(
    configManager config.Manager,
    redisManager redis.Manager,
) scheduler.Manager {
    
    type SchedulerConfig struct {
        Enabled        bool `mapstructure:"enabled"`
        MaxConcurrent  int  `mapstructure:"max_concurrent"`
        Distributed    struct {
            Enabled bool `mapstructure:"enabled"`
            Redis   struct {
                Connection string `mapstructure:"connection"`
            } `mapstructure:"redis"`
            Lock struct {
                Duration   time.Duration `mapstructure:"duration"`
                MaxRetries int           `mapstructure:"max_retries"`
                RetryDelay time.Duration `mapstructure:"retry_delay"`
                KeyPrefix  string        `mapstructure:"key_prefix"`
            } `mapstructure:"lock"`
        } `mapstructure:"distributed"`
    }
    
    var cfg SchedulerConfig
    configManager.UnmarshalKey("scheduler", &cfg)
    
    if !cfg.Enabled {
        return nil
    }
    
    manager := scheduler.NewManager()
    
    if cfg.Distributed.Enabled {
        client := redisManager.GetConnection(cfg.Distributed.Redis.Connection)
        
        options := &scheduler.RedisLockerOptions{
            LockDuration: cfg.Distributed.Lock.Duration,
            MaxRetries:   cfg.Distributed.Lock.MaxRetries,
            RetryDelay:   cfg.Distributed.Lock.RetryDelay,
            KeyPrefix:    cfg.Distributed.Lock.KeyPrefix,
        }
        
        locker := scheduler.NewRedisLocker(client, options)
        manager = manager.WithDistributedLocker(locker)
    }
    
    return manager
}
```

### 2. Production Setup Complete
```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "go.fork.vn/config"
    "go.fork.vn/redis"
    "go.fork.vn/scheduler"
    "go.fork.vn/log"
    "go.fork.vn/di"
)

func main() {
    // Setup DI container
    container := di.NewContainer()
    
    // Register providers
    configProvider := &config.ServiceProvider{}
    logProvider := &log.ServiceProvider{}
    redisProvider := &redis.ServiceProvider{}
    schedulerProvider := &scheduler.ServiceProvider{}
    
    // Create app instance
    app := &Application{container: container}
    
    // Register and boot providers
    providers := []interface {
        Register(interface{}) error
        Boot(interface{}) error
    }{
        configProvider,
        logProvider,
        redisProvider,
        schedulerProvider,
    }
    
    for _, provider := range providers {
        provider.Register(app)
        provider.Boot(app)
    }
    
    // Setup jobs
    container.Call(func(
        schedulerManager scheduler.Manager,
        logManager log.Manager,
    ) {
        setupJobs(schedulerManager, logManager)
        
        // Start scheduler
        if err := schedulerManager.Start(); err != nil {
            logManager.Fatal("Failed to start scheduler: %v", err)
        }
        
        logManager.Info("Scheduler started successfully")
        
        // Graceful shutdown
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        
        <-c
        logManager.Info("Shutting down scheduler...")
        schedulerManager.Stop()
        logManager.Info("Scheduler stopped")
    })
}

func setupJobs(schedulerManager scheduler.Manager, logManager log.Manager) {
    // Database backup job
    schedulerManager.At("02:00").Every(1).Day().
        Tag("backup").
        Tag("database").
        SingletonMode().
        Do(func() {
            logManager.Info("Starting database backup")
            // Implement backup logic
            logManager.Info("Database backup completed")
        })
    
    // Log cleanup job
    schedulerManager.Every(6).Hours().
        Tag("cleanup").
        Tag("logs").
        Do(func() {
            logManager.Info("Starting log cleanup")
            // Implement log cleanup logic
            logManager.Info("Log cleanup completed")
        })
    
    // Health check job
    schedulerManager.Every(5).Minutes().
        Tag("monitoring").
        Tag("health").
        Do(func() {
            // Implement health check logic
            logManager.Debug("Health check completed")
        })
}

type Application struct {
    container *di.Container
}

func (app *Application) Container() *di.Container {
    return app.container
}

func (app *Application) BasePath() string {
    return "/app"
}
```

## Troubleshooting

### 1. Jobs không chạy
```go
// Kiểm tra scheduler đã start chưa
if !manager.IsRunning() {
    log.Error("Scheduler is not running")
    manager.Start()
}

// Kiểm tra jobs đã được schedule chưa
jobs := manager.Jobs()
if len(jobs) == 0 {
    log.Warning("No jobs scheduled")
}

// Debug job timing
for _, job := range jobs {
    log.Debug("Job: %v, Next run: %v", job.Tags(), job.NextRun())
}
```

### 2. Distributed locking issues
```go
// Test Redis connection
func testRedisConnection(client redis.RedisClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return client.Ping(ctx).Err()
}

// Verify lock acquisition
options := &scheduler.RedisLockerOptions{
    LockDuration: 10 * time.Second,
    MaxRetries:   1,
    RetryDelay:   100 * time.Millisecond,
    KeyPrefix:    "test:",
}

locker := scheduler.NewRedisLocker(client, options)
lock, err := locker.Lock(context.Background(), "test-key")
if err != nil {
    log.Error("Failed to acquire lock: %v", err)
} else {
    defer lock.Unlock(context.Background())
    log.Info("Lock acquired successfully")
}
```

### 3. Memory leaks
```go
// Proper cleanup
func gracefulShutdown(manager scheduler.Manager) {
    // Stop accepting new jobs
    manager.Stop()
    
    // Clear all scheduled jobs
    manager.Clear()
    
    // Wait for running jobs to complete (implement timeout)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Monitor running jobs
    for {
        select {
        case <-ctx.Done():
            log.Warning("Forced shutdown after timeout")
            return
        default:
            if len(manager.Jobs()) == 0 {
                log.Info("All jobs completed")
                return
            }
            time.Sleep(1 * time.Second)
        }
    }
}
```
