package path

import (
	"reflect"
	"testing"
)

func Test_parseService(t *testing.T) {
	type args struct {
		r string
	}
	tests := []struct {
		name    string
		args    args
		want    *pathItem
		wantErr bool
	}{
		{
			"1", args{"$s[abc]"}, creServiceItem("abc"), false,
		},
		{
			"2", args{"$s[abc].$r[d]"}, nil, true,
		},
		{
			"3", args{"$r[abc].$r[d]"}, nil, false,
		},
		{
			"4", args{"$s[abc"}, nil, true,
		},
		{
			"5", args{"$s[abc.a]"}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseService(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseService() got = %v, want %v", got, tt.want)
			}
		})
	}
}
