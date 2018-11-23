package path

import (
	"github.com/pkg/errors"
)

func parseRepository(r string) (*pathItem, error) {
	// 尝试解析仓储
	repo, r, ok := parseNamespaces(r, repositoryPre, repositoryAliasLen)
	if !ok {
		return nil, errors.Errorf("路径(%v)尝试解析仓储时失败.", r)
	}
	if len(repo) == 1 && r != "" {
		return nil, errors.Errorf("路径(%v)尝试解析仓储时失败,有多余的字符存在", r)
	}
	if len(repo) == 1 {
		return creRepositoryItem(repo[0]), nil
	} else if len(repo) > 1 {
		return nil, errors.Errorf("路径(%v)尝试解析仓储时失败,仓储不能分级", r)
	}
	return nil, nil
}

// creRepositoryItem 创建仓储Item
func creRepositoryItem(id string) *pathItem {
	return &pathItem{
		curName: id,
		kind:    Repository,
		fmtPath: fmtRepositoryPath,
	}
}
