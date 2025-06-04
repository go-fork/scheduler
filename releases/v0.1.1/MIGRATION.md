# Migration Guide - v0.1.1

## Overview
Hướng dẫn này giúp bạn nâng cấp từ phiên bản v0.1.0 lên v0.1.1. Phiên bản v0.1.1 tập trung vào cải thiện tài liệu, cập nhật dependencies và sửa lỗi, không có breaking changes.

## Prerequisites
- Go 1.18 hoặc mới hơn
- Phiên bản v0.1.0 đã cài đặt

## Quick Migration Checklist
- [ ] Cập nhật dependency đến v0.1.1
- [ ] Kiểm tra mã nguồn với go.fork.vn/di v0.1.3
- [ ] Chạy tests để đảm bảo tương thích
- [ ] Xem tài liệu mới cho các tính năng cập nhật

## Breaking Changes
Phiên bản v0.1.1 không có breaking changes so với v0.1.0.

### Interface Updates
ServiceProvider interface đã được cập nhật để phù hợp với go.fork.vn/di v0.1.3. Nếu bạn đang triển khai interface này, hãy cập nhật chữ ký hàm:

```go
// Old way (v0.1.0)
func (p *YourProvider) Register(app interface{})
func (p *YourProvider) Boot(app interface{})

// New way (v0.1.1)
func (p *YourProvider) Register(app di.Application)
func (p *YourProvider) Boot(app di.Application)
```

### Documentation Changes
Tài liệu đã được tái cấu trúc và mở rộng. Ví dụ mã cũ vẫn hoạt động, nhưng bạn nên sử dụng tài liệu mới để tìm hiểu các tính năng chi tiết hơn.

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get -u go.fork.vn/scheduler@v0.1.1
go get -u go.fork.vn/di@v0.1.3
go get -u go.fork.vn/config@v0.1.3
go get -u go.fork.vn/redis@v0.1.2
go mod tidy
```

### Step 2: Kiểm tra ServiceProvider nếu đã triển khai
Nếu bạn đã mở rộng hoặc triển khai ServiceProvider, hãy cập nhật chữ ký hàm:

```go
// Cập nhật từ interface{} sang di.Application
func (p *YourProvider) Register(app di.Application) {
    // Your code here
}

// Cập nhật từ interface{} sang di.Application
func (p *YourProvider) Boot(app di.Application) {
    // Your code here
}
```

### Step 3: Xem tài liệu mới
Khám phá tài liệu mới được tổ chức trong thư mục docs:

```bash
ls -la docs/
```

### Step 4: Chạy Tests
```bash
go test ./...
```

## Vấn đề thường gặp và giải pháp

### Issue 1: Kiểu không tương thích với ServiceProvider
**Vấn đề**: `cannot use app (variable of type interface{}) as di.Application value`  
**Giải pháp**: Cập nhật chữ ký hàm Register và Boot như hướng dẫn ở trên

### Issue 2: Lỗi panic khi inject container
**Vấn đề**: Scheduler panic khi inject container vì không xử lý nil container
**Giải pháp**: Đã được sửa trong v0.1.1, cập nhật phiên bản sẽ giải quyết vấn đề

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/scheduler@v0.1.1)
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
