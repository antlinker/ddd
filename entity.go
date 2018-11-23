package ddd

import (
	"github.com/antlinker/ddd/path"
)

var (
	_ Entity = &BaseEntity{}
)

// BaseEntity 基础实例
type BaseEntity struct {
	_entity
}

type _entity struct {
	node
}

func (r *_entity) Init(self DomainNode, parent DomainNode, domainID string) {
	r.init(self, parent, path.Entity, domainID, true)
}
