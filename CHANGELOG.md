# Changelog

## [Unreleased]

### Added
- Automation scripts for release management and project maintenance
- CODEOWNERS, FUNDING, and issue templates
- Comprehensive CI/CD workflows

### Fixed
- Updated copyright attribution in LICENSE file
- Fixed references from mongodb to scheduler in CODEOWNERS, issue templates, and release workflow

### Changed
- Updated `go.fork.vn/config` from v0.1.0 to v0.1.3
- Updated `go.fork.vn/di` from v0.1.0 to v0.1.3
- Updated `go.fork.vn/redis` from v0.1.0 to v0.1.2

### Dependencies
- **go.fork.vn/config v0.1.3**: Latest configuration management improvements
- **go.fork.vn/di v0.1.3**: Enhanced dependency injection features
- **go.fork.vn/redis v0.1.2**: Updated Redis connectivity and distributed locking

## v0.1.0 - 2025-05-31

### Added
- **Task Scheduling System**: Comprehensive task scheduling and management system for Go applications
- **Multiple Scheduling Methods**: Support for interval-based, time-based, and cron expression scheduling
- **Distributed Locking**: Redis-based distributed locking for cluster environments with automatic renewal
- **Singleton Mode**: Prevent parallel execution of same jobs across distributed systems
- **Job Management**: Tag-based job grouping, cancellation, and health monitoring
- **DI Integration**: Seamless integration with Dependency Injection container
- **Configuration-Driven**: Auto-start and Redis locker options through configuration
- **Fluent API**: Easy-to-use fluent interface for job configuration
- **Thread-Safe**: Safe concurrent job management and execution
- **Resource Management**: Automatic cleanup and proper resource disposal
- **Performance Optimized**: Efficient handling of multiple concurrent tasks
- **Error Resilience**: Robust error handling and recovery mechanisms
- **Health Monitoring**: Built-in health checks for scheduled tasks
- **Timezone Support**: Full timezone support for cron expressions
- **Memory Optimization**: Leak prevention for scheduled and canceled jobs
- **Mockery Integration**: Complete MockManager interface with 27 mocked methods
- **Testing Support**: Testify mock framework with expecter interface for fluent test setup
- **Configuration Support**: Dual config system with int values for config files and time.Duration for internal use
- **Auto-start Feature**: Configurable scheduler startup behavior through configuration
- **Lock Management**: Automatic Redis lock renewal in background goroutines
- **API Flexibility**: Support for both programmatic and configuration-driven setup

### Technical Details
- Initial release as standalone module `go.fork.vn/scheduler`
- Repository located at `github.com/Fork/scheduler`
- Built with Go 1.23.9
- Full test coverage and comprehensive documentation
- Based on go-co-op/gocron with Redis integration
- Automatic lock renewal after 2/3 of expiration time
- Background goroutines for lock management
- Complete mockery integration for testing
- Uses mockery v2.53.4+ compatible configuration
- Compile-time interface checks for type safety
- Memory leak prevention for scheduled and canceled jobs
- Robust error handling and recovery mechanisms

### Dependencies
- `go.fork.vn/di`: Dependency injection integration
- `go.fork.vn/config`: Configuration management
- `go.fork.vn/redis`: Redis connectivity and distributed locking

[Unreleased]: github.com/go-fork/scheduler/compare/v0.1.0...HEAD
[v0.1.0]: github.com/go-fork/scheduler/releases/tag/v0.1.0

## Features
- Add mockery configuration and complete mocks support
- Generate MockManager for scheduler.Manager interface
- Support testify mock framework with expecter interface
- Enable fluent test setup with automatic cleanup

## Added
- .mockery.yaml configuration file
- mocks/manager.go with 27 mocked methods
- Complete interface coverage for testing

## Technical Details
- Uses mockery v2.53.4+ compatible configuration
- Follows same pattern as config module
- Includes with-expecter support for better test experience
- Automatic mock expectations cleanup in tests"

## v0.0.4 - 2025-05-28

### Added

- **Configuration-driven system**: Added comprehensive configuration support through `Config` struct and `RedisLockerOptions`
- **Auto-start feature**: Added `auto_start` configuration option to control scheduler startup behavior
- **Dual config system**: 
  - `RedisLockerOptions` with int values for config files (seconds/milliseconds)
  - `RedisLockerOptionsTime` with time.Duration for internal use
  - `ToTimeDuration()` conversion method between the two
- **Configuration-driven distributed locking**: Automatic Redis locker setup when enabled in config
- **Enhanced ServiceProvider**: Added `NewSchedulerWithConfig()` function for configuration-driven setup
- **Improved dependency management**: Added proper require dependencies for config and redis providers

### Changed

- **BREAKING**: Updated `RedisLockerOptions` to use int values instead of time.Duration for better config file compatibility
- **Provider restructure**: Completely restructured ServiceProvider following separation of concerns
  - Provider now only loads config and passes to manager (no longer handles manager logic)
  - Added automatic distributed locking setup when enabled in config
  - Fixed interface compliance with `Providers()` method instead of `Provides()`
- **Manager enhancement**: Added `NewSchedulerWithConfig()` to support configuration-driven initialization
- **Test compatibility**: Updated all tests to work with new int-based config values

### Fixed

- **Static analysis**: Fixed all staticcheck warnings by removing unnecessary type assertions
- **Interface compliance**: Ensured ServiceProvider properly implements `di.ServiceProvider` interface
- **Test reliability**: Fixed `locker_test.go` to work with new configuration system
- **Build system**: Updated `go.mod` with proper dependencies and replace directives

### Technical Improvements

- Used compile-time interface checks instead of runtime type assertions
- Improved error handling in provider setup
- Enhanced documentation with configuration examples
- Better separation between provider and manager responsibilities

## v0.0.3 - 2025-05-25

### Added

- Tích hợp toàn bộ tính năng của thư viện gocron vào DI container của ứng dụng
- Hỗ trợ nhiều loại lịch trình: theo khoảng thời gian, theo thời điểm cụ thể, biểu thức cron
- Hỗ trợ chế độ singleton để tránh chạy song song cùng một task
- Hỗ trợ distributed locking với Redis cho môi trường phân tán
- Hỗ trợ tag để nhóm và quản lý các task
- API fluent cho trải nghiệm lập trình dễ dàng
--
- Khóa Redis được tự động gia hạn sau khi đã chạy được 2/3 thời gian hết hạn
- Việc gia hạn xảy ra trong một goroutine riêng biệt
- Khi job hoàn thành, khóa sẽ được giải phóng

### Added
- Task cancellation API
- Health monitoring for scheduled tasks

## [v0.0.3] - 2025-05-25

### Added
- New Scheduler Provider for managing periodic tasks
- Support for cron expression syntax
- Redis integration for distributed scheduling
- Fluent interface for job configuration
- Lock mechanism for distributed environments

### Changed
- Performance optimization for handling multiple concurrent tasks
- Improved error handling in periodic tasks

### Fixed
- Timezone issues in cron expressions
- Memory leak when many jobs are scheduled and canceled

## [v0.0.2] - 2025-05-22

### Added
- Job scheduling system based on go-co-op/gocron
- Support for multiple scheduling methods:
  - Time interval-based scheduling
  - Specific time-based scheduling
  - Cron expression support
- Single-run and repeated task support
- Tag-based job grouping and management
- Distributed locking with Redis for cluster environments
- Singleton mode to prevent parallel execution of the same job
- Dependency Injection integration through ServiceProvider
- Fluent API for job configuration
