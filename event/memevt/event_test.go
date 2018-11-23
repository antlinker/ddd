package memevt

import (
	"testing"
	"time"

	"github.com/antlinker/ddd"
)

func TestEvt(t *testing.T) {
	e := NewEventTrigger()
	r := make(chan struct{})
	e.AddListener("a", "a", func(evt ddd.Event) {
		t.Log(evt)
		close(r)
	})

	e.TriggerEvent(ddd.Event{
		Type:   "a",
		Action: "a",
		ID:     "1",
	})
	select {
	case <-r:
	case <-time.After(time.Millisecond):
		t.Error("等待事件触发超时")
	}
}
