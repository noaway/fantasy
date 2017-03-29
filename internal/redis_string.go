package internal

import (
	"github.com/garyburd/redigo/redis"
)

func APPEND() {

}

func (f *FRedis) Get(key string) (string, error) {
	if f.conn_nil() {
		return "", nil
	}
	c := f.conn.Get()
	defer c.Close()
	return redis.String(c.Do("GET", key))
}
