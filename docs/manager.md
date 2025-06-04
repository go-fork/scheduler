# Manager API

`Manager` là interface chính cho việc lập lịch và quản lý jobs. Nó cung cấp một API fluent cho việc tạo, lập lịch và quản lý các công việc định kỳ.

## Tạo Manager

```go
// Cách 1: Tạo manager mới với cấu hình mặc định
manager := scheduler.NewScheduler()

// Cách 2: Tạo manager với cấu hình tùy chỉnh
config := scheduler.DefaultConfig()
config.AutoStart = false
manager := scheduler.NewSchedulerWithConfig(config)
```

## API Lập lịch

### Khoảng thời gian (Interval)

```go
// Chạy mỗi 5 phút
manager.Every(5).Minutes().Do(func() {
    fmt.Println("Task chạy mỗi 5 phút")
})

// Chạy mỗi giờ
manager.Every(1).Hour().Do(func() {
    fmt.Println("Task chạy mỗi giờ")
})

// Chạy mỗi 30 giây
manager.Every(30).Seconds().Do(func() {
    fmt.Println("Task chạy mỗi 30 giây")
})

// Chạy mỗi ngày
manager.Every(1).Day().Do(func() {
    fmt.Println("Task chạy mỗi ngày")
})

// Chạy mỗi tuần
manager.Every(1).Week().Do(func() {
    fmt.Println("Task chạy mỗi tuần")
})
```

### Thời điểm cụ thể

```go
// Chạy lúc 10:30 hàng ngày
manager.At("10:30").Every(1).Day().Do(func() {
    fmt.Println("Task chạy lúc 10:30 mỗi ngày")
})

// Chạy nhiều thời điểm trong ngày
manager.At("09:00", "12:00", "18:00").Every(1).Day().Do(func() {
    fmt.Println("Task chạy 3 lần mỗi ngày")
})

// Định dạng giờ: "HH:MM" hoặc "HH:MM:SS"
manager.At("15:30:45").Every(1).Day().Do(func() {
    fmt.Println("Task chạy lúc 3:30:45 PM mỗi ngày")
})
```

### Cron Expressions

```go
// Cron expression (mỗi 6 giờ)
manager.Cron("0 */6 * * *").Do(func() {
    fmt.Println("Task chạy mỗi 6 giờ")
})

// Cron phức tạp (Thứ 2 đến Thứ 6, 9:00 AM)
manager.Cron("0 9 * * 1-5").Do(func() {
    fmt.Println("Task chạy lúc 9:00 AM từ T2-T6")
})

// Cron syntax có thể sử dụng các predefined schedules
manager.Cron("@daily").Do(func() {
    fmt.Println("Task chạy mỗi ngày lúc nửa đêm")
})
```

### Thời điểm bắt đầu

```go
// Bắt đầu job từ thời điểm cụ thể
startTime := time.Date(2025, 6, 15, 12, 0, 0, 0, time.Local)
manager.Every(1).Hour().StartAt(startTime).Do(func() {
    fmt.Println("Task chạy mỗi giờ, bắt đầu từ 15/06/2025 12:00")
})
```

## Quản lý Job

### Tagging

```go
// Thêm tags cho jobs
manager.Every(5).Minutes().
    Tag("backup").
    Tag("database").
    Do(func() {
        fmt.Println("Database backup")
    })

// Xóa jobs theo tag
err := manager.RemoveByTag("backup")
if err != nil {
    log.Error("Lỗi khi xóa jobs: %v", err)
}
```

### Singleton Mode

```go
// Đảm bảo chỉ một instance của job chạy cùng lúc
manager.Every(1).Hour().
    SingletonMode().
    Do(func() {
        // Task quan trọng
        fmt.Println("Critical task - only one instance will run")
    })
```

### Liệt kê Jobs

```go
// Liệt kê tất cả jobs
jobs := manager.Jobs()
for _, job := range jobs {
    fmt.Printf("Job: %s, Next run: %v\n", job.Tags(), job.NextRun())
}
```

### Xóa Jobs

```go
// Xóa tất cả jobs
manager.Clear()
```

## Quản lý Lifecycle

### Khởi động và Dừng

```go
// Khởi động scheduler
err := manager.Start()
if err != nil {
    log.Fatal("Không thể khởi động scheduler: %v", err)
}

// Kiểm tra trạng thái
if manager.IsRunning() {
    fmt.Println("Scheduler đang chạy")
}

// Dừng scheduler
manager.Stop()
```

### Đặt tên cho Scheduler

```go
// Đặt tên cho scheduler (hữu ích cho logging và monitoring)
manager := scheduler.NewScheduler().Name("ApplicationScheduler")
```

## Event Listeners

```go
// Định nghĩa event listener
type LoggingEventListener struct {
    logger Logger
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

// Đăng ký event listener
listener := &LoggingEventListener{logger: appLogger}
manager.RegisterEventListeners(listener)
```

## Truy cập Underlying Scheduler

```go
// Truy cập gocron scheduler bên dưới
gochronScheduler := manager.GetScheduler()

// Sử dụng các tính năng nâng cao (nếu cần)
gochronScheduler.SetMaxConcurrentJobs(10, gocron.WaitMode)
```
