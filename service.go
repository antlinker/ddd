package ddd

import "github.com/antlinker/ddd/path"

// type Repository interface {
// 	DomainID
// 	DomainPath
// 	ID() string
// }

var (
	_ DomainService = &BaseService{}
)

// BaseService 基础领域服务
type BaseService struct {
	_service
}

type _service struct {
	node
	id string
}

func (r *_service) Init(self DomainNode, parent DomainNode, domainID string) {
	r.init(self, parent, path.Service, domainID, false)
	// d := parent.(Domain)
	// d.regService(self.(DomainService))
}

func (r _service) DomainPath() path.Path {
	return r.path

}
