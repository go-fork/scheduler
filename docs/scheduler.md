# Scheduler Package - Tài liệu Kỹ thuật

## Tổng quan

Package `go.fork.vn/scheduler` cung cấp một hệ thống lập lịch linh hoạt và mạnh mẽ cho ứng dụng Go, với khả năng tích hợp dependency injection, distributed locking, và hỗ trợ các pattern cron phức tạp.

## Kiến trúc

### Core Components

#### 1. Manager Interface
```go
type Manager interface {
    // Job scheduling
    Every(interval uint) *gocron.Job
    At(times ...string) *gocron.Job
    Cron(cronExpression string) (*gocron.Job, error)
    
    // Lifecycle management  
    Start() error
    Stop() error
    IsRunning() bool
    Clear() error
    
    // Job management
    Jobs() []*gocron.Job
    RemoveByTag(tag string) error
    
    // Configuration
    Name(name string) Manager
    SingletonMode() Manager
    WithDistributedLocker(locker gocron.Locker) Manager
    RegisterEventListeners(eventListeners gocron.EventListener) Manager
}
```

#### 2. Distributed Locking System
Package cung cấp Redis-based distributed locking để đảm bảo chỉ một instance của scheduled job chạy tại một thời điểm:

```go
type RedisLocker struct {
    client  redis.RedisClient
    options *RedisLockerOptions
}

type RedisLockerOptions struct {
    LockDuration time.Duration
    MaxRetries   int
    RetryDelay   time.Duration
    KeyPrefix    string
}
```

#### 3. ServiceProvider
Tích hợp với DI container:

```go
type ServiceProvider struct{}

func (p *ServiceProvider) Register(app Application) error
func (p *ServiceProvider) Boot(app Application) error
func (p *ServiceProvider) Requires() []string
func (p *ServiceProvider) Providers() []string
```

## Features

### 1. Flexible Scheduling
- **Interval-based**: `Every(5).Minutes()`
- **Time-based**: `At("10:30").Every(1).Day()`
- **Cron expressions**: `Cron("0 */6 * * *")`
- **Complex patterns**: Monthly, weekly, daily schedules

### 2. Distributed Execution
- **Redis-based locking**: Prevents duplicate job execution across multiple instances
- **Configurable timeouts**: Customizable lock duration and retry mechanisms
- **Automatic cleanup**: Expired locks are automatically released

### 3. Job Management
- **Tagging system**: Group and manage related jobs
- **Singleton mode**: Ensure only one instance of a job runs
- **Event listeners**: Hook into job lifecycle events
- **Dynamic scheduling**: Add, remove, and modify jobs at runtime

### 4. Integration Capabilities
- **Config integration**: Uses config package for Redis connection settings
- **DI container**: Seamless integration with dependency injection
- **Redis provider**: Leverages redis package for underlying storage

## Advanced Features

### 1. Distributed Locking Implementation

#### Redis Lock Mechanism
```go
type RedisLocker struct {
    client  redis.RedisClient
    options *RedisLockerOptions
}

func (l *RedisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
    lockKey := l.options.KeyPrefix + key
    lockValue := generateLockValue()
    
    for i := 0; i <= l.options.MaxRetries; i++ {
        success, err := l.client.SetNX(ctx, lockKey, lockValue, l.options.LockDuration).Result()
        if err != nil {
            return nil, err
        }
        
        if success {
            return &RedisLock{
                client:   l.client,
                key:      lockKey,
                value:    lockValue,
                duration: l.options.LockDuration,
            }, nil
        }
        
        if i < l.options.MaxRetries {
            time.Sleep(l.options.RetryDelay)
        }
    }
    
    return nil, errors.New("failed to acquire lock after retries")
}
```

#### Lock Release Mechanism
```go
type RedisLock struct {
    client   redis.RedisClient
    key      string
    value    string
    duration time.Duration
}

func (l *RedisLock) Unlock(ctx context.Context) error {
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    
    result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
    if err != nil {
        return err
    }
    
    if result.(int64) == 0 {
        return errors.New("lock not owned by this instance")
    }
    
    return nil
}
```

### 2. Configuration Management

#### Default Redis Locker Options
```go
func DefaultRedisLockerOptions() *RedisLockerOptions {
    return &RedisLockerOptions{
        LockDuration: 30 * time.Second,
        MaxRetries:   3,
        RetryDelay:   100 * time.Millisecond,
        KeyPrefix:    "scheduler:lock:",
    }
}
```

#### Validation
```go
func ValidateRedisLockerOptions(options *RedisLockerOptions) error {
    if options.LockDuration <= 0 {
        return errors.New("lock duration must be positive")
    }
    
    if options.MaxRetries < 0 {
        return errors.New("max retries cannot be negative")
    }
    
    if options.RetryDelay < 0 {
        return errors.New("retry delay cannot be negative")
    }
    
    if options.KeyPrefix == "" {
        return errors.New("key prefix cannot be empty")
    }
    
    return nil
}
```

### 3. Event System

#### Event Listeners
```go
type EventListener interface {
    OnJobRun(jobID string, jobName string)
    OnJobComplete(jobID string, jobName string, duration time.Duration)
    OnJobError(jobID string, jobName string, err error)
}

// Custom event listener implementation
type LoggingEventListener struct {
    logger log.Manager
}

func (l *LoggingEventListener) OnJobRun(jobID, jobName string) {
    l.logger.Info("Job started: id=%s, name=%s", jobID, jobName)
}

func (l *LoggingEventListener) OnJobComplete(jobID, jobName string, duration time.Duration) {
    l.logger.Info("Job completed: id=%s, name=%s, duration=%v", jobID, jobName, duration)
}

func (l *LoggingEventListener) OnJobError(jobID, jobName string, err error) {
    l.logger.Error("Job failed: id=%s, name=%s, error=%v", jobID, jobName, err)
}
```

