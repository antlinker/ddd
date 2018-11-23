package path

import (
	"reflect"
	"testing"
)

func TestNewItem(t *testing.T) {
	type args struct {
		kind ItemKind
		id   string
	}
	tests := []struct {
		name string
		args args
		want Item
	}{
		// TODO: Add test cases.
		{"1", args{kind: 10, id: "aaa"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewItem(tt.args.kind, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
