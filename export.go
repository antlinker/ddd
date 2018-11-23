package ddd

import (
	"github.com/antlinker/ddd/log"
)

var (
	defaultDM = &domainManage{domains: make(map[string]Domain)}
)

// GetDomain 获取指定的领域
func GetDomain(domainid string) Domain {
	return defaultDM.GetDomain(domainid)
}

// RegDomain 注册领域
func RegDomain(d Domain) {
	defaultDM.RegDomain(d)
}

// // RegSubDomain 注册子域
// func RegSubDomain(d Domain, sub Domain) {
// 	d.regSubDomain(sub)
// }

// // RegAggregateRoot 向领域或者子域内注册聚合跟
// func RegAggregateRoot(d Domain, ar AggregateRoot) {
// 	d.regAggregateRoot(ar)
// }

// // SetRepoForARoot 为聚合根设置对应的仓储
// // 每个聚合根只能设置一个仓储
// func SetRepoForARoot(ar AggregateRoot, repo Repository) {
// 	ar.setRepository(repo)
// }

// // RegRepositoryByDomain 向领域注入仓储
// func RegRepositoryByDomain(d Domain, repo Repository) {
// 	d.regRepository(repo)
// }

// // RegService 向领域注入仓储
// func RegService(d Domain, svc DomainService) {
// 	d.regService(svc)
// }

// FindNode 查找领域对象节点
func FindNode(path string) (DomainNode, error) {
	return defaultDM.FindNode(path)
}

// FindDomain 通过路径查找领域节点
func FindDomain(path string) (Domain, error) {
	return defaultDM.FindDomain(path)
}

// FindAggregateRoot 通过路径查找聚合根
func FindAggregateRoot(path string) (AggregateRoot, error) {
	return defaultDM.FindAggregateRoot(path)
}

// FindAggregate 通过路径查找聚合
func FindAggregate(path string) (Aggregate, error) {
	return defaultDM.FindAggregate(path)
}

// FindDomainService 通过路径查找服务
func FindDomainService(path string) (DomainService, error) {
	return defaultDM.FindDomainService(path)
}

// FindEntity 通过路径查找实例
func FindEntity(path string) (Entity, error) {
	return defaultDM.FindEntity(path)
}

// FindRepository 通过路径查找仓储
func FindRepository(path string) (Repository, error) {
	return defaultDM.FindRepository(path)
}

// PrintDomain 输出当前领域内的对象
func PrintDomain(d DomainNode) {
	printDomain("", d)
}
func printDomain(pre string, d DomainNode) {
	log.Infof("%s%s->%s", pre, d.ItemKind().Name(), d.DomainID())
	children := d.getChildren()
	if children != nil {
		for _, v := range children {
			for _, v := range v {
				printDomain(pre+"    ", v)
			}
		}
	}
}
