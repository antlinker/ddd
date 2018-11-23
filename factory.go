package ddd

import (
	"fmt"

	"github.com/antlinker/ddd/path"
	"github.com/pkg/errors"
)

// DomainManager 领域管理者
// 通过领域管理者可以查询所有的领域对象
type DomainManager interface {
	GetDomain(domainid string) Domain
	RegDomain(d Domain)
	FindNode(path string) DomainNode
	FindDomain(path string) Domain
	FindAggregateRoot(path string) AggregateRoot
	FindAggregate(path string) Aggregate
	FindDomainService(path string) DomainService
	FindEntity(path string) Entity
	FindRepository(path string) Repository
}

type domainManage struct {
	domains map[string]Domain
}

func (m *domainManage) RegDomain(d Domain) {
	_, ok := m.domains[d.DomainID()]
	if ok {
		panic(fmt.Sprintf("已经注册了一个同名的领域 ：%v", d.DomainID()))
	}
	m.domains[d.DomainID()] = d
}

func (m *domainManage) GetDomain(id string) Domain {
	return m.domains[id]
}
func (m *domainManage) FindNode(pstr string) (DomainNode, error) {
	p := path.FromString(pstr)
	if p.IsInvalid() {
		return nil, errors.Errorf("错误的域路径%v", pstr)
	}
	i := p.Next()

	if i.Kind() != path.Domain {
		return nil, errors.Errorf("领域路径必须从领域开始")
	}

	d := m.GetDomain(i.CurName())
	if d == nil {
		return nil, ErrorNoFoundPath
	}
	var cur DomainNode = d
	for i = p.Next(); i != nil; i = p.Next() {
		switch i.Kind() {
		case path.Domain:
			switch cp := cur.(type) {
			case Domain:
				cur = cp.SubDomain(i.CurName())
			default:
				return nil, ErrorNoFoundPath
			}
		case path.AggregateRoot:
			switch cp := cur.(type) {
			case Domain:
				cur = cp.AggregateRootByID(i.CurName())
			default:
				return nil, ErrorNoFoundPath
			}
		case path.Aggregate:
			switch cp := cur.(type) {
			case AggregateRoot:
				cp.GetAggregate(i.CurName())
			default:
				return nil, ErrorNoFoundPath
			}
		case path.Service:
			switch cp := cur.(type) {
			case Domain:
				cur = cp.ServiceByID(i.CurName())
			default:
				return nil, ErrorNoFoundPath
			}
		case path.Repository:
			switch cp := cur.(type) {
			case Domain:
				cur = cp.RepositoryByID(i.CurName())
			case AggregateRoot:
				cur = cp.Repository()
			default:
				return nil, ErrorNoFoundPath
			}
		case path.Entity:
			switch cp := cur.(type) {
			case Domain:
				cur = cp.EntityByID(i.CurName())
			default:
				return nil, ErrorNoFoundPath
			}
		}
	}

	return cur, nil

}

func (m *domainManage) FindDomain(pstr string) (Domain, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(Domain); ok {
		return d, nil
	}
	return nil, ErrorNodeKindNotMatch

}
func (m *domainManage) FindAggregateRoot(pstr string) (AggregateRoot, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(AggregateRoot); ok {
		return d, nil
	}

	return nil, ErrorNodeKindNotMatch
}
func (m *domainManage) FindAggregate(pstr string) (Aggregate, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(Aggregate); ok {
		return d, nil
	}

	return nil, ErrorNodeKindNotMatch
}
func (m *domainManage) FindDomainService(pstr string) (DomainService, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(DomainService); ok {
		return d, nil
	}

	return nil, ErrorNodeKindNotMatch
}
func (m *domainManage) FindEntity(pstr string) (Entity, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(Entity); ok {
		return d, nil
	}

	return nil, ErrorNodeKindNotMatch
}
func (m *domainManage) FindRepository(pstr string) (Repository, error) {
	if n, err := m.FindNode(pstr); err != nil {
		return nil, err
	} else if d, ok := n.(Repository); ok {
		return d, nil
	}
	return nil, ErrorNodeKindNotMatch
}
