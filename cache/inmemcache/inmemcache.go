package inmemcache

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/antlinker/ddd"

	cache "github.com/patrickmn/go-cache"
)

// NewDDDCacher 创建缓存
func NewDDDCacher(defaultExpiration, cleanupInterval time.Duration) ddd.Cacher {
	return &dddCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

type dddCache struct {
	cache *cache.Cache
}

// Get Value
// 没有缓存 ok =false
// err !=nil 获取出错
func (c dddCache) Get(key string, dst interface{}) (ok bool, err error) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			err = fmt.Errorf("赋值失败:%v", err)
			return
		}
	}()
	if data, ok := c.cache.Get(key); ok {
		rv := reflect.ValueOf(dst)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return false, errors.New("dst类型不正确")
		}
		rv.Elem().Set(reflect.ValueOf(data))
		return true, nil

	}
	return false, nil

}

// Set Value
func (c dddCache) Set(key string, v interface{}) error {
	c.cache.Set(key, v, cache.NoExpiration)
	return nil
}

// Delete Cache
func (c dddCache) Del(key string) {
	c.cache.Delete(key)
}

// Auto Delete Set
func (c dddCache) SetExpire(key string, v interface{}, d time.Duration) error {
	c.cache.Set(key, v, d)
	return nil
}
