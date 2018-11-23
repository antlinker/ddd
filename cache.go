package ddd

import (
	"time"
)

var defaultCacher Cacher

// RegCacher 注册缓存器
func RegCacher(c Cacher) {
	defaultCacher = c
}

// Cache 缓存实现
func Cache() Cacher {
	return defaultCacher
}

// Cacher 缓存
type Cacher interface {
	// Get Value
	// 没有缓存 ok =false
	// err !=nil 获取出错
	Get(key string, dst interface{}) (ok bool, err error)
	// Set Value
	Set(key string, v interface{}) error
	// Delete Cache
	Del(key string)
	// Auto Delete Set
	SetExpire(key string, v interface{}, d time.Duration) error
}
