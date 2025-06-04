package scheduler

import (
	"go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/redis"
)

// ServiceProvider cung cấp dịch vụ scheduler và tích hợp với DI container.
//
// ServiceProvider là một implementation của interface di.ServiceProvider, cho phép tự động
// đăng ký scheduler manager vào DI container của ứng dụng. ServiceProvider thực hiện công việc:
//   - Tạo một scheduler manager mới sử dụng gocron
//   - Đăng ký scheduler manager vào DI container với key "scheduler"
//
// Việc cấu hình cụ thể và đăng ký các task được thực hiện bởi ứng dụng.
//
// Để sử dụng ServiceProvider, ứng dụng cần:
//   - Implement interface Container() *di.Container để cung cấp DI container
type ServiceProvider struct {
	providers []string
}

// NewServiceProvider trả về một ServiceProvider mới cho module scheduler.
//
// Hàm này khởi tạo và trả về một đối tượng ServiceProvider để sử dụng với DI container.
// ServiceProvider cho phép tự động đăng ký và cấu hình scheduler manager cho ứng dụng.
//
// Returns:
//   - di.ServiceProvider: Interface di.ServiceProvider đã được implement bởi ServiceProvider
//
// Example:
//
//	app.Register(scheduler.NewServiceProvider())
func NewServiceProvider() di.ServiceProvider {
	return &ServiceProvider{}
}

// Register đăng ký scheduler vào DI container.
//
// Register được gọi khi đăng ký ServiceProvider vào ứng dụng. Phương thức này
// tạo một scheduler manager mới và đăng ký vào DI container của ứng dụng.
//
// Params:
//   - app: di.Application - Đối tượng ứng dụng implements di.Application interface
//
// Luồng thực thi:
//  1. Lấy container từ app
//  2. Load cấu hình scheduler
//  3. Tạo scheduler manager mới
//  4. Cấu hình distributed locking nếu được bật
//  5. Đăng ký scheduler manager vào container với key "scheduler"
//
// Việc cấu hình và đăng ký các task sẽ được thực hiện bởi ứng dụng,
// cho phép mỗi ứng dụng tùy chỉnh scheduler theo nhu cầu riêng.
//
// Panics:
//   - Nếu không thể lấy container từ application
//   - Nếu không thể tạo scheduler manager
//   - Nếu không thể đăng ký scheduler vào container
//   - Nếu distributed locking được bật nhưng không thể cấu hình Redis locker
func (p *ServiceProvider) Register(app di.Application) {
	container := app.Container()
	if container == nil {
		panic("scheduler: DI container is nil - cannot register scheduler service")
	}

	// Load cấu hình scheduler với default fallback
	cfg := DefaultConfig()

	// Thử lấy cấu hình từ config provider (optional)
	if configInstance, err := container.Make("config"); err == nil {
		if configManager, ok := configInstance.(config.Manager); ok {
			// Load cấu hình từ file config với error handling
			if err := configManager.UnmarshalKey("scheduler", &cfg); err != nil {
				panic("scheduler: failed to load scheduler configuration: " + err.Error())
			}
		}
	}

	// Tạo scheduler manager với cấu hình
	manager := NewSchedulerWithConfig(cfg)
	if manager == nil {
		panic("scheduler: failed to create scheduler manager with config")
	}

	// Cấu hình distributed locking nếu được bật
	if cfg.DistributedLock.Enabled {
		redisInstance, err := container.Make("redis")
		if err != nil {
			panic("scheduler: distributed locking is enabled but redis service not found: " + err.Error())
		}

		redisManager, ok := redisInstance.(redis.Manager)
		if !ok {
			panic("scheduler: redis service is not a valid redis.Manager interface")
		}

		redisClient, err := redisManager.Client()
		if err != nil {
			panic("scheduler: failed to get redis client for distributed locking: " + err.Error())
		}

		locker, err := NewRedisLocker(redisClient, cfg.Options)
		if err != nil {
			panic("scheduler: failed to create Redis locker: " + err.Error())
		}

		manager = manager.WithDistributedLocker(locker)
		if manager == nil {
			panic("scheduler: failed to configure distributed locker on scheduler manager")
		}
	}

	// Đăng ký scheduler manager vào container
	container.Instance("scheduler", manager)

	p.providers = append(p.providers, "scheduler")
}

// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
//
// Boot là một lifecycle hook của di.ServiceProvider mà thực hiện sau khi tất cả
// các service provider đã được đăng ký xong.
//
// Trong trường hợp của SchedulerServiceProvider, Boot thực hiện:
// 1. Lấy scheduler manager từ container
// 2. Load cấu hình scheduler để kiểm tra AutoStart
// 3. Tự động start scheduler nếu AutoStart được bật
//
// Params:
//   - app: di.Application - Đối tượng ứng dụng implements di.Application interface
//
// Panics:
//   - Nếu không thể lấy container từ application
//   - Nếu không tìm thấy scheduler trong container
//   - Nếu scheduler không đúng type Manager interface
//   - Nếu không thể load cấu hình scheduler
//   - Nếu AutoStart được bật nhưng không thể start scheduler
func (p *ServiceProvider) Boot(app di.Application) {
	container := app.Container()
	if container == nil {
		panic("scheduler: DI container is nil during boot phase")
	}

	// Lấy scheduler manager từ container
	instance, err := container.Make("scheduler")
	if err != nil {
		panic("scheduler: scheduler service not found in container during boot: " + err.Error())
	}

	scheduler, ok := instance.(Manager)
	if !ok {
		panic("scheduler: registered scheduler service is not a valid Manager interface")
	}

	// Kiểm tra xem scheduler đã được start chưa
	if scheduler.IsRunning() {
		return // Scheduler đã được start rồi, không cần làm gì thêm
	}

	// Lấy cấu hình để kiểm tra AutoStart
	cfg := DefaultConfig()

	// Thử lấy cấu hình từ config provider (required cho AutoStart)
	configInstance, err := container.Make("config")
	if err != nil {
		panic("scheduler: config service not found but required for AutoStart feature: " + err.Error())
	}

	configManager, ok := configInstance.(config.Manager)
	if !ok {
		panic("scheduler: config service is not a valid config.Manager interface")
	}

	// Load cấu hình từ file config với error handling
	if err := configManager.UnmarshalKey("scheduler", &cfg); err != nil {
		panic("scheduler: failed to load scheduler configuration during boot: " + err.Error())
	}

	// Chỉ auto start nếu được cấu hình để làm vậy
	if cfg.AutoStart {
		scheduler.StartAsync()
	}
}

func (p *ServiceProvider) Requires() []string {
	return []string{
		"config",
		"redis",
	}
}

func (p *ServiceProvider) Providers() []string {
	return p.providers
}
