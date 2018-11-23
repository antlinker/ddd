package mgorepo

import (
	"github.com/antlinker/ddd"
	"github.com/antlinker/ddd/log"
	"github.com/globalsign/mgo"
)

var defIndexRepo = &indexRepo{}

// RegColForIndex 注册需要创建索引的集合
func RegColForIndex(col Coll) {
	defIndexRepo.cols = append(defIndexRepo.cols, col)
}

// CreateAllIndex 创建索引
func CreateAllIndex() {
	ddd.ExecContextForRepo(defIndexRepo.CreateAllIndex)
}

type indexRepo struct {
	BaseRepo
	cols []Coll
}

// CreateIndex 创建索引
func (r *indexRepo) CreateAllIndex(ctx ddd.Context) {
	for _, v := range r.cols {
		r.CreateIndex(ctx, v)
	}
}

// CreateIndex 创建索引
func (r *indexRepo) CreateIndex(ctx ddd.Context, col Coll) {
	colind, ok := col.(CollIndex)
	if !ok {
		return
	}
	indexes := colind.Indexes()
	if len(indexes) == 0 {
		return
	}
	_ = r.DB(ctx).ExecSync(col.CollName(), func(c *mgo.Collection) error {
		for _, v := range indexes {
			if err := c.EnsureIndex(v); err != nil {
				log.Warnf("创建集合(%v)索引失败:%v", col.CollName(), err)
				continue
			}
			log.Infof("创建集合(%v)索引成功:%v", col.CollName(), v)
		}
		return nil
	})
}
