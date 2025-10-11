package internal

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq"
)


func NewStore(pgDSN, redisAddr string) (*Store, error){
    db, err := sql.Open("postgres", pgDSN)
    if err!=nil { return nil, err }
    redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
    ctx,_ := context.WithTimeout(context.Background(),5*time.Second)
    if err := redisClient.Ping(ctx).Err(); err!=nil { return nil, err }
    return &Store{DB: db, RDB: redisClient}, nil
}

// Save order in Redis cache
func (s *Store) CacheOrder(ctx context.Context, o *Order) {
    if s.RDB == nil {
        // Redis not initialized (e.g., in tests), skip caching
        return
    }

    key := fmt.Sprintf("order:%d", o.ID)
    b, _ := json.Marshal(o)
    s.RDB.Set(ctx, key, b, time.Hour)
}

func (s *Store) GetCachedOrder(ctx context.Context, id int) (*Order, error) {
    if s.RDB == nil {
        return nil, fmt.Errorf("Redis not initialized")
    }

    key := fmt.Sprintf("order:%d", id)
    val, err := s.RDB.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }

    var o Order
    json.Unmarshal([]byte(val), &o)
    return &o, nil
}

