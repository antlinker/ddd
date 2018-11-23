package mgorepo

import (
	"reflect"

	"github.com/antlinker/ddd"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// BaseRepo 基础操作
type BaseRepo struct {
	ddd.BaseRepository
}

// DB 获取数据库操作
func (r *BaseRepo) DB(ctx ddd.Context) MyMgoer {
	db := FromContext(ctx)
	if db == nil {
		if d := r.Domain(); d != nil {
			db = FromDomain(d)
			if db == nil {
				panic("没有初始化mongodb数据库操作")
			}
		}
	}
	return db
}

// CreateIndex 创建索引
func (r *BaseRepo) CreateIndex(ctx ddd.Context, col Coll) {
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
				ctx.Logger().Warnf("创建索引失败:%v", err)
			}
		}
		return nil
	})
}

// NextID 获取下一个唯一ID
func (r *BaseRepo) NextID() string {
	return bson.NewObjectId().Hex()
}

// DBInsert 插入数据
func (r *BaseRepo) DBInsert(ctx ddd.Context, obj Coll) error {
	return r.DB(ctx).ExecSync(obj.CollName(), func(c *mgo.Collection) error {
		return c.Insert(obj)
	})
}

// DBUpsert 插入或更新数据
func (r *BaseRepo) DBUpsert(ctx ddd.Context, obj Coll) error {
	return r.DB(ctx).ExecSync(obj.CollName(), func(c *mgo.Collection) error {
		_, err := c.UpsertId(obj.CollID(), bson.M{"$set": obj})
		return err
	})
}

// DBUpdate 更新查找到的一条记录
// noset 参数只取第一个，默认false 设置为true不在最外层包装$set 否则自动包装$set
//		主要针对包含$inc $min $max $unset $mul $rename $setOnInsert等评级操作符时使用
func (r *BaseRepo) DBUpdate(ctx ddd.Context, collname string, query interface{}, update interface{}, noset ...bool) error {
	return r.DB(ctx).ExecSync(collname, func(c *mgo.Collection) error {
		if len(noset) > 0 && noset[0] {
			return c.Update(query, update)
		}
		return c.Update(query, bson.M{"$set": update})

	})
}

// DBInc 设置某个字段自增
func (r *BaseRepo) DBInc(ctx ddd.Context, collname string, query interface{}, update interface{}) error {
	return r.DBUpdate(ctx, collname, query, bson.M{"$inc": update}, true)
}

// DBUpdateAll 更新查找到的多条记录
// noset 参数只取第一个，默认false 设置为true不在最外层包装$set 否则自动包装$set
//		主要针对包含$inc $min $max $unset $mul $rename $setOnInsert等评级操作符时使用
func (r *BaseRepo) DBUpdateAll(ctx ddd.Context, collname string, query interface{}, update interface{}, noset ...bool) error {
	return r.DB(ctx).ExecSync(collname, func(c *mgo.Collection) error {
		if len(noset) > 0 && noset[0] {
			_, err := c.UpdateAll(query, update)
			return err
		}
		_, err := c.UpdateAll(query, bson.M{"$set": update})
		return err

	})
}

// DBUpdateID 更新指定记录
// noset 参数只取第一个，默认false 设置为true不在最外层包装$set 否则自动包装$set
//		主要针对包含$inc $min $max $unset $mul $rename $setOnInsert等评级操作符时使用
func (r *BaseRepo) DBUpdateID(ctx ddd.Context, collname string, id interface{}, update interface{}, noset ...bool) error {
	return r.DB(ctx).ExecSync(collname, func(c *mgo.Collection) error {
		if len(noset) > 0 && noset[0] {

			return c.UpdateId(id, update)
		}
		return c.UpdateId(id, bson.M{"$set": update})

	})
}

// DBUpdateCollAggregate 更新指定聚合记录 并根据配置自动刷新缓存
// ca 需要更新的聚合
// update 更新语句
// noset 参数只取第一个，默认false 设置为true不在最外层包装$set 否则自动包装$set
//		主要针对包含$inc $min $max $unset $mul $rename $setOnInsert等评级操作符时使用
func (r *BaseRepo) DBUpdateCollAggregate(ctx ddd.Context, ca CollAggregate, update interface{}, noset ...bool) error {

	return r.DB(ctx).ExecSync(ca.CollName(), func(c *mgo.Collection) error {
		var err error
		id := ca.CollID()
		if len(noset) > 0 && noset[0] {

			err = c.UpdateId(id, update)
		} else {

			err = c.UpdateId(id, bson.M{"$set": update})
		}
		if err != nil {
			return err
		}
		if err = r.dbGetOne(c, id, ca); err != nil {
			return err
		}
		_, err = r.refeshCache(ctx, ca)
		return err
	})
}

