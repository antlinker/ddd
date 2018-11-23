package path

import (
	"fmt"

	"github.com/pkg/errors"
)

// NewDomainPath 创建一个域路径
func NewDomainPath(id string) Path {
	di := &domainItem{
		pathItem: pathItem{
			curName: id,
			kind:    Domain,
		},
		curpath: id,
	}
	_ = di.Path()
	p := &_path{
		root: di,
	}
	_ = p.Path()
	return p
}
func parseDomain(r string) (*domainItem, *domainItem, string, error) {
	domains, r, ok := parseNamespaces(r, domainPre, domainAliasLen)
	if !ok {
		return nil, nil, "", errors.Errorf("路径(%v)尝试解析领域时失败。", r)
	}
	var first *domainItem
	var curItem *domainItem
	for i, v := range domains {
		if i == 0 {
			curItem = &domainItem{
				pathItem: pathItem{
					curName: v,
					kind:    Domain,
				},
				curpath: v,
			}
			_ = curItem.Path()
			first = curItem
		} else {
			c := &domainItem{
				curpath: curItem.curpath + string(pathSpe) + v,
				pathItem: pathItem{
					curName: v,
					kind:    Domain,
					parent:  curItem,
				},
			}
			curItem.next = c
			curItem = c
		}

	}
	return first, curItem, r, nil
}

// creDomainItem 创建新的domain item
func creDomainItem(id string) *domainItem {
	return &domainItem{
		curpath: id,
		pathItem: pathItem{
			curName: id,
			kind:    Domain,
		},
	}
}

type domainItem struct {
	pathItem
	curpath string
	path    string
}

func (i *domainItem) Path() string {
	if i.path != "" {
		return i.path
	}
	i.path = fmt.Sprintf(fmtDomainPath, i.curpath)
	return i.path
}
func (i *domainItem) Append(in Item) {
	i.next = in
	in.setParent(i)
	_ = in.resetPath()
}
func (i *domainItem) resetPath() string {
	i.path = fmt.Sprintf(fmtDomainPath, i.curpath)
	return i.path
}
func (i *domainItem) setParent(in Item) {
	i.parent = in
	din, ok := in.(*domainItem)
	if ok {
		i.curpath = din.curpath + string(pathSpe) + i.curName
	}
}
func (i *domainItem) parse(r string) error {

	// 解析聚合根
	ars, r, err := parseAggregateRoot(r)
	if err != nil {
		return err
	}

	if ars != nil {
		ars.parent = i
		i.next = ars
		ars.Path()
		if r != "" {
			return ars.parse(r)
		}
		return nil
	}

	// 尝试解析仓储
	repo, err := parseRepository(r)
	if err != nil {
		return err
	}
	if repo != nil {
		i.next = repo
		repo.parent = i
		i.next.Path()
		return nil
	}

	// 尝试解析服务
	service, err := parseService(r)
	if err != nil {
		return err
	}
	if service != nil {
		i.next = service
		service.parent = i
		_ = i.next.Path()
		return nil
	}
	// 尝试解析实例
	entiry, err := parseEntity(r)
	if err != nil {
		return err
	}
	if entiry != nil {
		i.next = entiry
		entiry.parent = i
		i.next.Path()
		return nil
	}
	return errors.Errorf("解析领域失败，有不能解析的路径(%v)存在", r)

}
