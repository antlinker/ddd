package mgorepo

import "testing"

func Test_decodeUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"1", args{url: "127.0.0.1"}, "127.0.0.1"},
		{"2", args{url: "mongo://127.0.0.1"}, "mongo://127.0.0.1"},
		{"3", args{url: "mongo://aaa@aaa:127.0.0.1"}, "mongo://*@*:127.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeUrl(tt.args.url); got != tt.want {
				t.Errorf("decodeUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
