package ddd

import (
	"context"
	"sync"
	"time"

	"github.com/antlinker/ddd/log"
)

// Context 上下文
type Context interface {
	context.Context
	UID() string
	Put(key interface{}, val interface{})
	Get(key interface{}) (interface{}, bool)
	Logger() ILogger
	TraceID() string
	Trigger(etype, action, from string, data interface{})
}

// NewTraceContext 创建一个Context
func NewTraceContext(ctx context.Context, traceid, uid string, l ILogger) Context {
	return &_context{
		traceid: traceid,
		uid:     uid,
		ctx:     ctx,
		ext:     make(map[interface{}]interface{}),
		log:     l,
	}
}

// NewContext 创建一个Context
func NewContext(ctx context.Context, uid string, l ILogger) Context {
	return &_context{
		uid: uid,
		ctx: ctx,
		ext: make(map[interface{}]interface{}),
		log: l,
	}
}

// _context 上下文
type _context struct {
	traceid string
	uid     string
	log     ILogger
	ctx     context.Context
	ext     map[interface{}]interface{}
	sync.RWMutex
}

func (c *_context) Put(key interface{}, val interface{}) {
	c.Lock()
	c.ext[key] = val
	c.Unlock()
}
func (c *_context) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	v, ok := c.ext[key]
	c.RUnlock()
	return v, ok
}

// UID 返回uid
func (c *_context) UID() string {
	return c.uid
}

// TraceID 返回traceid
func (c *_context) TraceID() string {
	return c.traceid
}

func (c *_context) Logger() ILogger {
	return c.log
}

// Done ctx Done
func (c *_context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err ctx Err
func (c *_context) Err() error {
	return c.ctx.Err()
}

// Deadline ctx Deadline
func (c *_context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Value ctx Value
func (c *_context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *_context) Trigger(etype, action, from string, data interface{}) {
	evt := Event{
		Type:   etype,
		Action: action,
		From:   from,
		Data:   data,
	}
	evt.TraceID = c.traceid

	TriggerForEvent(evt)
}

// TraceLog 带有追踪ID的日志
type TraceLog interface {
	SetTraceID(traceid string)
}

// WrapTraceID 追踪ID包装到上下文中
func WrapTraceID(c Context, traceid string) {
	if tc, ok := c.(*_context); ok {
		tc.traceid = traceid
		if l, ok := tc.log.(TraceLog); ok {
			l.SetTraceID(traceid)
		}
	}
}

// WrapLogger 日志管理器包装到上下文中
func WrapLogger(c Context, logger ILogger) {
	if tc, ok := c.(*_context); ok {
		tc.log = logger
	}
}

var (
	defaultExecContext execContext
)

// SetContextForRepo 设置一个执行时有Context并能运行Repo方法的环境
func SetContextForRepo(exeHanler ExecCtxForRepoHandler) {
	defaultExecContext = execContext{
		exeHanler: exeHanler,
	}
}

// ExecContextForRepo 有Context并可以运行Repo的环境执行hanler
func ExecContextForRepo(handler ContextHandler) {
	defaultExecContext.ExecContextForRepo(handler)
}

// ExecCtxForRepoHandler 执行时有Context并能运行Repo方法的环境函数
type ExecCtxForRepoHandler func(ctx Context, handler ContextHandler)

// ContextHandler 执行函数
type ContextHandler func(ctx Context)

type execContext struct {
	exeHanler ExecCtxForRepoHandler
}

var (
	newTraceLogHandler = func() ILogger {
		return log.NewTraceLog("")
	}
	releaseTraceLog = func(l ILogger) {
		log.ReleaseTraceLog(l)
	}
	newLogHandler = func() ILogger {
		return log.Logger()
	}
)

// ILogger 日志类型log.ILogger别名
type ILogger = log.ILogger

// NewLogger 日志生成器
type NewLogger func() ILogger

// ReleaseTraceLog 追踪日志释放
type ReleaseTraceLog func(ILogger)

// SetTraceLogFactory 设置追踪日志生产工厂
func SetTraceLogFactory(log NewLogger, rlog ReleaseTraceLog) {
	newTraceLogHandler = log
	releaseTraceLog = rlog
}

func (c execContext) ExecContextForRepo(handler ContextHandler) {
	if newTraceLogHandler != nil {
		l := newTraceLogHandler()
		defer releaseTraceLog(l)
		ctx := NewTraceContext(context.Background(), "", "", l)
		c.exeHanler(ctx, handler)
		return
	}
	if newLogHandler != nil {

		ctx := NewContext(context.Background(), "", newLogHandler())
		c.exeHanler(ctx, handler)
		return
	}

	ctx := NewContext(context.Background(), "", nil)
	c.exeHanler(ctx, handler)

}
