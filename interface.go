package ddd

import (
	"time"

	"github.com/antlinker/ddd/path"
)

// DomainPath 域路径
// 用来唯一定位域子域，领域服务，仓储，聚合根，聚合，实例的路径
type DomainPath interface {
	DomainPath() path.Path
}

// DomainID 领域唯一标识
type DomainID interface {
	// 获取域唯一标识
	DomainID() string
}

// DomainNode 领域节点
type DomainNode interface {
	DomainID
	DomainPath
	Domain() Domain
	Init(self DomainNode, parent DomainNode, DomainID string)
	Parent() DomainNode
	Trigger(c Context, etype, action string, data interface{})
	setParent(parent DomainNode)
	appendChildren(c DomainNode)
	getChildren() (dn map[path.ItemKind]map[string]DomainNode)
	ItemKind() path.ItemKind
	// resetParent(parent DomainNode)
}

// Domain 领域
type Domain interface {
	DomainNode
	SubDomainManager
	AggregateRootManager
	IsRoot() bool
	// 通过服务id获取域服务
	ServiceByID(serviceID string) DomainService
	// 获取域内所有服务
	Services() []DomainService
	// 获取域内所有的实例名
	EntityNames() []string
	EntityByID(id string) Entity
	RepositoryByID(repoID string) Repository
	// 领域全局数据存储
	GlobalStorer() GlobalStorer

	//regRepository(repo Repository)
	// 向领域注册领域服务
	//regService(repo DomainService)
}

// SubDomainManager 子域管理
type SubDomainManager interface {
	ParentDomain() Domain
	// 获取所有子域
	SubDomains() map[string]Domain
	// 获取指定子域
	// 返回值为nil id对应的子域不存在
	SubDomain(id string) Domain
	// 注册子域
	//	regSubDomain(domain Domain)
}

// AggregateRootManager 领域聚合根管理
type AggregateRootManager interface {
	// 通过聚合根唯一标识获取聚合根
	AggregateRootByID(id string) AggregateRoot
	// 获取域内所有聚合根
	AggregateRoots() []AggregateRoot
	// 注册聚合根
	//	regAggregateRoot(ar AggregateRoot)
}

// DomainService 领域服务
type DomainService interface {
	DomainNode
}

// Entity 实例
type Entity interface {
	DomainNode
}

// Aggregate 聚合
type Aggregate interface {
	Entity
}

// AggregateRoot 聚合根
type AggregateRoot interface {
	DomainNode
	// 通过聚合id获取聚合实例
	GetAggregate(aggregateID string) Aggregate
	// 聚合根对应的仓储
	Repository() Repository
	// 设置聚合根对应的仓储，每个聚合根只有唯一的仓储
	setRepository(repo Repository)
}

// Repository 仓储
type Repository interface {
	DomainNode
	GetAggregate(aggregateID string) Aggregate
	AggregateRoot() AggregateRoot
	CacheSet(a DomainNode, expire ...time.Duration) error
	CacheGet(a DomainNode) (bool, error)
	CacheDelete(a DomainNode)
}
