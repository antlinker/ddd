package ddd

import (
	"fmt"
	"time"

	"github.com/antlinker/ddd/path"
)

// type Repository interface {
// 	DomainID
// 	DomainPath
// 	ID() string
// }

var (
	_ Repository = &BaseRepository{}
)

// BaseRepository 基础仓储
type BaseRepository struct {
	_repository
}

type _repository struct {
	node
	aggRoot AggregateRoot
}

func (r *_repository) Init(self DomainNode, parent DomainNode, domainID string) {
	r.init(self, parent, path.Repository, domainID, false)
	// switch n := self.(type) {
	// case Domain:
	// 	n.regRepository(self.(Repository))
	// case AggregateRoot:
	// 	n.setRepository(self.(Repository))
	// }
}

func (r _repository) DomainPath() path.Path {
	return r.path

}
func (r _repository) GetAggregate(aggregateID string) Aggregate {
	panic("需要每一个仓储实现该方法")
}

// AggregateRoot 获取上级聚合根
func (r *_repository) AggregateRoot() AggregateRoot {
	if r.aggRoot != nil {
		return r.aggRoot
	}
	for p := r.Parent(); p != nil; p = p.Parent() {
		switch d := p.(type) {
		case AggregateRoot:
			r.aggRoot = d
			return d
		}
	}
	return nil
}

// CacheSet 缓存聚合
// 参数a 必须指定DomainID()
// expire 有效时间，单位秒
func (r _repository) CacheSet(a DomainNode, expire ...time.Duration) error {
	c := Cache()
	if c == nil {
		return nil
	}
	key := r.key(a)
	if len(expire) > 0 && expire[0] > 0 {
		return Cache().SetExpire(key, a, expire[0])
	}
	return c.Set(key, a)
}

// CacheGet 获取缓存聚合
// a 参数必须指定DomainID
func (r _repository) CacheGet(a DomainNode) (bool, error) {
	c := Cache()
	if c == nil {
		return false, nil
	}
	key := r.key(a)
	return c.Get(key, a)
}

// CacheDelete 删除缓存聚合
// a 参数必须指定DomainID
func (r _repository) CacheDelete(a DomainNode) {
	c := Cache()
	if c == nil {
		return
	}
	key := r.key(a)
	c.Del(key)
}
func (r _repository) key(a DomainNode) string {
	if a.Parent() != nil {
		// 已经设置上级
		return a.DomainPath().Path()
	}
	switch a.(type) {
	case Aggregate:
		ar := r.AggregateRoot()
		if ar == nil {
			panic(fmt.Sprintf("仓储 %v 未初始化，或不是一个聚合根的仓储", r))
		}
		a.Init(a, ar, a.DomainID())
		return a.DomainPath().Path()
	case Entity:
		d := r.Domain()
		if d == nil {
			// 该情况说明该repo 还没有设置上级
			panic(fmt.Sprintf("仓储 %v 未初始化", r))
		}
		e := path.NewItem(path.Entity, a.DomainID())
		d.DomainPath().Append(e)
		return e.Path()
	}
	return ""
}
