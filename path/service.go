package path

import (
	"github.com/pkg/errors"
)

func parseService(r string) (*pathItem, error) {
	// 尝试解析服务
	service, r, ok := parseNamespaces(r, servicePre, serviceAliasLen)
	if !ok {
		return nil, errors.Errorf("路径(%v)尝试解析服务时失败.", r)
	}
	if len(service) == 1 && r != "" {
		return nil, errors.Errorf("路径(%v)尝试解析服务时失败,有多余的字符存在", r)
	}
	if len(service) == 1 {
		return creServiceItem(service[0]), nil
	} else if len(service) > 1 {
		return nil, errors.Errorf("路径(%v)尝试解析服务时失败,服务不能分级", r)
	}
	return nil, nil
}

// creServiceItem 创建服务Item
func creServiceItem(id string) *pathItem {
	return &pathItem{
		curName: id,
		kind:    Service,
		fmtPath: fmtServicePath,
	}
}

// type serviceItem struct {
// 	pathItem
// }

// func (i *serviceItem) Append(in Item) {
// 	i.next = in
// 	in.setParent(i)
// 	_ = in.resetPath()
// }
