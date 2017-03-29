package internal

import (
	"fmt"

	"fantasy/iconfig"
	"github.com/garyburd/redigo/redis"
)

type IFRedis interface {
	Open(iconfig.IConfig, string) error
	// OpenWithConfig(config RedisConfig) error
}

type FRedis struct {
	conn *redis.Pool
	Logger
}

func (f *FRedis) Open(c iconfig.IConfig, name string) error {
	ip, err := c.String(name, "ip")
	if err != nil {
		return err
	}
	port, err := c.String(name, "port")
	if err != nil {
		return err
	}
	maxIdle, err := c.Int(name, "maxIdle")
	if err != nil {
		return err
	}
	rpool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", fmt.Sprintf("%v:%v", ip, port))
	}, maxIdle)
	f.conn = rpool
	return nil
}

func (f *FRedis) error_info(prefix, key string, err error) {
	f.Errorf("ERROR: %v key:%v err:%v", prefix, key, err)
}

func (f *FRedis) conn_nil() bool {
	if f.conn == nil {
		f.SetPrefix("redis: ")
		f.Error("redis conn is nil")
		return true
	}
	return false
}
