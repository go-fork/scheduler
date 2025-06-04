# Release Notes - v0.1.0

## Overview
Phiên bản đầu tiên chính thức của Go Scheduler package, cung cấp hệ thống lập lịch mạnh mẽ và linh hoạt cho các ứng dụng Go với tích hợp dependency injection, distributed locking và hỗ trợ các pattern cron phức tạp.

## What's New
### 🚀 Features
- **Task Scheduling System**: Hệ thống lập lịch và quản lý task toàn diện cho ứng dụng Go
- **Multiple Scheduling Methods**: Hỗ trợ lập lịch theo khoảng thời gian, theo thời điểm cụ thể, và biểu thức cron
- **Distributed Locking**: Redis-based distributed locking cho môi trường cluster với auto-renewal
- **Singleton Mode**: Ngăn chặn thực thi song song của cùng một job trên nhiều hệ thống phân tán
- **Job Management**: Nhóm job theo tag, hủy bỏ và giám sát sức khỏe
- **DI Integration**: Tích hợp liền mạch với Dependency Injection container
- **Configuration-Driven**: Tùy chọn auto-start và Redis locker thông qua cấu hình
- **Fluent API**: Giao diện dễ sử dụng để cấu hình job

### 📚 Documentation
- Tài liệu đầy đủ về cách sử dụng và cấu hình scheduler
- Ví dụ tích hợp với hệ thống DI container
- Hướng dẫn chi tiết về distributed locking với Redis
- Tài liệu API cho tất cả các chức năng chính

## Breaking Changes
### ⚠️ Important Notes
Không có breaking changes vì đây là phiên bản đầu tiên.

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Added
- go.fork.vn/di: v0.1.3
- go.fork.vn/config: v0.1.3
- go.fork.vn/redis: v0.1.2
- go-co-op/gocron: v1.37.0
- github.com/redis/go-redis/v9: v9.9.0

## Performance
- **Thread-Safe**: Xử lý đồng thời job một cách an toàn
- **Resource Management**: Tự động dọn dẹp và giải phóng tài nguyên đúng cách
- **Performance Optimized**: Xử lý hiệu quả nhiều task đồng thời
- **Memory Optimization**: Ngăn leak cho job đã lập lịch và bị hủy

## Security
- Cơ chế distributed locking an toàn với Redis Lua scripts
- Xử lý key trong Redis với prefix tùy chỉnh để tránh xung đột

## Testing
- Hơn 50 test cases bao gồm tất cả các chức năng chính
- Test coverage trên 80% cho toàn bộ package
- MockManager interface với 27 phương thức được mock
- Hỗ trợ testify mock framework với expecter interface

## Contributors
Thanks to all contributors who made this release possible:
- @cluster
- @fork-team

## Download
- Source code: [github.com/go-fork/scheduler/releases/tag/v0.1.0](https://github.com/go-fork/scheduler/releases/tag/v0.1.0)
- Documentation: [pkg.go.dev/go.fork.vn/scheduler@v0.1.0](https://pkg.go.dev/go.fork.vn/scheduler@v0.1.0)

---
Release Date: 2025-06-04
