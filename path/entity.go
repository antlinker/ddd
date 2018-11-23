package path

import (
	"github.com/pkg/errors"
)

func parseEntity(r string) (*pathItem, error) {
	// 尝试解实例
	repo, r, ok := parseNamespaces(r, entityPre, entityAliasLen)
	if !ok {
		return nil, errors.Errorf("路径(%v)尝试解析实例时失败.", r)
	}
	if len(repo) == 1 && r != "" {
		return nil, errors.Errorf("路径(%v)尝试解析实例时失败,有多余的字符存在", r)
	}
	if len(repo) == 1 {
		return creEntityItem(repo[0]), nil
	} else if len(repo) > 1 {
		return nil, errors.Errorf("路径(%v)尝试解析实例时失败,实例不能分级", r)
	}
	return nil, nil
}

// creEntityItem 创建实例Item
func creEntityItem(id string) *pathItem {
	return &pathItem{
		curName: id,
		kind:    Entity,
		fmtPath: fmtEntityPath,
	}
}
