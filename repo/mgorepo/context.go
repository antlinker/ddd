package mgorepo

import (
	"github.com/antlinker/ddd"
)

type dbkeys struct {
}

var (
	dbkey = dbkeys{}
)

// DomainInitDB 领域初始化mongodb数据库操作
func DomainInitDB(domain ddd.Domain, myMgo MyMgoer) {
	domain.GlobalStorer().Put(dbkey, myMgo)
}

// SetContextDB 向context中追加数据库操作
func SetContextDB(ctx ddd.Context, domain ddd.Domain) {
	if tmp := domain.GlobalStorer().Get(dbkey); tmp != nil {
		db := tmp.(MyMgoer)
		ctx.Put(dbkey, db.Clone())
	}

}

// FromDomain 从领域中获取数据库操作
func FromDomain(domain ddd.Domain) MyMgoer {
	if tmp := domain.GlobalStorer().Get(dbkey); tmp != nil {
		db := tmp.(MyMgoer)
		return db.Clone()
	}
	return nil
}

// FromContext 从context中获取 数据库操作
func FromContext(ctx ddd.Context) MyMgoer {
	db, ok := ctx.Get(dbkey)
	if !ok {
		return nil
	}
	return db.(MyMgoer)
}

// ReleaseDB 释放数据库
func ReleaseDB(ctx ddd.Context) {
	db := FromContext(ctx)
	if db != nil {
		db.Release()
	}
}

// NewExecCtxForRepoHandler 创建 ddd.ExecCtxForRepoHandler函数
func NewExecCtxForRepoHandler(d ddd.Domain) ddd.ExecCtxForRepoHandler {
	return func(ctx ddd.Context, handler ddd.ContextHandler) {
		SetContextDB(ctx, d)
		defer ReleaseDB(ctx)
		handler(ctx)
	}
}
