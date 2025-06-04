package scheduler

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	configMocks "go.fork.vn/config/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
	redisMocks "go.fork.vn/redis/mocks"
)

func TestServiceProviderRegister(t *testing.T) {
	// Tạo mock objects
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)

	// Setup expectations
	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("config").Return(nil, assert.AnError)
	mockContainer.EXPECT().Instance("scheduler", mock.AnythingOfType("*scheduler.manager"))

	// Tạo service provider
	provider := NewServiceProvider()

	// Test Register
	provider.Register(mockApp)

	// Verify expectations
	mockApp.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
}

func TestServiceProviderRegisterWithConfig(t *testing.T) {
	// Tạo mock objects
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)
	mockConfig := configMocks.NewMockManager(t)

	// Setup expectations
	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("config").Return(mockConfig, nil)
	mockConfig.EXPECT().UnmarshalKey("scheduler", mock.AnythingOfType("*scheduler.Config")).Return(nil)
	mockContainer.EXPECT().Instance("scheduler", mock.AnythingOfType("*scheduler.manager"))

	// Tạo service provider
	provider := NewServiceProvider()

	// Test Register với config
	provider.Register(mockApp)

	// Verify expectations
	mockApp.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

func TestServiceProviderRegisterWithDistributedLock(t *testing.T) {
	// Test that distributed lock configuration panics when Redis client ping fails
	// since we can't easily mock a working Redis client in unit tests

	// Tạo mock objects
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)
	mockConfig := configMocks.NewMockManager(t)
	mockRedis := redisMocks.NewMockManager(t)
	mockRedisClient := &redis.Client{} // This will be nil and cause ping to fail

	// Tạo config với distributed lock enabled
	cfg := DefaultConfig()
	cfg.DistributedLock.Enabled = true

	// Setup expectations
	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("config").Return(mockConfig, nil)
	mockConfig.EXPECT().UnmarshalKey("scheduler", mock.AnythingOfType("*scheduler.Config")).Run(func(key string, target interface{}) {
		// Simulate config loading
		if config, ok := target.(*Config); ok {
			*config = cfg
		}
	}).Return(nil)
	mockContainer.EXPECT().Make("redis").Return(mockRedis, nil)
	mockRedis.EXPECT().Client().Return(mockRedisClient, nil)

	// Tạo service provider
	provider := NewServiceProvider()

	// Test should panic because Redis ping will fail with nil client
	assert.Panics(t, func() {
		provider.Register(mockApp)
	})

	// Verify expectations (note: Instance expectation removed since it panics before reaching that point)
	mockApp.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestServiceProviderRegisterPanics(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *diMocks.MockApplication
	}{
		{
			name: "panic when container is nil",
			setupMock: func() *diMocks.MockApplication {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.EXPECT().Container().Return(nil)
				return mockApp
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := tt.setupMock()
			provider := NewServiceProvider()

			assert.Panics(t, func() {
				provider.Register(mockApp)
			})
		})
	}
}

func TestServiceProviderBoot(t *testing.T) {
	// Tạo mock objects
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)
	mockConfig := configMocks.NewMockManager(t)

	// Tạo scheduler manager cho test
	manager := NewScheduler()

	// Setup expectations
	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("scheduler").Return(manager, nil)
	mockContainer.EXPECT().Make("config").Return(mockConfig, nil)
	mockConfig.EXPECT().UnmarshalKey("scheduler", mock.AnythingOfType("*scheduler.Config")).Return(nil)

	// Tạo service provider
	provider := NewServiceProvider()

	// Test Boot
	provider.Boot(mockApp)

	// Verify expectations
	mockApp.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
	mockConfig.AssertExpectations(t)

	// Cleanup
	if manager.IsRunning() {
		manager.Stop()
	}
}

func TestServiceProviderBootWithAutoStart(t *testing.T) {
	// Tạo mock objects
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)
	mockConfig := configMocks.NewMockManager(t)

	// Tạo scheduler manager cho test
	manager := NewScheduler()

	// Tạo config với AutoStart enabled
	cfg := DefaultConfig()
	cfg.AutoStart = true

	// Setup expectations
	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("scheduler").Return(manager, nil)
	mockContainer.EXPECT().Make("config").Return(mockConfig, nil)
	mockConfig.EXPECT().UnmarshalKey("scheduler", mock.AnythingOfType("*scheduler.Config")).Run(func(key string, target interface{}) {
		// Simulate config loading
		if config, ok := target.(*Config); ok {
			*config = cfg
		}
	}).Return(nil)

	// Tạo service provider
	provider := NewServiceProvider()

	// Test Boot với AutoStart
	provider.Boot(mockApp)

	// Verify scheduler đã được start
	assert.True(t, manager.IsRunning(), "Scheduler should be running with AutoStart enabled")

	// Verify expectations
	mockApp.AssertExpectations(t)
	mockContainer.AssertExpectations(t)
	mockConfig.AssertExpectations(t)

	// Cleanup
	manager.Stop()
}

func TestServiceProviderBootPanics(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *diMocks.MockApplication
	}{
		{
			name: "panic when container is nil",
			setupMock: func() *diMocks.MockApplication {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.EXPECT().Container().Return(nil)
				return mockApp
			},
		},
		{
			name: "panic when scheduler not found",
			setupMock: func() *diMocks.MockApplication {
				mockApp := diMocks.NewMockApplication(t)
				mockContainer := diMocks.NewMockContainer(t)
				mockApp.EXPECT().Container().Return(mockContainer)
				mockContainer.EXPECT().Make("scheduler").Return(nil, assert.AnError)
				return mockApp
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := tt.setupMock()
			provider := NewServiceProvider()

			assert.Panics(t, func() {
				provider.Boot(mockApp)
			})
		})
	}
}

func TestServiceProviderRequires(t *testing.T) {
	provider := NewServiceProvider()
	requires := provider.Requires()

	expectedRequires := []string{"config", "redis"}
	assert.Equal(t, expectedRequires, requires)
}

func TestServiceProviderProviders(t *testing.T) {
	provider := NewServiceProvider()

	// Initially should be empty
	providers := provider.Providers()
	assert.Empty(t, providers)

	// After registration should contain "scheduler"
	mockApp := diMocks.NewMockApplication(t)
	mockContainer := diMocks.NewMockContainer(t)

	mockApp.EXPECT().Container().Return(mockContainer)
	mockContainer.EXPECT().Make("config").Return(nil, assert.AnError)
	mockContainer.EXPECT().Instance("scheduler", mock.AnythingOfType("*scheduler.manager"))

	provider.Register(mockApp)

	providers = provider.Providers()
	assert.Contains(t, providers, "scheduler")
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()

	assert.NotNil(t, provider)

	// Kiểm tra provider implement đúng interface
	var _ di.ServiceProvider = provider
}
