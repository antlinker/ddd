package path

import (
	"github.com/pkg/errors"
)

func parseAggregate(r string) (*pathItem, error) {
	// 尝试解聚合
	agg, r, ok := parseNamespaces(r, aggregatePre, aggregateAliasLen)
	if !ok {
		return nil, errors.Errorf("路径(%v)尝试解析聚合时失败.", r)
	}
	if len(agg) == 1 && r != "" {
		return nil, errors.Errorf("路径(%v)尝试解析聚合时失败,有多余的字符存在", r)
	}
	if len(agg) == 1 {
		return creAggregateItem(agg[0]), nil
	} else if len(agg) > 1 {
		return nil, errors.Errorf("路径(%v)尝试解析聚合时失败,聚合不能分级", r)
	}
	return nil, nil
}

// creAggregateItem 创建聚合Item
func creAggregateItem(id string) *pathItem {
	return &pathItem{
		curName: id,
		kind:    Aggregate,
		fmtPath: fmtAggregatePath,
	}
}
