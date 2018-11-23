package ddd

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/antlinker/ddd/log"
)

// Event 事件信息
type Event struct {
	Type    string      // 事件类型
	TraceID string      // 追中ID
	Action  string      // 行为
	ID      string      // 事件编号
	Occtime time.Time   // 发生时间
	From    string      // 事件来源
	Data    interface{} // 事件数据
}

// EventTrigger 事件触发器
type EventTrigger interface {
	TriggerEvent(e Event)
	EventListenManager
}

// EventListener 监听函数
type EventListener func(evt Event) error

// EventListenManager 事件监听管理器
type EventListenManager interface {
	AddListener(etype, act string, l EventListener)
	RemoveListener(etype, act string, l ...EventListener)
}

var evtflow int64

// TriggerForEvent 触发事件
func TriggerForEvent(evt Event) {
	id := atomic.AddInt64(&evtflow, 1)
	evt.ID = fmt.Sprintf("%09d%03d", time.Now().Unix(), id)
	evt.Occtime = time.Now()
	if evt.TraceID == "" {
		evt.TraceID = evt.ID
	}
	defEvtTrigger.TriggerEvent(evt)
}

// Trigger 触发事件
func Trigger(etype, action, from string, data interface{}) {
	evt := Event{
		Type:   etype,
		Action: action,
		From:   from,
		Data:   data,
	}
	TriggerForEvent(evt)
}

var (
	defEvtTrigger EventTrigger
)

// RegEventTrigger 注册事件触发器
func RegEventTrigger(et EventTrigger) {
	defEvtTrigger = et
}

// AddEventListener 注册事件监听
func AddEventListener(etype, act string, l EventListener) {
	defEvtTrigger.AddListener(etype, act, l)
}

// RemoveEventListener 注册事件监听
func RemoveEventListener(etype, act string, l EventListener) {
	defEvtTrigger.RemoveListener(etype, act, l)
}

// CreEventListener 创建事件监听
func CreEventListener(handler func(ctx Context, evt Event) error) EventListener {
	return func(evt Event) (err error) {
		start := time.Now()

		ExecContextForRepo(func(ctx Context) {
			WrapTraceID(ctx, evt.TraceID)
			var logger log.ILogger
			defer func() {

				if err != nil {
					_, file, line, _ := runtime.Caller(2)
					// 执行事件监听失败
					logger.Errorf("exec %v:%v listener  (%v:%v:%v),(%v :%v): error:%+v ", file, line, evt.Type, evt.Action, evt.ID, start, time.Since(start), err)
					return
				}
				logger.Debugf("exec listener (%v:%v:%v),(%v :%v): ok.", evt.Type, evt.Action, evt.ID, start, time.Since(start))
			}()
			logger = ctx.Logger()
			err = handler(ctx, evt)
		})
		return
	}
}

const (
	EvtTypeDDD    = "ddd"
	EvtActStartOK = "startok"
)

// StartOK 服务启动成功时调用改事件
func StartOK() {
	Trigger(EvtTypeDDD, EvtActStartOK, "srv", nil)
}
