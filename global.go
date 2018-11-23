package ddd

import "sync"

// GlobalStorer 领域全局存储
type GlobalStorer interface {
	Put(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{})
}

type globalStore struct {
	sync.Map
}

// 领域全局数据存储
func (d *globalStore) Put(key interface{}, value interface{}) {
	d.Store(key, value)
}

// 领域全局数据 获取
func (d *globalStore) Get(key interface{}) interface{} {
	key, ok := d.Load(key)
	if ok {
		return key
	}
	return nil
}

// 领域全局数据 删除
func (d *globalStore) Remove(key interface{}) {
	d.Delete(key)
}
