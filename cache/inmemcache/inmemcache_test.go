package inmemcache

import (
	"testing"
	"time"

	"github.com/antlinker/ddd"
)

func TestNewDDDCache(t *testing.T) {
	cache := NewDDDCacher(time.Minute, time.Minute)
	i := "aaaa"
	if err := cache.Set("a", i); err != nil {
		t.Errorf("设置失败:%v", err)
		return
	}
	var b string
	if ok, err := cache.Get("a", &b); err != nil || !ok {
		t.Errorf("读取失败:%v", err)
		return
	}
	if b != i {
		t.Errorf("读取失败:想要一个%v 实际获得 %v", i, b)
		return
	}

}

func TestNewDDDCacheDomainNode(t *testing.T) {
	cache := NewDDDCacher(time.Minute, time.Minute)
	i := &ddd.BaseDomain{}
	i.Init(i, nil, "aaaa")
	if err := cache.Set("a", i); err != nil {
		t.Errorf("设置失败:%v", err)
		return
	}
	var b ddd.Domain
	if ok, err := cache.Get("a", &b); err != nil || !ok {
		t.Errorf("读取失败:%v", err)
		return
	}
	if b.DomainID() != i.DomainID() {
		t.Errorf("读取失败:想要一个%v 实际获得 %v", i.DomainID(), b.DomainID())
		return
	}
	cache.Del("a")
	if ok, _ := cache.Get("a", &b); ok {
		t.Errorf("没有删除成功")
		return
	}

}
