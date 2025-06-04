# Release Notes - v0.1.1

## Overview
Phiên bản v0.1.1 tập trung vào việc cải thiện tài liệu, cập nhật dependencies và sửa lỗi cho package scheduler. Phiên bản này nâng cao tính ổn định và hiệu suất của hệ thống, đồng thời cải thiện trải nghiệm phát triển với tài liệu toàn diện hơn.

## What's New
### 🚀 Features
- Thêm automation scripts cho việc quản lý phát hành và bảo trì dự án
- Thêm CODEOWNERS, FUNDING và issue templates cho quản lý dự án tốt hơn
- Thêm comprehensive CI/CD workflows

### 🐛 Bug Fixes
- Sửa thông tin bản quyền trong LICENSE file
- Sửa các references từ mongodb đến scheduler trong CODEOWNERS, issue templates và release workflow
- Sửa vấn đề với ServiceProvider Interface để tương thích với go.fork.vn/di v0.1.3
- Sửa lỗi type mismatch trong provider_test.go (*scheduler.schedulerManager -> *scheduler.manager)
- Cải thiện distributed lock test để xử lý Redis client validation

### 🔧 Improvements
- Nâng cao xử lý lỗi với panic messages rõ ràng cho các lỗi quan trọng
- Tăng cường test coverage với thêm config_test.go

### 📚 Documentation
- Tái cấu trúc toàn bộ tài liệu thành các module có tổ chức: index, overview, config, provider, manager và with_distributed_lock
- Cải thiện hướng dẫn sử dụng distributed lock với ví dụ chi tiết
- Thêm tài liệu về cách cấu hình và troubleshooting

## Breaking Changes
### ⚠️ Important Notes
Không có breaking changes trong phiên bản này.

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- go.fork.vn/config: v0.1.0 → v0.1.3
- go.fork.vn/di: v0.1.0 → v0.1.3
- go.fork.vn/redis: v0.1.0 → v0.1.2

### Dependencies details
- **go.fork.vn/config v0.1.3**: Latest configuration management improvements
- **go.fork.vn/di v0.1.3**: Enhanced dependency injection features
- **go.fork.vn/redis v0.1.2**: Updated Redis connectivity and distributed locking

## Performance
- Benchmark improvement: X% faster in scenario Y
- Memory usage: X% reduction in scenario Z

## Security
- Security fix for vulnerability X
- Updated dependencies with security patches

## Testing
- Added X new test cases
- Improved test coverage to X%

## Contributors
Thanks to all contributors who made this release possible:
- @contributor1
- @contributor2

## Download
- Source code: [go.fork.vn/scheduler@v0.1.1]
- Documentation: [pkg.go.dev/go.fork.vn/scheduler@v0.1.1]

---
Release Date: 2025-06-04
