
package internal
import (
    "database/sql"
    "github.com/go-redis/redis/v8"
)

type Store struct {
    DB    *sql.DB
    Cache map[int]*Order
    RDB   *redis.Client
}
func (s *Store) Close() error {
    if s.DB != nil {
        return s.DB.Close()
    }
    return nil
}