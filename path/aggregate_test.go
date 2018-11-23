package path

import (
	"reflect"
	"testing"
)

func Test_parseAggregate(t *testing.T) {
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
			"1", args{r: "$aggregate[abc]"}, creAggregateItem("abc"), false,
		},
		{
			"2", args{r: "$aggregate[abc].$aggregate[d]"}, nil, true,
		},
		{
			"3", args{r: "$ss[abc].$r[d]"}, nil, false,
		},
		{
			"4", args{r: "$aggregate[abc"}, nil, true,
		},
		{
			"5", args{r: "$aggregate[abc.a]"}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAggregate(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAggregate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
