# Scheduler Package

Package `go.fork.vn/scheduler` cung cấp một hệ thống lập lịch linh hoạt và mạnh mẽ cho ứng dụng Go, với khả năng tích hợp dependency injection, distributed locking, và hỗ trợ các pattern cron phức tạp.

## Mục lục

- [Tổng quan](overview.md) - Tổng quan về kiến trúc và tính năng scheduler
- [Cấu hình](config.md) - Chi tiết về cấu hình và các tùy chọn
- [Provider](provider.md) - Cách sử dụng ServiceProvider với DI
- [Manager](manager.md) - API Manager và các phương thức lập lịch
- [Distributed Lock](with_distributed_lock.md) - Sử dụng distributed locking trong môi trường phân tán

## Cài đặt

```bash
go get go.fork.vn/scheduler@latest
```

## Sử dụng nhanh

```go
package main

import (
    "fmt"
    "time"
    
    "go.fork.vn/scheduler"
)

func main() {
    // Tạo scheduler manager
    manager := scheduler.NewScheduler()
    
    // Lập lịch job chạy mỗi 5 phút
    manager.Every(5).Minutes().Do(func() {
        fmt.Println("Task chạy mỗi 5 phút")
    })
    
    // Lập lịch job chạy vào 10:30 AM hàng ngày
    manager.At("10:30").Every(1).Day().Do(func() {
        fmt.Println("Task chạy lúc 10:30 AM hàng ngày")
    })
    
    // Lập lịch job sử dụng cron expression
    manager.Cron("0 */6 * * *").Do(func() {
        fmt.Println("Task chạy mỗi 6 giờ")
    })
    
    // Khởi động scheduler
    manager.Start()
    
    // Giữ ứng dụng chạy
    select {}
}
```

## Với DI Container

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
    
    // Lấy scheduler từ container
    container := app.Container()
    manager := container.Get("scheduler").(scheduler.Manager)
    
    // Sử dụng manager để lập lịch jobs
    // ...
    
    // Giữ ứng dụng chạy
    select {}
}
```

## Xem thêm tài liệu

Vui lòng xem các trang tài liệu chi tiết để biết thêm về các tính năng và cách sử dụng khác của scheduler package.
