package internal

import (
	"github.com/garyburd/redigo/redis"
)

// 删除给定的一个或多个 key 。
// 不存在的 key 会被忽略。
// 可用版本：
// >= 1.0.0
// 返回值：
// 被删除 key 的数量。
func (f *FRedis) Del(key string) (int, error) {
	if f.conn_nil() {
		return 0, nil
	}
	c := f.conn.Get()
	defer c.Close()
	return redis.Int(c.Do("DEL", key))
}

// 序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
// 序列化生成的值有以下几个特点：
// 它带有 64 位的校验和，用于检测错误， RESTORE 在进行反序列化之前会先检查校验和。
// 值的编码格式和 RDB 文件保持一致。
// RDB 版本会被编码在序列化值当中，如果因为 Redis 的版本不同造成 RDB 格式不兼容，那么 Redis 会拒绝对这个值进行反序列化操作。
// 序列化的值不包括任何生存时间信息。
// 可用版本：
// >= 2.6.0
// 返回值：
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (f *FRedis) Dump(key string) (string, error) {
	if f.conn_nil() {
		return "", nil
	}
	c := f.conn.Get()
	defer c.Close()
	return redis.String(c.Do("DUMP", key))
}

//检查给定 key 是否存在。
// 返回值:
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (f *FRedis) Exists(key string) bool {
	if f.conn_nil() {
		return false
	}
	c := f.conn.Get()
	defer c.Close()
	res, err := redis.Int(c.Do("EXISTS", key))
	if err != nil {
		f.error_info("EXISTS", key, err)
		return false
	}
	return res == 1
}

// 返回值:
// 设置成功返回 1 。
// 当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
func (f *FRedis) Expire(key string, seconds int) bool {
	if f.conn_nil() {
		return false
	}
	c := f.conn.Get()
	defer c.Close()
	_, err := redis.Int(c.Do("EXPIRE", key, seconds))
	if err != nil {
		f.error_info("EXPIRE", key, err)
		return false
	}
	return true
}

// EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 不同在于 EXPIREAT 命令接受的时间参数是 UNIX 时间戳(unix timestamp)。
// 可用版本：
// >= 1.2.0
// 返回值：
// 如果生存时间设置成功，返回 1 。
// 当 key 不存在或没办法设置生存时间，返回 0
func (f *FRedis) Expireat(key string, timestamp int) bool {
	if f.conn_nil() {
		return false
	}
	c := f.conn.Get()
	defer c.Close()
	res, err := redis.Int(c.Do("EXPIREAT", key, timestamp))
	if err != nil {
		f.error_info("EXPIRE", key, err)
		return false
	}
	return res == 1
}
