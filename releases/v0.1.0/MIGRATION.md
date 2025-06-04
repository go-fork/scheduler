# Migration Guide - v0.1.0

## Overview
Đây là phiên bản đầu tiên của Go Scheduler package, nên không có migration cần thiết từ phiên bản trước đó. Hướng dẫn này giải thích cách cài đặt và sử dụng package.

## Prerequisites
- Go 1.18 hoặc mới hơn
- Redis (chỉ khi sử dụng distributed locking)

## Quick Installation Checklist
- [ ] Add go.fork.vn/scheduler vào go.mod
- [ ] Khởi tạo scheduler trong ứng dụng của bạn
- [ ] Cấu hình thông qua config hoặc code
- [ ] Chạy tests để đảm bảo tương thích

## Getting Started

### Installation
```go
go get go.fork.vn/scheduler@v0.1.0
go mod tidy
```

### Basic Usage
```go
import (
    "go.fork.vn/scheduler"
)

// Tạo mới scheduler
manager := scheduler.NewScheduler()

// Lập lịch job
manager.Every(5).Minutes().Do(func() {
    fmt.Println("Task runs every 5 minutes")
})

// Khởi động scheduler
manager.Start()
```

### DI Container Integration

```go
import (
    "go.fork.vn/di"
    "go.fork.vn/scheduler"
)

// Đăng ký provider
app := di.New()
app.Register(scheduler.NewServiceProvider())
app.Boot()

// Lấy scheduler từ container
container := app.Container()
sched := container.Get("scheduler").(scheduler.Manager)

// Lập lịch tasks
sched.Every(10).Minutes().Do(func() {
    fmt.Println("Task runs every 10 minutes")
})
```

### Configuration Example

```yaml
# config/app.yaml
scheduler:
  # Tự động khởi động khi boot
  auto_start: true
  
  # Distributed locking
  distributed_lock:
    enabled: true
    
  # Redis locker options
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60
    max_retries: 5
    retry_delay: 200
```

## Important Notes

### RedisLockerOptions
Cấu hình của Redis Locker sử dụng các giá trị int thay vì time.Duration để tương thích tốt hơn với file cấu hình:
- `lock_duration`: Thời gian lock tính bằng giây
- `retry_delay`: Thời gian giữa các lần retry tính bằng milliseconds

### Error Handling
Luôn kiểm tra lỗi khi sử dụng các phương thức có thể trả về error:

```go
if err := sched.Start(); err != nil {
    log.Fatal("Failed to start scheduler:", err)
}
```

## Getting Help
Nếu bạn gặp vấn đề khi sử dụng package, vui lòng mở issue trên GitHub repository.

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/scheduler@v0.1.0)
- Search [existing issues](https://github.com/go-fork/scheduler/issues)
- Create a [new issue](https://github.com/go-fork/scheduler/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/scheduler@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
