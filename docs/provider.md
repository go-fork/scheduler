# Service Provider

Scheduler package cung cấp một ServiceProvider để tích hợp liền mạch với DI container của ứng dụng. ServiceProvider này tự động đăng ký, cấu hình và khởi động scheduler dựa trên file cấu hình của ứng dụng.

## Đăng ký ServiceProvider

```go
package main

import (
    "go.fork.vn/di"
    "go.fork.vn/scheduler"
)

func main() {
    app := di.New()
    app.Register(scheduler.NewServiceProvider())
    app.Boot()
    
    // Giữ ứng dụng chạy
    select {}
}
```

## Dependencies

ServiceProvider khai báo các dependencies của nó thông qua phương thức `Requires()`:

```go
func (p *ServiceProvider) Requires() []string {
    return []string{"config", "redis"}
}
```

- `config`: Cần thiết để tải cấu hình scheduler
- `redis`: Cần thiết khi sử dụng distributed locking (nếu bật trong cấu hình)

## Các dịch vụ đăng ký

ServiceProvider đăng ký các dịch vụ sau vào container:

```go
func (p *ServiceProvider) Providers() []string {
    return []string{"scheduler"}
}
```

- `scheduler`: Một instance của `scheduler.Manager`, có thể được truy xuất từ container

## Luồng thực thi

### Register Method

Phương thức `Register` được gọi khi đăng ký ServiceProvider vào ứng dụng:

```go
func (p *ServiceProvider) Register(app di.Application) {
    // 1. Lấy container từ app
    container := app.Container()
    
    // 2. Load cấu hình scheduler từ config provider
    cfg := DefaultConfig()
    if configManager, err := container.Make("config"); err == nil {
        configManager.(config.Manager).UnmarshalKey("scheduler", &cfg)
    }
    
    // 3. Tạo scheduler manager mới với cấu hình
    manager := NewSchedulerWithConfig(cfg)
    
    // 4. Cấu hình Redis-based distributed locking nếu được bật
    if cfg.DistributedLock.Enabled {
        redisManager := container.Make("redis").(redis.Manager)
        redisClient, _ := redisManager.Client()
        
        locker, _ := NewRedisLocker(redisClient, cfg.Options)
        manager = manager.WithDistributedLocker(locker)
    }
    
    // 5. Đăng ký scheduler manager vào container
    container.Instance("scheduler", manager)
    
    p.providers = append(p.providers, "scheduler")
}
```

### Boot Method

Phương thức `Boot` được gọi sau khi tất cả các ServiceProvider đã được đăng ký:

```go
func (p *ServiceProvider) Boot(app di.Application) {
    // 1. Lấy container từ app
    container := app.Container()
    
    // 2. Lấy scheduler manager từ container
    manager := container.Make("scheduler").(Manager)
    
    // 3. Tải cấu hình để kiểm tra tùy chọn AutoStart
    cfg := DefaultConfig()
    if configManager, err := container.Make("config"); err == nil {
        configManager.(config.Manager).UnmarshalKey("scheduler", &cfg)
    }
    
    // 4. Tự động khởi động scheduler nếu AutoStart được bật
    if cfg.AutoStart {
        manager.Start()
    }
}
```

## Xử lý lỗi

ServiceProvider xử lý lỗi và panic trong các tình huống quan trọng:

1. Khi container không khả dụng:
   ```go
   if container == nil {
       panic("scheduler: DI container is nil - cannot register scheduler service")
   }
   ```

2. Khi không thể tạo scheduler manager:
   ```go
   if manager == nil {
       panic("scheduler: failed to create scheduler manager with config")
   }
   ```

3. Khi distributed locking được bật nhưng Redis không khả dụng:
   ```go
   if err != nil {
       panic("scheduler: distributed locking is enabled but redis service not found: " + err.Error())
   }
   ```

## Các tùy chọn cấu hình

ServiceProvider sử dụng các tùy chọn cấu hình từ `config.Config`:

```yaml
scheduler:
  auto_start: true
  distributed_lock:
    enabled: true
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60
    max_retries: 5
    retry_delay: 200
```

## Sử dụng từ Container

Sau khi ServiceProvider đã được đăng ký và khởi động, bạn có thể lấy scheduler manager từ container:

```go
container := app.Container()
manager := container.Make("scheduler").(scheduler.Manager)

// Sử dụng manager để lập lịch jobs
manager.Every(5).Minutes().Do(func() {
    fmt.Println("Task chạy mỗi 5 phút")
})
```
