package memevt

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/antlinker/ddd"
	"github.com/antlinker/taskpool"
)

var (
	eventTaskOperate taskpool.AsyncTaskOperater
)

func getEventTaskExecor() taskpool.AsyncTaskOperater {
	if eventTaskOperate == nil {
		eventTaskOperate = taskpool.CreateAsyncTaskOperater("事件任务", &eventTaskExecor{}, &taskpool.AsyncTaskOption{
			AsyncMaxWaitTaskNum: 1024,
			//最大异步任务go程数量
			MaxAsyncPoolNum: 128,
			MinAsyncPoolNum: 8,
			//最大空闲时间
			AsyncPoolIdelTime: 10 * time.Second,
			//任务最大失败次数
			AsyncTaskMaxFaildedNum: 1,
		})
	}
	return eventTaskOperate
}

//监听组 该组是线程不安全的
type listenerGroup struct {
	listenerList      *list.List
	asyncTaskOperater taskpool.AsyncTaskOperater
}

func createListenerGroup() *listenerGroup {
	group := &listenerGroup{}
	group.listenerList = list.New()
	group.asyncTaskOperater = getEventTaskExecor()
	return group
}

//向组内增加监听
func (l *listenerGroup) AddListener(listener ddd.EventListener) {
	if l.listenerList == nil {
		l.listenerList = list.New()

		return
	}
	l.listenerList.PushFront(listener)
	// _, ok := l.checkElement(listener)
	// if !ok {
	// 	l.listenerMap[listener] = l.listenerList.PushFront(listener)
	// }

}

//从组内移除监听
func (l *listenerGroup) RemoveListener(listener ddd.EventListener) {
	// elem, ok := l.checkElement(listener)
	// if ok {
	// 	l.listenerList.Remove(elem)
	// 	delete(l.listenerMap, listener)
	// }
	// TODO: 暂时还没有实现移除事件
	panic("暂时还没有实现移除事件")
}

//FireListener 触发监听事件
func (l *listenerGroup) TriggerEvent(event ddd.Event) {
	for elem := l.listenerList.Front(); elem != nil; elem = elem.Next() {
		listener := elem.Value.(ddd.EventListener)
		l.asyncTaskOperater.ExecAsyncTask(createTriggerEvent(event, listener))

	}

}

// //检测组内是否有该监听
// func (l *listenerGroup) checkElement(listener ddd.EventListener) (*list.Element, bool) {
// 	if l.listenerMap == nil {
// 		return nil, false
// 	}
// 	elem, ok := l.listenerMap[listener]
// 	return elem, ok

// }

// NewEventTrigger 创建一个事件触发器
func NewEventTrigger() ddd.EventTrigger {
	return &dddEventTrigger{
		listenergroupMap: make(map[string]*listenerGroup),
		evtactions:       make(map[string][]string),
	}
}

// dddEventTrigger 事件发生器 线程安全的
type dddEventTrigger struct {
	sync.Mutex
	listenergroupMap map[string]*listenerGroup
	evtactions       map[string][]string
}

// AddListener 增加一个监听器
// eventtype 事件类型
// listener 监听器
func (eg *dddEventTrigger) AddListener(eventtype string, method string, listener ddd.EventListener) {
	eg.Lock()
	defer eg.Unlock()
	if method == "" {

		listenerGroup := eg.tryInit(eventtype)
		listenerGroup.AddListener(listener)
		return
	}
	em, ok := eg.evtactions[eventtype]
	if !ok {
		em = make([]string, 0)
		eg.evtactions[eventtype] = em
	}
	listenerGroup := eg.tryInit(eventtype + "." + method)
	listenerGroup.AddListener(listener)
}

// RemoveListener 移除一个监听
// eventtype 事件类型
// listener 监听器
func (eg *dddEventTrigger) RemoveListener(eventtype string, method string, listener ...ddd.EventListener) {
	eg.Lock()
	defer eg.Unlock()
	if eg.listenergroupMap == nil {
		return
	}
	if eventtype == "" {
		return
	}
	if method == "" {
		if len(listener) > 0 {
			listenerGroup := eg.tryInit(eventtype)
			for _, l := range listener {
				listenerGroup.RemoveListener(l)
			}
		} else {
			delete(eg.listenergroupMap, eventtype)
		}
		if et, ok := eg.evtactions[eventtype]; ok {
			if len(listener) > 0 {
				for _, v := range et {
					listenerGroup := eg.tryInit(eventtype + "." + v)
					for _, l := range listener {
						listenerGroup.RemoveListener(l)
					}
				}
			} else {
				for _, v := range et {
					delete(eg.listenergroupMap, eventtype+"."+v)
				}
				delete(eg.evtactions, eventtype)
			}
		}
		return
	}
	if len(listener) > 0 {
		listenerGroup := eg.tryInit(eventtype + "." + method)
		for _, l := range listener {
			listenerGroup.RemoveListener(l)
		}
	} else {
		delete(eg.listenergroupMap, eventtype+"."+method)
		delete(eg.evtactions, eventtype)
	}
}

// FireListener 触发事件
// eventtype 事件类型
// method 触发的方法
// event 事件
// param 事件自定义参数
func (eg *dddEventTrigger) TriggerEvent(evt ddd.Event) {
	//eg.Lock()
	//defer eg.Unlock()
	if lg, ok := eg.listenergroupMap[evt.Type]; ok {
		lg.TriggerEvent(evt)
	}
	if lg, ok := eg.listenergroupMap[evt.Type+"."+evt.Action]; ok {
		lg.TriggerEvent(evt)
	}
}

//尝试初始化,如果已经初始化不在初始化
func (eg *dddEventTrigger) tryInit(eventtype string) *listenerGroup {
	if eg.listenergroupMap == nil {

		eg.listenergroupMap = make(map[string]*listenerGroup)
		lg := createListenerGroup()
		eg.listenergroupMap[eventtype] = lg
		return lg

	}
	listenerGroup, ok := eg.listenergroupMap[eventtype]
	if !ok {

		listenerGroup = createListenerGroup()
		eg.listenergroupMap[eventtype] = listenerGroup

	}
	return listenerGroup
}
func createTriggerEvent(event ddd.Event, listener ddd.EventListener) *triggerEvent {
	return &triggerEvent{
		event:    event,
		listener: listener,
	}
}

type triggerEvent struct {
	taskpool.BaseTask
	event    ddd.Event
	listener ddd.EventListener
}

// // Listener 监听接口
// type Listener interface {
// }

type eventTaskExecor struct {
}

func (e *eventTaskExecor) ExecTask(task taskpool.Task) error {
	even, ok := task.(*triggerEvent)
	if !ok {
		return fmt.Errorf("错误的任务:%v", task)
	}

	return even.listener(even.event)
}
