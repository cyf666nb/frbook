# 高并发 Go 编写注意事项

## 一、数据库连接池配置

### 1.1 MySQL 连接池

```go
sqlDB.SetMaxOpenConns(100)     // 最大打开连接数
sqlDB.SetMaxIdleConns(10)      // 最大空闲连接数
sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最大生命周期
sqlDB.SetConnMaxIdleTime(time.Minute * 10)  // 空闲连接最大存活时间
```

**注意事项:**
- `MaxOpenConns` 应根据数据库服务器配置和并发量设置,一般不超过数据库最大连接数的 80%
- `MaxIdleConns` 建议设置为 `MaxOpenConns` 的 10%-25%
- `ConnMaxLifetime` 设置连接最大生命周期,避免长时间连接导致的资源泄漏
- `ConnMaxIdleTime` 设置空闲连接超时,释放不活跃的连接

### 1.2 Redis 连接池

```go
client := redis.NewClient(&redis.Options{
    PoolSize:     100,              // 连接池大小
    MinIdleConns: 10,               // 最小空闲连接数
    MaxRetries:   3,                // 最大重试次数
    DialTimeout:  5 * time.Second,  // 连接超时
    ReadTimeout:  3 * time.Second,  // 读超时
    WriteTimeout: 3 * time.Second,  // 写超时
    PoolTimeout:  4 * time.Second,  // 获取连接超时
})
```

**注意事项:**
- `PoolSize` 建议设置为 `CPU核心数 * 10`
- `MinIdleConns` 保持一定数量的空闲连接,避免频繁创建
- 设置合理的超时时间,防止慢查询阻塞连接池

---

## 二、Context 超时控制

### 2.1 请求级别超时

```go
func (h *Handler) GetBook(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    book, err := h.service.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            response.Error(c, 504, "请求超时")
            return
        }
        response.Error(c, 500, err.Error())
        return
    }
    response.Success(c, book)
}
```

### 2.2 数据库操作超时

```go
func (r *Repository) FindByID(ctx context.Context, id uint64) (*model.Book, error) {
    var book model.Book
    err := r.db.WithContext(ctx).First(&book, id).Error
    if err != nil {
        return nil, err
    }
    return &book, nil
}
```

### 2.3 Redis 操作超时

```go
func (r *RedisClient) GetWithTimeout(ctx context.Context, key string) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    return r.Get(ctx, key).Result()
}
```

**注意事项:**
- 所有 I/O 操作都应传入 context
- 设置合理的超时时间,避免级联超时
- 使用 `defer cancel()` 释放资源

---

## 三、分布式锁实现

### 3.1 基于 Redis 的分布式锁

```go
func (r *RedisClient) Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
    return r.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisClient) Unlock(ctx context.Context, key string, value string) error {
    script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `
    return r.Eval(ctx, script, []string{key}, value).Err()
}
```

### 3.2 使用示例

```go
func (s *Service) CreateOrder(ctx context.Context, bookID uint64) error {
    lockKey := fmt.Sprintf("lock:book:%d", bookID)
    lockValue := uuid.New().String()
    
    locked, err := s.redis.Lock(ctx, lockKey, lockValue, 10*time.Second)
    if err != nil || !locked {
        return errors.New("获取锁失败")
    }
    defer s.redis.Unlock(ctx, lockKey, lockValue)
    
    // 业务逻辑
    return s.doCreateOrder(ctx, bookID)
}
```

**注意事项:**
- 锁必须设置过期时间,防止死锁
- 解锁时使用 Lua 脚本保证原子性
- 锁的 value 应为唯一值,防止误删其他请求的锁
- 业务处理时间不应超过锁的过期时间

---

## 四、并发安全

### 4.1 避免全局变量

```go
// 错误示例
var userCache = make(map[uint64]*User)

// 正确示例
type UserService struct {
    cache sync.Map  // 或使用 sync.RWMutex 保护 map
}
```

### 4.2 使用 sync.Map

```go
type Cache struct {
    data sync.Map
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, value)
}

func (c *Cache) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}
```

### 4.3 使用 sync.RWMutex

```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (m *SafeMap) Get(key string) interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.data[key]
}

