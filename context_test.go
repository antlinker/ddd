package ddd

import (
	"context"
	"testing"
)

func TestNewContext(t *testing.T) {
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{ctx: context.Background(), uid: "abc"}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.args.ctx, tt.args.uid); got.UID() != tt.want {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestContextValue(t *testing.T) {
	uid := "aa"
	ctx := NewContext(context.Background(), uid)

	if ctx.UID() != uid {
		t.Errorf("ctx.UID() = %v, want %v", ctx.UID(), uid)
		return
	}
	key := "1111"
	value := "2222"
	ctx.Put(key, value)

	v, ok := ctx.Get(key)
	if !ok {
		t.Errorf("v, ok := ctx.Get(key) ok = %v, want %v", ok, true)
		return
	}
	v1, ok := v.(string)
	if !ok {
		t.Errorf("v1,ok:= v.(string) ok = %v, want %v", ok, true)
		return
	}
	if v1 != value {
		t.Errorf("v, ok := ctx.Get(key) v = %v, want %v", v1, value)
		return
	}
}