## Performance Considerations

### 1. Lock Efficiency
- **Lua scripts**: Atomic operations reduce race conditions
- **Configurable timeouts**: Balance between safety and performance
- **Retry mechanisms**: Handle temporary Redis unavailability

### 2. Memory Management
- **Job cleanup**: Automatic removal of completed jobs
- **Connection pooling**: Efficient Redis connection usage
- **Resource disposal**: Proper cleanup on scheduler shutdown

### 3. Scalability
- **Horizontal scaling**: Multiple instances can safely run with distributed locking
- **Load distribution**: Jobs spread across available instances
- **Fault tolerance**: Continue operation if some instances fail

## Integration Patterns

### 1. With Config Package
```go
// Configuration structure
type SchedulerConfig struct {
    Enabled         bool          `mapstructure:"enabled"`
    MaxConcurrent   int           `mapstructure:"max_concurrent"`
    LockDuration    time.Duration `mapstructure:"lock_duration"`
    RetryDelay      time.Duration `mapstructure:"retry_delay"`
    MaxRetries      int           `mapstructure:"max_retries"`
}

// Using config in scheduler setup
func setupScheduler(configManager config.Manager) scheduler.Manager {
    var cfg SchedulerConfig
    configManager.UnmarshalKey("scheduler", &cfg)
    
    if !cfg.Enabled {
        return nil
    }
    
    manager := scheduler.NewManager()
    if cfg.MaxConcurrent > 0 {
        manager.SetMaxRunnableJobs(cfg.MaxConcurrent)
    }
    
    return manager
}
```

### 2. With Redis Package
```go
// Redis locker setup
func setupDistributedScheduler(
    redisManager redis.Manager, 
    schedulerManager scheduler.Manager,
) scheduler.Manager {
    
    client := redisManager.GetConnection("default")
    
    options := &scheduler.RedisLockerOptions{
        LockDuration: 60 * time.Second,
        MaxRetries:   5,
        RetryDelay:   200 * time.Millisecond,
        KeyPrefix:    "myapp:scheduler:",
    }
    
    locker := scheduler.NewRedisLocker(client, options)
    return schedulerManager.WithDistributedLocker(locker)
}
```

### 3. With Logging Package
```go
// Event listener with logging
func setupSchedulerWithLogging(
    schedulerManager scheduler.Manager,
    logManager log.Manager,
) scheduler.Manager {
    
    eventListener := &LoggingEventListener{logger: logManager}
    return schedulerManager.RegisterEventListeners(eventListener)
}
```

## Testing Support

### 1. Mock Manager
```go
type MockManager struct {
    mock.Mock
}

func (m *MockManager) Every(interval uint) *gocron.Job {
    args := m.Called(interval)
    return args.Get(0).(*gocron.Job)
}

func (m *MockManager) Start() error {
    args := m.Called()
    return args.Error(0)
}

// Usage in tests
func TestJobScheduling(t *testing.T) {
    mockScheduler := &MockManager{}
    mockJob := &gocron.Job{}
    
    mockScheduler.On("Every", uint(5)).Return(mockJob)
    mockScheduler.On("Start").Return(nil)
    
    // Test your code
    result := schedulePeriodicJob(mockScheduler)
    
    assert.NoError(t, result)
    mockScheduler.AssertExpectations(t)
}
```

### 2. Test Utilities
```go
// Helper for testing scheduled jobs
func TestScheduledJobExecution(t *testing.T) {
    manager := scheduler.NewManager()
    executed := false
    
    job, err := manager.Cron("* * * * * *").Do(func() {
        executed = true
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, job)
    
    manager.Start()
    defer manager.Stop()
    
    // Wait for job execution
    time.Sleep(2 * time.Second)
    assert.True(t, executed)
}
```

## Best Practices

### 1. Resource Management
```go
// Always clean up scheduler resources
defer manager.Stop()
defer manager.Clear()
```

### 2. Error Handling
```go
// Proper error handling in scheduled jobs
job, err := manager.Every(5).Minutes().Do(func() {
    if err := performTask(); err != nil {
        log.Error("Scheduled task failed: %v", err)
        // Consider implementing retry logic or alerting
    }
})

if err != nil {
    log.Error("Failed to schedule job: %v", err)
}
```

### 3. Job Naming and Tagging
```go
// Use descriptive names and tags for better management
manager.Every(10).Minutes().
    Tag("maintenance").
    Tag("database").
    Name("cleanup-old-records").
    Do(cleanupOldRecords)

manager.Every(1).Hour().
    Tag("backup").
    Name("backup-user-data").
    Do(backupUserData)
```

### 4. Distributed Environment Considerations
```go
// Always use distributed locking in multi-instance deployments
if isDistributedEnvironment() {
    locker := scheduler.NewRedisLocker(redisClient, options)
    manager = manager.WithDistributedLocker(locker)
}

// Use singleton mode for critical jobs
manager.Every(1).Day().
    Tag("critical").
    SingletonMode().
    Do(criticalDailyTask)
```

### 5. Configuration Management
```go
// Environment-specific configuration
func getSchedulerConfig() *SchedulerConfig {
    env := os.Getenv("ENVIRONMENT")
    
    switch env {
    case "production":
        return &SchedulerConfig{
            LockDuration: 300 * time.Second, // 5 minutes
            MaxRetries:   10,
            RetryDelay:   1 * time.Second,
        }
    case "development":
        return &SchedulerConfig{
            LockDuration: 30 * time.Second,
            MaxRetries:   3,
            RetryDelay:   100 * time.Millisecond,
        }
    default:
        return DefaultSchedulerConfig()
    }
}
```