// DBDestroyID 销魂指定记录，物理删除
func (r *BaseRepo) DBDestroyID(ctx ddd.Context, collname string, id interface{}) error {
	if err := r.DB(ctx).ExecSync(collname, func(c *mgo.Collection) error {
		return c.RemoveId(id)
	}); err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

// DBGetOne 获取数据库操作
func (r *BaseRepo) DBGetOne(ctx ddd.Context, obj Coll) error {
	return r.DB(ctx).ExecSync(obj.CollName(), func(c *mgo.Collection) error {
		return r.dbGetOne(c, obj.CollID(), obj)
	})

	//return r.dbGetOne(ctx, obj, obj)
}

// DBGetOne 获取数据库操作
func (r *BaseRepo) dbGetOne(c *mgo.Collection, id interface{}, result interface{}) error {
	return c.FindId(id).One(result)
}

// DBExecSync 同步执行数据库操作
func (r *BaseRepo) DBExecSync(ctx ddd.Context, collname string, task func(coll *mgo.Collection) error) error {
	return r.DB(ctx).ExecSync(collname, task)
}

// DBExecAsync 异步执行数据库操作
func (r *BaseRepo) DBExecAsync(ctx ddd.Context, collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{}) {
	r.DB(ctx).ExecAsync(collname, task, opt...)
}

// DBQueryPage 进行分页查询
func (r *BaseRepo) DBQueryPage(ctx ddd.Context, collname string, pi PageInfo, pageidname string, query bson.M, field interface{}, sort []string, result interface{}) (pr PageResult, err error) {
	var nsort []string
	if pi.Desc {
		nsort = append(nsort, "-"+pageidname)
	} else {
		nsort = append(nsort, pageidname)
	}
	if sort != nil {
		nsort = append(nsort, sort...)
	}
	if pi.Mode == 1 {
		return r.dbQueryPage1(ctx, collname, pi, pageidname, query, field, nsort, result)
	}
	return r.dbQueryPage0(ctx, collname, pi, pageidname, query, field, nsort, result)

}

// DBQueryNum 查询匹配数量
func (r *BaseRepo) DBQueryNum(ctx ddd.Context, collname string, query bson.M) (num int, err error) {

	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) (err error) {
		num, err = c.Find(query).Count()
		return err
	})
	return

}

// DBQueryPage 进行分页查询
func (r *BaseRepo) dbQueryPage0(ctx ddd.Context, collname string, pi PageInfo, pageidname string, query bson.M, field interface{}, sort []string, result interface{}) (pr PageResult, err error) {
	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) error {
		ctx.Logger().Debugf("dbQueryPage0 %s find :%v", collname, query)
		pr.Total, err = c.Find(query).Count()
		if err != nil {
			if err == mgo.ErrNotFound {
				err = nil
			}
			return err
		}
		if pr.Total == 0 {
			return nil
		}
		q := c.Find(query)
		if field != nil {
			q = q.Select(field)
		}
		if len(sort) > 0 {
			q = q.Sort(sort...)
		}
		if pi.Current == 0 {
			pi.Current = 1
		}
		if pi.Current > 0 {
			start := (pi.Current - 1) * pi.PageSize
			q = q.Skip(start)
		}
		if pi.PageSize > 0 {
			q = q.Limit(pi.PageSize)
		}
		return q.All(result)
	})
	if err == mgo.ErrNotFound {
		err = nil
		pr.End = 1
		return
	}
	pr.Current = pi.Current
	pr.PageSize = pi.PageSize
	if result != nil {
		e := reflect.ValueOf(result).Elem()
		if e.Kind() == reflect.Slice {
			len := e.Len()
			if pi.PageSize > len {
				pr.End = 1
			} else if pi.PageSize == len && pr.Total == pr.Current*pr.PageSize {
				pr.End = 1
			}
		}
	}

	return
}

// dbQueryPage1 进行分页查询 1 模式
func (r *BaseRepo) dbQueryPage1(ctx ddd.Context, collname string, pi PageInfo, pageidname string, query bson.M, field interface{}, sort []string, result interface{}) (pr PageResult, err error) {

	if pi.PageID != nil {
		if pi.Desc {
			if pi.Direct == 0 {
				query[pageidname] = bson.M{"$lt": pi.PageID}
			} else {
				query[pageidname] = bson.M{"$gt": pi.PageID}
			}
		} else {
			if pi.Direct == 0 {
				query[pageidname] = bson.M{"$gt": pi.PageID}
			} else {
				query[pageidname] = bson.M{"$lt": pi.PageID}
			}
		}
	}
	ctx.Logger().Debugf("dbQueryPage1 %s find :%v", collname, query)
	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) error {
		q := c.Find(query)
		if field != nil {
			q = q.Select(field)
		}
		if len(sort) > 0 {
			q = q.Sort(sort...)
		}
		if pi.PageSize > 0 {
			q = q.Limit(pi.PageSize)
		}
		return q.All(result)

	})
	if err == mgo.ErrNotFound {
		err = nil
		pr.End = 1
		return
	}
	pr.Current = pi.Current
	pr.PageSize = pi.PageSize
	if result != nil {
		e := reflect.ValueOf(result).Elem()
		if e.Kind() == reflect.Slice {
			len := e.Len()
			if pi.PageSize > len {
				pr.End = 1
			}
		}
	}

	return
}

