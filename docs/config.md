# Cấu hình Scheduler

Scheduler package được thiết kế để dễ dàng cấu hình thông qua file config và code. Tài liệu này mô tả chi tiết các tùy chọn cấu hình có sẵn.

## Cấu trúc Config

Scheduler sử dụng struct `Config` để cấu hình các tùy chọn:

```go
type Config struct {
    // AutoStart xác định có tự động khởi động scheduler khi ứng dụng boot không
    // ServiceProvider sẽ tự động gọi scheduler.StartAsync() trong Boot() method nếu true
    AutoStart bool `mapstructure:"auto_start" yaml:"auto_start"`

    // DistributedLock chứa cấu hình cho distributed locking
    DistributedLock DistributedLockConfig `mapstructure:"distributed_lock" yaml:"distributed_lock"`

    // Options chứa cấu hình RedisLockerOptions cho distributed locking
    Options RedisLockerOptions `mapstructure:"options" yaml:"options"`
}
```

### DistributedLockConfig

```go
type DistributedLockConfig struct {
    // Enabled xác định có bật distributed locking không
    // Chỉ cần thiết khi chạy scheduler trên nhiều instance trong môi trường phân tán
    Enabled bool `mapstructure:"enabled" yaml:"enabled"`
}
```

### RedisLockerOptions

```go
type RedisLockerOptions struct {
    // KeyPrefix là tiền tố được thêm vào trước mỗi khóa trong Redis
    KeyPrefix string `mapstructure:"key_prefix" yaml:"key_prefix"`

    // LockDuration là thời gian một khóa sẽ tồn tại trước khi tự động hết hạn (giây)
    LockDuration int `mapstructure:"lock_duration" yaml:"lock_duration"`

    // MaxRetries là số lần thử tối đa khi gặp lỗi khi tương tác với Redis
    MaxRetries int `mapstructure:"max_retries" yaml:"max_retries"`

    // RetryDelay là thời gian chờ giữa các lần thử (milliseconds)
    RetryDelay int `mapstructure:"retry_delay" yaml:"retry_delay"`
}
```

## Giá trị mặc định

Giá trị mặc định được cung cấp thông qua các hàm:

```go
// DefaultConfig trả về cấu hình mặc định cho scheduler
func DefaultConfig() Config {
    return Config{
        AutoStart: true,
        DistributedLock: DistributedLockConfig{
            Enabled: false,
        },
        Options: DefaultRedisLockerOptions(),
    }
}

// DefaultRedisLockerOptions trả về các tùy chọn mặc định cho Redis Locker
func DefaultRedisLockerOptions() RedisLockerOptions {
    return RedisLockerOptions{
        KeyPrefix:    "scheduler_lock:",
        LockDuration: 30, // 30 seconds
        MaxRetries:   3,
        RetryDelay:   100, // 100 milliseconds
    }
}
```

## Cấu hình qua File

### Định dạng YAML

```yaml
scheduler:
  # Tự động khởi động scheduler khi ứng dụng boot
  auto_start: true

  # Distributed locking với Redis
  distributed_lock:
    enabled: true
  
  # Cài đặt Redis Locker
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60      # seconds
    max_retries: 5
    retry_delay: 200       # milliseconds
```

### Định dạng JSON

```json
{
  "scheduler": {
    "auto_start": true,
    "distributed_lock": {
      "enabled": true
    },
    "options": {
      "key_prefix": "myapp_scheduler:",
      "lock_duration": 60,
      "max_retries": 5,
      "retry_delay": 200
    }
  }
}
```

## Cấu hình qua Code

### Sử dụng DefaultConfig

```go
cfg := scheduler.DefaultConfig()
cfg.AutoStart = false
cfg.DistributedLock.Enabled = true
cfg.Options.LockDuration = 60
cfg.Options.MaxRetries = 5

manager := scheduler.NewSchedulerWithConfig(cfg)
```

### Tùy chỉnh RedisLockerOptions

```go
options := scheduler.RedisLockerOptions{
    KeyPrefix:    "custom_prefix:",
    LockDuration: 120,     // 2 minutes
    MaxRetries:   10,
    RetryDelay:   500,     // 500ms
}

// Convert to time.Duration based options
timeOptions := options.ToTimeDuration()
```

## Tích hợp với Config Package

Khi sử dụng ServiceProvider, scheduler sẽ tự động load cấu hình từ config provider:

```go
// Đăng ký các providers theo thứ tự phụ thuộc
app.Register(config.NewServiceProvider())
app.Register(redis.NewServiceProvider())  // Required cho distributed locking
app.Register(scheduler.NewServiceProvider())

// Khởi động ứng dụng - scheduler sẽ tự động tải cấu hình
app.Boot()
```

## Môi trường cụ thể

Bạn có thể tùy chỉnh cấu hình dựa trên môi trường:

```go
func getSchedulerConfig(env string) scheduler.Config {
    config := scheduler.DefaultConfig()
    
    switch env {
    case "production":
        config.Options.LockDuration = 300 // 5 phút
        config.Options.MaxRetries = 10
        config.Options.RetryDelay = 1000  // 1 giây
    case "development":
        config.Options.LockDuration = 30  // 30 giây
        config.Options.MaxRetries = 3
        config.Options.RetryDelay = 100   // 100ms
    }
    
    return config
}
```
