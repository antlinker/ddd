package mgorepo

import (
	"time"

	"github.com/globalsign/mgo"

	"github.com/antlinker/ddd"
)

// Coll 实例对应集合接口
type Coll interface {
	CollName() string
	CollID() string
}

// CollAggregate 聚合对应集合接口
type CollAggregate interface {
	Coll
	ToAggregate() ddd.Aggregate
}

// CollAggregateCache 缓存设置
type CollAggregateCache interface {
	// 是否可以缓存
	IsCache() bool
}

// CollAggregateCacheExpire 设置缓存时长
type CollAggregateCacheExpire interface {
	CacheExpire() time.Duration
}

// CollIndex 索引
type CollIndex interface {
	Indexes() []mgo.Index
}
