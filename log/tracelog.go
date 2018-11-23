package log

import (
	"fmt"
	"sync"

	"github.com/antlinker/ddd/util"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return &traceLog{
			log: logger,
		}
	},
}

// NewTraceLog 创建追踪日志
func NewTraceLog(traceid string) (l ILogger) {
	return newTraceLog(traceid)
}

//  创建追踪日志
func newTraceLog(traceid string) (l *traceLog) {
	tmp := pool.Get()
	if traceid == "" {
		traceid = util.NewV4().String()
	}
	if l = tmp.(*traceLog); l != nil {
		l.traceID = traceid
	}
	return l
}

// ReleaseTraceLog 释放追踪日志
func ReleaseTraceLog(l ILogger) {
	if tl, ok := l.(*traceLog); ok {
		releaseTraceLog(tl)
	}
}

// releaseTraceLog 释放追踪日志
func releaseTraceLog(l *traceLog) {
	pool.Put(l)
}

func wrap(uid, traceid string, v ...interface{}) (tmp []interface{}) {
	tmp = append(tmp, uid, traceid)
	tmp = append(tmp, v...)
	return
}
func wrapf(uid, traceid string, format string) (f string) {
	f = fmt.Sprintf("[uid:%v][tid:%v]%v", uid, traceid, format)
	return
}

// WrapTraceID 向日志中包装追踪ID
func WrapTraceID(l ILogger, traceid string) {
	if tmp, ok := l.(*traceLog); ok {
		tmp.traceID = traceid
	}
}

type traceLog struct {
	log     ILogger
	uid     string
	traceID string
}

func (l traceLog) SetTraceID(traceid string) {
	l.traceID = traceid
}
func (l traceLog) Debug(v ...interface{}) {
	l.log.Debug(wrap(l.uid, l.traceID, v...))
}
func (l traceLog) Debugf(format string, v ...interface{}) {

	l.log.Debugf(wrapf(l.uid, l.traceID, format), v...)
}
func (l traceLog) Error(v ...interface{}) {
	l.log.Error(wrap(l.uid, l.traceID, v...))
}
func (l traceLog) Errorf(format string, v ...interface{}) {
	l.log.Errorf(wrapf(l.uid, l.traceID, format), v...)
}
func (l traceLog) Info(v ...interface{}) {
	l.log.Info(wrap(l.uid, l.traceID, v...))
}
func (l traceLog) Infof(format string, v ...interface{}) {
	l.log.Infof(wrapf(l.uid, l.traceID, format), v...)
}
func (l traceLog) Warn(v ...interface{}) {
	l.log.Warn(wrap(l.uid, l.traceID, v...))
}
func (l traceLog) Warnf(format string, v ...interface{}) {
	l.log.Warnf(wrapf(l.uid, l.traceID, format), v...)
}
func (l traceLog) Fatal(v ...interface{}) {
	l.log.Fatal(wrap(l.uid, l.traceID, v...))
}
func (l traceLog) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(wrapf(l.uid, l.traceID, format), v...)
}
func (l traceLog) InputLevel() Level {
	return l.log.InputLevel()
}
func (l traceLog) SetInputLevel(lvl Level) {
	l.log.SetInputLevel(lvl)
}
