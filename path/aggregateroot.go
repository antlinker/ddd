package path

import (
	"github.com/pkg/errors"
)

func parseAggregateRoot(r string) (*aggregateRootItem, string, error) {
	ars, r, ok := parseNamespaces(r, aggregateRootPre, aggregateRootAliasLen)
	if !ok {
		return nil, "", errors.Errorf("路径(%v)尝试聚合根时失败。", r)
	}
	if len(ars) > 1 {
		return nil, "", errors.Errorf("路径(%v)尝试聚合根时失败:聚合根不能分级", r)
	}
	if len(ars) == 1 {
		return creAggregateRootItem(ars[0]), r, nil
	}
	return nil, r, nil
}

// creAggregateRootItem 创建聚合根Item
func creAggregateRootItem(id string) *aggregateRootItem {
	return &aggregateRootItem{
		pathItem: pathItem{
			curName: id,
			kind:    AggregateRoot,
			fmtPath: fmtAggregateRootPath,
		},
	}
}

type aggregateRootItem struct {
	pathItem
}

func (i *aggregateRootItem) Append(in Item) {

	i.next = in
	in.setParent(i)
	_ = in.resetPath()
}

func (i *aggregateRootItem) parse(r string) error {
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
	// 尝试解析聚合
	agg, err := parseAggregate(r)
	if err != nil {
		return err
	}
	if agg != nil {
		agg.parent = i
		i.next = agg
		_ = agg.Path()

		return nil
	}

	return errors.Errorf("解析聚合根(%v)后的路径(%v)失败", i.curName, r)

}