// DBQuery 进行数据库查询，返回所有查询结果到result中
func (r *BaseRepo) DBQuery(ctx ddd.Context, collname string, query interface{}, field interface{}, sort []string, result interface{}) (err error) {
	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) error {
		q := c.Find(query)
		if field != nil {
			q = q.Select(field)
		}
		if len(sort) > 0 {
			q = q.Sort(sort...)
		}
		ctx.Logger().Debugf("DBQuery ALL %s <find:%v>  <select:%v>  <sort:%v>", collname, query, field, sort)
		return q.All(result)

	})
	return
}

// DBQueryOne 进行数据库查询，返回一个查询结果到result中
func (r *BaseRepo) DBQueryOne(ctx ddd.Context, collname string, query interface{}, field interface{}, result interface{}) (err error) {
	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) error {
		q := c.Find(query)
		if field != nil {
			q = q.Select(field)
		}
		ctx.Logger().Debugf("DBQuery One %s <find:%v>  <select:%v>", collname, query, field)
		return q.One(result)

	})
	return
}

// DBExists 进行数据库查询，返回所有查询结果到result中
func (r *BaseRepo) DBExists(ctx ddd.Context, collname string, query interface{}) (ok bool, err error) {
	err = r.DBExecSync(ctx, collname, func(c *mgo.Collection) error {
		if n, err := c.Find(query).Count(); err != nil && err != mgo.ErrNotFound {
			return err
		} else if n > 0 {
			ok = true
		}
		return nil

	})
	return
}

// RefeshCache 刷新缓存
func (r *BaseRepo) RefeshCache(ctx ddd.Context, c CollAggregate) (a ddd.Aggregate, err error) {
	// 从数据库查询

	err = r.DBGetOne(ctx, c)
	if err != nil {
		return
	}
	return r.refeshCache(ctx, c)
}

// RefeshCache 刷新缓存
func (r *BaseRepo) refeshCache(ctx ddd.Context, c CollAggregate) (a ddd.Aggregate, err error) {
	// 从数据库查询
	a = c.ToAggregate()
	a.Init(a, r.AggregateRoot(), c.CollID())
	if cc, ok := c.(CollAggregateCache); ok {
		if cc.IsCache() {
			// 添加到缓存
			if cce, ok := c.(CollAggregateCacheExpire); ok {
				err1 := r.CacheSet(a, cce.CacheExpire())
				if err1 != nil {
					ctx.Logger().Warnf("设置缓存失败：%v", err1)
				}
			} else {
				err1 := r.CacheSet(a)
				if err1 != nil {
					ctx.Logger().Warnf("设置缓存失败：%v", err1)
				}
			}

		}
	}

	return
}

// DeleteCache 删除缓存
func (r *BaseRepo) DeleteCache(ctx ddd.Context, c CollAggregate) (err error) {
	// 从数据库查询

	a := c.ToAggregate()
	a.Init(a, r.AggregateRoot(), c.CollID())
	r.CacheDelete(a)
	return nil
}

// GetByCollAggregate 获取一个聚合信息
// 该方法会自动缓存
func (r *BaseRepo) GetByCollAggregate(ctx ddd.Context, c CollAggregate) (a ddd.Aggregate, err error) {
	a = c.ToAggregate()
	// 初始化聚合
	a.Init(a, r.AggregateRoot(), c.CollID())
	// 获取缓存
	if ok, err1 := r.CacheGet(a); ok {
		return
	} else if err1 != nil {
		ctx.Logger().Warnf("获取缓存失败:%v", err1)
	}
	// 不存在缓存设置缓存，并返回结果
	return r.RefeshCache(ctx, c)

}

// DestroyByCollAggregate 彻底删除指定信息
// 同时会销毁缓存
func (r *BaseRepo) DestroyByCollAggregate(ctx ddd.Context, c CollAggregate) error {
	_ = r.DeleteCache(ctx, c)
	return r.DBDestroyID(ctx, c.CollName(), c.CollID())
}

// // ParseQueryPage 为查询条件追加分页查询
// func (r *BaseRepo) ParseQueryPage(query bson.M, field string, pi PageInfo) bson.M {
// 	if pi.Mode == 1 {
// 		if pi.PageID != "" {
// 			if pi.Direct == 0 {

// 				query[field] = bson.M{"$gt": pi.PageID}
// 			} else {
// 				query[field] = bson.M{"$lt": pi.PageID}

// 			}
// 		}
// 	}
// 	return query
// }
