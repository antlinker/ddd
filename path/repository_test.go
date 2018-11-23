package path

import (
	"reflect"
	"testing"
)

func Test_parseRepository(t *testing.T) {
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
			"1", args{"$repo[abc]"}, creRepositoryItem("abc"), false,
		},
		{
			"2", args{"$repo[abc].$r[d]"}, nil, true,
		},
		{
			"3", args{"$ss[abc].$r[d]"}, nil, false,
		},
		{
			"4", args{"$repo[abc"}, nil, true,
		},
		{
			"5", args{"$repo[abc.a]"}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRepository(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}
