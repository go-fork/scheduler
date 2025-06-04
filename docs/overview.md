# Tổng quan về Scheduler

## Kiến trúc

Scheduler package được xây dựng dựa trên thư viện [go-co-op/gocron](https://github.com/go-co-op/gocron) với các lớp bọc và tích hợp để cung cấp trải nghiệm mượt mà trong hệ sinh thái `go.fork.vn`. Kiến trúc gồm các thành phần chính:

### Core Components

#### 1. Manager Interface
Interface chính định nghĩa tất cả phương thức lập lịch và quản lý jobs:

```go
type Manager interface {
    // Job scheduling
    Every(interval interface{}) Manager
    Second() Manager
    Seconds() Manager
    Minutes() Manager
    Hours() Manager
    Days() Manager
    Weeks() Manager
    At(time string) Manager
    StartAt(time time.Time) Manager
    Cron(cronExpression string) Manager
    Do(job interface{}, params ...interface{}) error
    Tag(tags ...string) Manager
    
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
    RegisterEventListeners(eventListeners ...gocron.EventListener)
    GetScheduler() *gocron.Scheduler
}
```

#### 2. Distributed Locking System
Hệ thống khóa phân tán dựa trên Redis để đảm bảo chỉ một instance của job chạy trong môi trường phân tán:

```go
type RedisLocker struct {
    client  redis.RedisClient
    options RedisLockerOptionsTime
}
```

#### 3. ServiceProvider
Tích hợp với DI container của ứng dụng:

```go
type ServiceProvider struct{}

func (p *ServiceProvider) Register(app di.Application)
func (p *ServiceProvider) Boot(app di.Application)
func (p *ServiceProvider) Requires() []string
func (p *ServiceProvider) Providers() []string
```

## Tính năng

### 1. Flexible Scheduling
- **Khoảng thời gian**: `Every(5).Minutes()`
- **Thời điểm cụ thể**: `At("10:30").Every(1).Day()`
- **Cron expressions**: `Cron("0 */6 * * *")`
- **Mẫu phức tạp**: Lịch hàng tháng, hàng tuần, hàng ngày

### 2. Distributed Execution
- **Redis-based locking**: Ngăn chặn việc thực thi trùng lặp trên nhiều instances
- **Timeouts tùy chỉnh**: Thời gian lock và cơ chế retry có thể tùy chỉnh
- **Tự động dọn dẹp**: Khóa hết hạn được tự động giải phóng

### 3. Job Management
- **Hệ thống tag**: Nhóm và quản lý các job liên quan
- **Chế độ singleton**: Đảm bảo chỉ một instance của job chạy
- **Event listeners**: Gắn kết với các sự kiện vòng đời của job
- **Dynamic scheduling**: Thêm, xóa và sửa đổi jobs lúc runtime

### 4. Integration Capabilities
- **Config integration**: Sử dụng config package cho các cài đặt kết nối Redis
- **DI container**: Tích hợp liền mạch với dependency injection
- **Redis provider**: Tận dụng redis package cho lưu trữ cơ bản

## Performance Considerations

- **Efficient lock management**: Sử dụng Lua scripts để thực hiện các thao tác nguyên tử
- **Balanced timeouts**: Cân bằng giữa an toàn và hiệu suất với các tùy chọn cấu hình
- **Retry mechanisms**: Xử lý các tình huống Redis tạm thời không khả dụng
- **Resource management**: Dọn dẹp tài nguyên đúng cách khi scheduler shutdown

## Best Practices

1. **Luôn dọn dẹp tài nguyên**:
   ```go
   defer manager.Stop()
   defer manager.Clear()
   ```

2. **Xử lý lỗi trong scheduled jobs**:
   ```go
   manager.Every(5).Minutes().Do(func() {
       if err := performTask(); err != nil {
           log.Error("Scheduled task failed: %v", err)
       }
   })
   ```

3. **Đặt tên và tag jobs cho quản lý tốt hơn**:
   ```go
   manager.Every(10).Minutes().
       Tag("maintenance").
       Tag("database").
       Do(cleanupOldRecords)
   ```

4. **Sử dụng distributed locking trong môi trường phân tán**:
   ```go
   locker := scheduler.NewRedisLocker(redisClient, options)
   manager = manager.WithDistributedLocker(locker)
   ```

5. **Sử dụng chế độ singleton cho các job quan trọng**:
   ```go
   manager.SingletonMode().Every(1).Day().Do(criticalDailyTask)
   ```
