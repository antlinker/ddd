package ddd

import "github.com/antlinker/ddd/path"

var (
	_ Aggregate = &BaseAggregate{}
)

// BaseAggregate 基础聚合
type BaseAggregate struct {
	_aggregate
}
type _aggregate struct {
	_entity
}

func (r *_aggregate) Init(self DomainNode, parent DomainNode, domainID string) {
	r.init(self, parent, path.Aggregate, domainID, true)
}
func (r *_aggregate) Repository() Repository {
	if p := r.Parent(); p != nil {
		if pa, ok := p.(AggregateRoot); ok {
			return pa.Repository()
		}
	}
	return nil
}
