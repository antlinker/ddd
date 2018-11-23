package ddd

import "github.com/antlinker/ddd/path"

var (
	_ AggregateRoot = &BaseAggregateRoot{}
)

// BaseAggregateRoot 基础聚合根
type BaseAggregateRoot struct {
	_aggregateRoot
}

type _aggregateRoot struct {
	node
	repo Repository
}

func (r *_aggregateRoot) Init(self DomainNode, parent DomainNode, domainID string) {
	r.init(self, parent, path.AggregateRoot, domainID, false)
	// d := parent.(Domain)
	// d.regAggregateRoot(self.(AggregateRoot))
}

// 通过聚合名和id获取
func (r _aggregateRoot) GetAggregate(aggregateID string) Aggregate {
	return r.repo.GetAggregate(aggregateID)
}

// 聚合根对应的仓储
func (r _aggregateRoot) Repository() Repository {
	if s, ok := r.getNodes(path.Repository); ok {
		for _, n := range s {
			if nr, o := n.(Repository); o {
				return nr
			}
		}
	}
	return nil
}
func (r *_aggregateRoot) setRepository(repo Repository) {
	r.repo = repo
}

// func (r *_aggregateRoot) SetParent(parent DomainNode) {
// 	r.node.resetParent(parent)
// 	if r.repo != nil {
// 		r.repo.resetParent(r)
// 	}
// }
