package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test default values
	assert.True(t, config.AutoStart, "AutoStart should be true by default")
	assert.False(t, config.DistributedLock.Enabled, "DistributedLock should be disabled by default")

	// Test default Redis locker options
	expectedOptions := DefaultRedisLockerOptions()
	assert.Equal(t, expectedOptions, config.Options, "Options should match default Redis locker options")
}

func TestRedisLockerOptionsToTimeDuration(t *testing.T) {
	options := RedisLockerOptions{
		KeyPrefix:    "test_prefix:",
		LockDuration: 60, // 60 seconds
		MaxRetries:   5,
		RetryDelay:   250, // 250 milliseconds
	}

	timeOptions := options.ToTimeDuration()

	// Test conversions
	assert.Equal(t, "test_prefix:", timeOptions.KeyPrefix, "KeyPrefix should remain unchanged")
	assert.Equal(t, 60*time.Second, timeOptions.LockDuration, "LockDuration should be converted to time.Duration")
	assert.Equal(t, 5, timeOptions.MaxRetries, "MaxRetries should remain unchanged")
	assert.Equal(t, 250*time.Millisecond, timeOptions.RetryDelay, "RetryDelay should be converted to time.Duration")
}

func TestRedisLockerOptionsToTimeDurationWithZeroValues(t *testing.T) {
	options := RedisLockerOptions{
		KeyPrefix:    "zero_test:",
		LockDuration: 0,
		MaxRetries:   0,
		RetryDelay:   0,
	}

	timeOptions := options.ToTimeDuration()

	// Test zero value conversions
	assert.Equal(t, "zero_test:", timeOptions.KeyPrefix)
	assert.Equal(t, 0*time.Second, timeOptions.LockDuration)
	assert.Equal(t, 0, timeOptions.MaxRetries)
	assert.Equal(t, 0*time.Millisecond, timeOptions.RetryDelay)
}

func TestConfigStructFields(t *testing.T) {
	// Test that Config struct can be properly instantiated
	config := Config{
		AutoStart: false,
		DistributedLock: DistributedLockConfig{
			Enabled: true,
		},
		Options: RedisLockerOptions{
			KeyPrefix:    "custom_prefix:",
			LockDuration: 45,
			MaxRetries:   2,
			RetryDelay:   150,
		},
	}

	// Verify all fields are set correctly
	assert.False(t, config.AutoStart)
	assert.True(t, config.DistributedLock.Enabled)
	assert.Equal(t, "custom_prefix:", config.Options.KeyPrefix)
	assert.Equal(t, 45, config.Options.LockDuration)
	assert.Equal(t, 2, config.Options.MaxRetries)
	assert.Equal(t, 150, config.Options.RetryDelay)
}

func TestDistributedLockConfig(t *testing.T) {
	// Test enabled
	enabledConfig := DistributedLockConfig{Enabled: true}
	assert.True(t, enabledConfig.Enabled)

	// Test disabled
	disabledConfig := DistributedLockConfig{Enabled: false}
	assert.False(t, disabledConfig.Enabled)
}

func TestRedisLockerOptionsTimeStruct(t *testing.T) {
	// Test that RedisLockerOptionsTime can be properly instantiated
	timeOptions := RedisLockerOptionsTime{
		KeyPrefix:    "time_test:",
		LockDuration: 2 * time.Minute,
		MaxRetries:   10,
		RetryDelay:   500 * time.Millisecond,
	}

	assert.Equal(t, "time_test:", timeOptions.KeyPrefix)
	assert.Equal(t, 2*time.Minute, timeOptions.LockDuration)
	assert.Equal(t, 10, timeOptions.MaxRetries)
	assert.Equal(t, 500*time.Millisecond, timeOptions.RetryDelay)
}

func TestConfigWithDefaultRedisOptions(t *testing.T) {
	// Test creating config with explicit default Redis options
	config := Config{
		AutoStart: true,
		DistributedLock: DistributedLockConfig{
			Enabled: true,
		},
		Options: DefaultRedisLockerOptions(),
	}

	expectedOptions := DefaultRedisLockerOptions()
	assert.Equal(t, expectedOptions, config.Options)
	assert.True(t, config.AutoStart)
	assert.True(t, config.DistributedLock.Enabled)
}

func TestRedisLockerOptionsValidValues(t *testing.T) {
	tests := []struct {
		name     string
		options  RedisLockerOptions
		expected RedisLockerOptionsTime
	}{
		{
			name: "Standard values",
			options: RedisLockerOptions{
				KeyPrefix:    "scheduler:",
				LockDuration: 30,
				MaxRetries:   3,
				RetryDelay:   100,
			},
			expected: RedisLockerOptionsTime{
				KeyPrefix:    "scheduler:",
				LockDuration: 30 * time.Second,
				MaxRetries:   3,
				RetryDelay:   100 * time.Millisecond,
			},
		},
		{
			name: "Large values",
			options: RedisLockerOptions{
				KeyPrefix:    "large_test:",
				LockDuration: 300, // 5 minutes
				MaxRetries:   100,
				RetryDelay:   5000, // 5 seconds
			},
			expected: RedisLockerOptionsTime{
				KeyPrefix:    "large_test:",
				LockDuration: 300 * time.Second,
				MaxRetries:   100,
				RetryDelay:   5000 * time.Millisecond,
			},
		},
		{
			name: "Minimal values",
			options: RedisLockerOptions{
				KeyPrefix:    "min:",
				LockDuration: 1,
				MaxRetries:   1,
				RetryDelay:   1,
			},
			expected: RedisLockerOptionsTime{
				KeyPrefix:    "min:",
				LockDuration: 1 * time.Second,
				MaxRetries:   1,
				RetryDelay:   1 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ToTimeDuration()
			assert.Equal(t, tt.expected, result)
		})
	}
}
