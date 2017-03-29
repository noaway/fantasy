package internal

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"fantasy/iconfig"
	"github.com/noaway/config"
)

var DefaultConfigManager *ConfigManager = NewConfigManager()

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		rediss: make(map[string]IFRedis),
	}
}

type ConfigManager struct {
	sync.Mutex
	config iconfig.IConfig
	rediss map[string]IFRedis
}

func (c *ConfigManager) RegisterRedis(redis IFRedis) {
	c.Lock()
	defer c.Unlock()

	name := getTypeMoudleName(redis)
	if redis == nil {
		panic("Register redis driver is nil")
	}
	if _, dup := c.rediss[name]; dup {
		panic("Register called twice for redis driver " + name)
	}
	c.rediss[name] = redis
	fmt.Println("register redis with", name)
}

func (c *ConfigManager) OpenConfig(name string) error {
	cfg, err := config.ReadDefault(name)
	if err != nil {
		return err
	}
	c.config = cfg

	c.Lock()
	defer c.Unlock()
	//初始化日志
	NewLogger(cfg, "")

	// open registered redis
	for name, redis := range c.rediss {
		fmt.Println("::::", name)
		if err = redis.Open(c.config, name); err != nil {
			fmt.Println("open redis", name, "failed with", err)
			panic(err)
		}
		fmt.Println("open redis:", name)
	}
	return nil
}

func getTypeMoudleName(v interface{}) string {
	var path string
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Struct:
		path = value.Type().PkgPath()
	case reflect.Ptr:
		path = value.Elem().Type().PkgPath()
	default:
		panic("err type")
	}
	moudle := strings.Split(path, "/")

	if len(moudle) > 0 {
		if moudle[len(moudle)-1] == "common" {
			panic(path)
		}
		return moudle[len(moudle)-1]
	}
	return path
}
