package path

import (
	"reflect"
	"testing"
)

func Test_parseAggregateRoot(t *testing.T) {
	type args struct {
		r string
	}
	tests := []struct {
		name    string
		args    args
		want    *aggregateRootItem
		want1   string
		wantErr bool
	}{
		{
			"1", args{r: "$aroot[abc]"}, creAggregateRootItem("abc"), "", false,
		},
		{
			"2", args{r: "$aroot[abc].$aroot[d]"}, creAggregateRootItem("abc"), "$aroot[d]", false,
		},
		{
			"3", args{r: "$ss[abc].$r[d]"}, nil, "$ss[abc].$r[d]", false,
		},
		{
			"4", args{r: "$aroot[abc"}, nil, "", true,
		},
		{
			"5", args{r: "$aroot[abc.a]"}, nil, "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseAggregateRoot(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAggregateRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAggregateRoot() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseAggregateRoot() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
