# Distributed Locking

Scheduler package cung cấp một hệ thống khóa phân tán dựa trên Redis để đảm bảo rằng các jobs được lập lịch chỉ chạy trên một instance trong môi trường phân tán. Tài liệu này mô tả cách sử dụng và cấu hình tính năng distributed locking.

## Tổng quan về Distributed Locking

Khi chạy ứng dụng trên nhiều máy chủ hoặc containers, các jobs được lập lịch có thể chạy trùng lặp trên nhiều instance, gây ra các vấn đề về tính nhất quán dữ liệu. Distributed locking đảm bảo rằng một job chỉ chạy trên một instance tại một thời điểm.

## Cách hoạt động

Scheduler sử dụng Redis để triển khai distributed locking:

1. Trước khi thực thi job, scheduler cố gắng lấy một khóa trong Redis
2. Nếu lấy được khóa, job sẽ được thực thi
3. Nếu không lấy được khóa (đã được instance khác lấy), job sẽ bị bỏ qua trên instance hiện tại
4. Sau khi job hoàn thành, khóa sẽ được giải phóng

## Kích hoạt Distributed Locking qua Config

Cách đơn giản nhất để kích hoạt distributed locking là thông qua file cấu hình:

```yaml
scheduler:
  distributed_lock:
    enabled: true
  options:
    key_prefix: "myapp_scheduler:"
    lock_duration: 60      # seconds
    max_retries: 5
    retry_delay: 200       # milliseconds
```

Khi sử dụng ServiceProvider, tính năng distributed locking sẽ được tự động cấu hình dựa trên cài đặt này.

## Cấu hình Distributed Locking bằng Code

```go
import (
    "github.com/redis/go-redis/v9"
    "go.fork.vn/scheduler"
)

// Khởi tạo Redis client
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Tạo Redis Locker với tùy chọn tùy chỉnh
options := scheduler.RedisLockerOptions{
    KeyPrefix:    "myapp_scheduler:",
    LockDuration: 60,   // seconds
    MaxRetries:   5,
    RetryDelay:   200,  // milliseconds
}

// Tạo locker
locker, err := scheduler.NewRedisLocker(redisClient, options)
if err != nil {
    log.Fatal("Failed to create Redis locker:", err)
}

// Tạo scheduler với distributed locker
manager := scheduler.NewScheduler().WithDistributedLocker(locker)

// Lập lịch jobs như bình thường
manager.Every(5).Minutes().Do(func() {
    fmt.Println("This job will only run on one instance at a time")
})

// Khởi động scheduler
manager.Start()
```

## Chi tiết Cài đặt

### RedisLockerOptions

```go
type RedisLockerOptions struct {
    // KeyPrefix là tiền tố được thêm vào trước mỗi khóa trong Redis
    KeyPrefix string

    // LockDuration là thời gian một khóa sẽ tồn tại trước khi tự động hết hạn (giây)
    LockDuration int

    // MaxRetries là số lần thử tối đa khi gặp lỗi khi tương tác với Redis
    MaxRetries int

    // RetryDelay là thời gian chờ giữa các lần thử (milliseconds)
    RetryDelay int
}
```

### Cơ chế khóa

Scheduler sử dụng Redis với các lệnh `SET NX` và Lua scripts để đảm bảo các thao tác nguyên tử:

```go
// Lấy khóa
success, err := client.SetNX(ctx, lockKey, lockValue, options.LockDuration).Result()

// Giải phóng khóa
script := `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`
result, err := client.Eval(ctx, script, []string{lockKey}, lockValue).Result()
```

## Tính năng Tự động Gia hạn Khóa

Khi một job chạy lâu hơn thời gian `LockDuration`, scheduler triển khai cơ chế tự động gia hạn khóa:

1. Sau khi đã chạy được 2/3 thời gian `LockDuration`, khóa sẽ được tự động gia hạn
2. Việc gia hạn xảy ra trong một goroutine riêng biệt
3. Khi job hoàn thành, khóa sẽ được giải phóng
4. Nếu instance gặp sự cố, khóa sẽ tự động hết hạn sau `LockDuration`

## Best Practices

### 1. Chọn LockDuration phù hợp

`LockDuration` nên dài hơn thời gian chạy dự kiến của job:

```go
options := scheduler.RedisLockerOptions{
    LockDuration: 300,  // 5 phút cho jobs dài
    MaxRetries:   3,
    RetryDelay:   100,
    KeyPrefix:    "scheduler:",
}
```

### 2. Xử lý lỗi trong jobs

```go
manager.Every(5).Minutes().Do(func() {
    defer func() {
        if r := recover(); r != nil {
            log.Error("Job panic: %v", r)
        }
    }()
    
    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
    defer cancel()
    
    if err := performTaskWithContext(ctx); err != nil {
        log.Error("Job failed: %v", err)
    }
})
```

### 3. Kết hợp với Singleton Mode

```go
manager.Every(1).Hour().
    SingletonMode().
    Do(func() {
        // Công việc quan trọng, SingletonMode đảm bảo không chạy đồng thời 
        // trên một instance, distributed locker đảm bảo không chạy đồng thời 
        // trên nhiều instances
        performCriticalTask()
    })
```

### 4. Kiểm tra kết nối Redis

```go
// Kiểm tra kết nối Redis trước khi thiết lập locker
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := redisClient.Ping(ctx).Err()
if err != nil {
    log.Fatal("Cannot connect to Redis:", err)
}

locker, err := scheduler.NewRedisLocker(redisClient, options)
```

### 5. Monitor Redis Keys

Bạn có thể theo dõi các khóa Redis để debug và giám sát:

```
KEYS scheduler_lock:*
TTL scheduler_lock:job123
```
