package mgorepo

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/antlinker/ddd/log"
	util "github.com/antlinker/taskpool"

	"github.com/globalsign/mgo"
)

// MgoTask 执行任务任务
type MgoTask struct {
	util.BaseTask
	collname string
	task     func(coll *mgo.Collection, opt ...interface{}) error
	opt      []interface{}
}

func (t *MgoTask) call(coll *mgo.Collection, opt ...interface{}) error {
	return t.task(coll, opt...)
}

// CreateTask 创建一个mongodb任务,指定集,合任务
func CreateTask(collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{}) *MgoTask {
	return &MgoTask{
		collname: collname,
		task:     task,
		opt:      opt,
	}
}

// CreateDBTask 创建一个mongodb任务,指定数据库,集合,任务
func CreateDBTask(collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{}) *MgoTask {

	return &MgoTask{
		collname: collname,
		task:     task,
	}
}

// MyMgoer mgo数据库操作
type MyMgoer interface {
	Clone() MyMgoer
	Release()
	ExecSync(collname string, task func(coll *mgo.Collection) error) error

	//ExecDBAsync(db, collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{})
	ExecAsyncTask(task *MgoTask)
	ExecAsync(collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{})
	Session() *mgo.Session
}

var lock = &sync.Mutex{}

// MongodbConfig mongodb 配置
type MongodbConfig struct {
	URL                    string `json:"url" yaml:"url"`
	AsyncMaxPoolNum        int64  `json:"asyncMaxPoolNum" yaml:"asyncMaxPoolNum"`               //异步任务最大执行数量
	AsyncMinPoolNum        int64  `json:"asyncMinPoolNum" yaml:"asyncMinPoolNum"`               //异步任务最小执行数量
	AsyncMaxWaitTaskNum    int    `json:"asyncTaskMaxFaildedNum" yaml:"asyncTaskMaxFaildedNum"` //异步任务最大等待数量
	AsyncPoolIdelTime      int    `json:"asyncPoolIdelTime" yaml:"asyncPoolIdelTime"`           //最大空闲时间,空闲时间到自动回收 单位秒
	AsyncTaskMaxFaildedNum uint64 `json:"asyncTaskMaxExeNum" yaml:"asyncTaskMaxExeNum"`         //异步任务最大重复执行次数,超过后放弃
	DBName                 string `json:"dbName" yaml:"dbName"`
	Debug                  bool   `json:"debug" yaml:"debug"` // debug 模式
}

// CreateMongodbConfigForEnv 通过环境变量创建mongodb连接信息
func CreateMongodbConfigForEnv(pre string) *MongodbConfig {
	tmp := os.Getenv(pre + "_MGO_URL")
	if tmp == "" {
		return nil
	}
	var cfg MongodbConfig
	cfg.URL = tmp
	tmp = os.Getenv(pre + "_MGO_DBNAME")
	if tmp != "" {
		cfg.DBName = tmp
	}
	tmp = os.Getenv(pre + "_MGO_MAXPOOLNUM")
	if tmp != "" {
		cfg.AsyncMaxPoolNum, _ = strconv.ParseInt(tmp, 0, 64)
	}
	tmp = os.Getenv(pre + "_MGO_MINPOOLNUM")
	if tmp != "" {
		cfg.AsyncMinPoolNum, _ = strconv.ParseInt(tmp, 0, 64)
	}
	tmp = os.Getenv(pre + "_MGO_MAXWAITTASKNUM")
	if tmp != "" {
		cfg.AsyncMaxWaitTaskNum, _ = strconv.Atoi(tmp)
	}
	tmp = os.Getenv(pre + "_MGO_POOLIDELTIME")
	if tmp != "" {
		cfg.AsyncPoolIdelTime, _ = strconv.Atoi(tmp)
	}
	tmp = os.Getenv(pre + "_MGO_MAXFAILDEDNUM")
	if tmp != "" {
		cfg.AsyncTaskMaxFaildedNum, _ = strconv.ParseUint(tmp, 0, 64)
	}
	return &cfg
}

type mgologHandler struct {
}

