package path

import (
	"reflect"
	"testing"
)

func Test_parseEntity(t *testing.T) {
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
			"1", args{r: "$e[abc]"}, creEntityItem("abc"), false,
		},
		{
			"2", args{r: "$e[abc].$e[d]"}, nil, true,
		},
		{
			"3", args{r: "$ss[abc].$r[d]"}, nil, false,
		},
		{
			"4", args{r: "$e[abc"}, nil, true,
		},
		{
			"5", args{r: "$e[abc.a]"}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseEntity(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