func (m *SafeMap) Set(key string, value interface{}) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
}
```

---

## 五、Goroutine 管理

### 5.1 使用 errgroup 管理并发

```go
import "golang.org/x/sync/errgroup"

func (s *Service) GetBookDetails(ctx context.Context, ids []uint64) ([]*model.BookDetail, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]*model.BookDetail, len(ids))
    
    for i, id := range ids {
        i, id := i, id  // 捕获循环变量
        g.Go(func() error {
            detail, err := s.GetDetail(ctx, id)
            if err != nil {
                return err
            }
            results[i] = detail
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### 5.2 防止 Goroutine 泄漏

```go
func (s *Service) ProcessOrder(ctx context.Context, orderID uint64) {
    go func() {
        select {
        case <-ctx.Done():
            return  // 上下文取消时退出
        default:
            s.doProcess(orderID)
        }
    }()
}
```

### 5.3 使用 Worker Pool

```go
type WorkerPool struct {
    tasks chan func()
    wg    sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
    p := &WorkerPool{
        tasks: make(chan func(), 100),
    }
    for i := 0; i < size; i++ {
        p.wg.Add(1)
        go p.worker()
    }
    return p
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for task := range p.tasks {
        task()
    }
}

func (p *WorkerPool) Submit(task func()) {
    p.tasks <- task
}

func (p *WorkerPool) Shutdown() {
    close(p.tasks)
    p.wg.Wait()
}
```

---

## 六、内存优化

### 6.1 对象复用

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    buf.Write(data)
    return buf.String()
}
```

### 6.2 避免不必要的内存分配

```go
// 错误示例
func Concat(strs []string) string {
    var result string
    for _, s := range strs {
        result += s  // 每次都创建新字符串
    }
    return result
}

// 正确示例
func Concat(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}
```

### 6.3 预分配切片容量

```go
// 错误示例
func Convert(items []Item) []Result {
    var results []Result
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}

// 正确示例
func Convert(items []Item) []Result {
    results := make([]Result, 0, len(items))
    for _, item := range items {
        results = append(results, process(item))
    }
    return results
}
```

---

## 七、错误处理

### 7.1 错误包装

```go
import "errors"

func (s *Service) GetBook(ctx context.Context, id uint64) (*model.Book, error) {
    book, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find book %d: %w", id, err)
    }
    return book, nil
}
```

### 7.2 错误判断

```go
import "errors"

if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrBookNotFound
}

var customErr *CustomError
if errors.As(err, &customErr) {
    return customErr.Code
}
```

---

## 八、性能监控

### 8.1 pprof 性能分析

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    // ...
}
```

### 8.2 指标采集

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(requestDuration)
}
```

---

## 九、限流与熔断

### 9.1 令牌桶限流

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(rps), rps),
    }
}

func (r *RateLimiter) Allow() bool {
    return r.limiter.Allow()
}
```

### 9.2 熔断器

```go
import "github.com/afex/hystrix-go/hystrix"

func init() {
    hystrix.ConfigureCommand("get_book", hystrix.CommandConfig{
        Timeout:               3000,
        MaxConcurrentRequests: 100,
        ErrorPercentThreshold: 50,
    })
}

func (s *Service) GetBook(id uint64) (*model.Book, error) {
    var book *model.Book
    err := hystrix.Do("get_book", func() error {
        var err error
        book, err = s.repo.FindByID(context.Background(), id)
        return err
    }, func(err error) error {
        return s.getFromCache(id)
    })
    return book, err
}
```

---

## 十、最佳实践总结

1. **连接池管理**: 合理配置数据库和 Redis 连接池参数
2. **超时控制**: 所有 I/O 操作都设置超时,使用 context 传递
3. **分布式锁**: 使用 Redis SETNX + Lua 脚本实现原子性操作
4. **并发安全**: 避免共享状态,使用 channel 或 sync 包进行同步
5. **资源管理**: 使用 defer 确保资源释放,避免泄漏
6. **错误处理**: 使用错误包装和判断,提供清晰的错误信息
7. **性能监控**: 集成 pprof 和 Prometheus 进行性能监控
8. **限流熔断**: 实现限流和熔断机制,保护系统稳定性
9. **优雅关闭**: 实现优雅关闭,确保请求处理完成
10. **代码规范**: 遵循 Go 代码规范,保持代码可读性和可维护性