func (mgologHandler) Output(calldepth int, s string) error {
	log.Debugf("[mgo] %d=>%v", calldepth, s)
	return nil
}

// CreateMyMgoForCfg 通过配置参数创建db
func CreateMyMgoForCfg(cfg MongodbConfig) MyMgoer {

	lock.Lock()
	defer lock.Unlock()
	if cfg.Debug {
		mgo.SetDebug(true)
		mgo.SetLogger(mgologHandler{})
	}
	mymgo := &myMgo{
		url:    cfg.URL,
		dbname: cfg.DBName,
	}
	options := &util.AsyncTaskOption{
		AsyncTaskMaxFaildedNum: cfg.AsyncTaskMaxFaildedNum,
		AsyncMaxWaitTaskNum:    cfg.AsyncMaxWaitTaskNum,
		MaxAsyncPoolNum:        cfg.AsyncMaxPoolNum,
		MinAsyncPoolNum:        cfg.AsyncMinPoolNum,
		AsyncPoolIdelTime:      time.Duration(cfg.AsyncPoolIdelTime) * time.Second,
	}
	mymgo.asyncTaskOperater = util.CreateAsyncTaskOperater("mongodb读写任务", mymgo, options)
	log.Debugf("连接MonggoDB:%v", decodeUrl(mymgo.url))
	tmp, err := mgo.Dial(mymgo.url)
	if err != nil {
		panic(errors.New("创建mongodb连接失败1:" + fmt.Sprintf("%v\n%v", err, cfg)))
	}
	mymgo.session = tmp
	// Optional. Switch the session to a monotonic behavior.
	mymgo.session.SetMode(mgo.Eventual, true)
	mymgo.session.SetMode(mgo.Strong, true)
	mymgo.db = mymgo.session.DB(mymgo.dbname)
	return mymgo
}

const (
	rexptpl = "^mongo\\://\\w+?@\\w+?\\:(\\S+)$"
)

var (
	rexp, _ = regexp.Compile(rexptpl)
)

func decodeUrl(url string) string {

	return rexp.ReplaceAllString(url, "mongo://*@*:$1")
}

// CreateMyMgo 创建数据库连接,其他使用默认设置
func CreateMyMgo(url, dbname string) MyMgoer {
	defaultOpt := MongodbConfig{
		URL:                    url,
		DBName:                 dbname,
		AsyncTaskMaxFaildedNum: 1,
		AsyncMaxWaitTaskNum:    10,
		AsyncMaxPoolNum:        32,
		AsyncMinPoolNum:        1,
		AsyncPoolIdelTime:      30,
	}
	return CreateMyMgoForCfg(defaultOpt)
}

type myMgo struct {
	url               string
	session           *mgo.Session
	db                *mgo.Database
	asyncTaskOperater util.AsyncTaskOperater
	dbname            string
}

func (m myMgo) Session() *mgo.Session {
	return m.session.Clone()
}
func (m myMgo) Clone() MyMgoer {
	m.session = m.session.Clone()
	return &m
}

func (m *myMgo) Release() {
	m.session.Close()
}
func (m *myMgo) execTask(task *MgoTask) error {
	return task.call(m.db.C(task.collname), task.opt...)

}

//同步执行
func (m *myMgo) ExecSync(collname string, task func(coll *mgo.Collection) error) error {
	return task(m.db.C(collname))
}

//异步执行
func (m *myMgo) ExecAsyncTask(task *MgoTask) {
	m.asyncTaskOperater.ExecAsyncTask(task)

}

// //异步执行
// func (m *myMgo) ExecDBAsync(db, collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{}) {
// 	m.ExecAsyncTask(CreateDBTask(taskName, db, collname, task, opt...))
// }

//异步执行
func (m *myMgo) ExecAsync(collname string, task func(coll *mgo.Collection, opt ...interface{}) error, opt ...interface{}) {
	m.ExecAsyncTask(CreateTask(collname, task, opt...))
}

func (m *myMgo) ExecTask(task util.Task) error {
	mtask, ok := task.(*MgoTask)
	if ok {

		return m.execTask(mtask)
	}
	return nil
}
