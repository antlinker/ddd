package ddd

import (
	"github.com/antlinker/ddd/path"
)

var (
	_ Domain = &BaseDomain{}
)

// BaseDomain 基础领域实现，其他领域实现需要集成该结构体
type BaseDomain struct {
	domain
}

// domain 实现Domain接口

type domain struct {
	node
	globalStore globalStore
	// children       map[string]Domain
	// services       map[string]DomainService
	// repositorys    map[string]Repository
	// aggregateRoots map[string]AggregateRoot
	//parent Domain
}

func (d *domain) Init(self DomainNode, parent DomainNode, domainID string) {
	if parent == nil {
		d.path = path.NewDomainPath(domainID)
	}
	d.init(self, parent, path.Domain, domainID, false)
	// if d.aggregateRoots == nil {
	// 	d.aggregateRoots = make(map[string]AggregateRoot)
	// }
	// d.path = path.NewDomainPath(domainID)
	// d.pathItem = d.path.Next()
	// d.domainID = domainID
	// if parent != nil {
	// 	dp, ok := parent.(Domain)
	// 	if !ok {
	// 		panic(fmt.Sprintf("初始化领域%s失败，上级节点不是一个有效的领域节点", domainID))
	// 	}
	// 	d.parent = dp
	// 	d.resetParent(parent)
	// }
	// d.regSubDomain(self.(Domain))
}
func (d *domain) GlobalStorer() GlobalStorer {
	return &d.globalStore
}
func (d *domain) IsRoot() bool {
	return d.parent == nil
}

// 通过服务id获取域服务
func (d *domain) ServiceByID(serviceID string) DomainService {
	if s, ok := d.getNode(path.Service, serviceID); ok {
		if ds, o := s.(DomainService); o {
			return ds
		}
	}
	return nil
}

// 获取域内所有服务
func (d *domain) Services() (out []DomainService) {
	if m, ok := d.getNodes(path.Service); ok {
		for _, s := range m {
			if ds, o := s.(DomainService); o {
				out = append(out, ds)
			}
		}
	}
	return
}

// 获取域内所有的实例名
func (d *domain) EntityNames() []string {
	// TODO 未实行
	return nil
}
func (d *domain) EntityByID(id string) Entity {
	// TOTO 未实现
	return nil
}
func (d *domain) RepositoryByID(repoID string) Repository {
	if s, ok := d.getNode(path.Repository, repoID); ok {
		if ds, o := s.(Repository); o {
			return ds
		}
	}
	return nil
}

// 通过聚合根唯一标识获取聚合根
func (d *domain) AggregateRootByID(id string) AggregateRoot {
	if s, ok := d.getNode(path.AggregateRoot, id); ok {
		if ds, o := s.(AggregateRoot); o {
			return ds
		}
	}
	return nil
}

// 获取域内所有聚合根
func (d *domain) AggregateRoots() (out []AggregateRoot) {
	if m, ok := d.getNodes(path.AggregateRoot); ok {
		for _, s := range m {
			if ds, o := s.(AggregateRoot); o {
				out = append(out, ds)
			}
		}
	}
	return
}

func (d *domain) ParentDomain() Domain {
	if d.parent != nil {
		return d.parent.(Domain)
	}
	return nil
}

func (d *domain) SubDomains() (out map[string]Domain) {
	out = make(map[string]Domain)
	if m, ok := d.getNodes(path.Domain); ok {
		for _, s := range m {
			if ds, o := s.(Domain); o {
				out[ds.DomainID()] = ds
			}
		}
	}
	return
}

// 获取指定子域
func (d *domain) SubDomain(id string) Domain {
	if s, ok := d.getNode(path.Domain, id); ok {
		if ds, o := s.(Domain); o {
			return ds
		}
	}
	return nil
}
