package path

import (
	"reflect"
	"testing"
)

func Test_parseNamespaces(t *testing.T) {
	type args struct {
		pathstr   string
		startflag string
		aliaslen  int
	}
	tests := []struct {
		name          string
		args          args
		wantDomains   []string
		wantRemainder string
		wantOk        bool
	}{
		{"1", args{"$domain[a]", "$domain", 2}, []string{"a"}, "", true},
		{"2", args{"$domain[ a ]", "$domain", 2}, []string{"a"}, "", true},
		{"3", args{"$domain[ a]", "$domain", 2}, []string{"a"}, "", true},
		{"4", args{"$domain[a ]", "$domain", 2}, []string{"a"}, "", true},
		{"5", args{"$domain [a]", "$domain", 2}, []string{"a"}, "", true},
		{"6", args{"$domain [ a]", "$domain", 2}, []string{"a"}, "", true},
		{"7", args{"$domain [ a ]", "$domain", 2}, []string{"a"}, "", true},
		{"8", args{"$domain [a ]", "$domain", 2}, []string{"a"}, "", true},
		{"9", args{"$domain [a ].$abc", "$domain", 2}, []string{"a"}, "$abc", true},
		{"10", args{"$domain [a ].", "$domain", 2}, []string{"a"}, "", true},
		{"11", args{"$domain [a ] .", "$domain", 2}, []string{"a"}, "", true},
		{"12", args{"$domain [a   ]   .", "$domain", 2}, []string{"a"}, "", true},
		{"13", args{"$domain [abcd   ]   .", "$domain", 2}, []string{"abcd"}, "", true},
		{"14", args{"$domain [abcd . 123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"15", args{"$domain [abcd.123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"16", args{"$domain [   abcd .   123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"17", args{"$domain [   abcd .   a   ]   .", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"18", args{"$domain [   abcd .   a   ]   .", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"19", args{"$domain [   abcd .a].", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"20", args{"$domain [  ].", "$domain", 2}, nil, "", true},
		{"21", args{"$domain [].", "$domain", 2}, nil, "", true},
		{"22", args{"$domain[].", "$domain", 2}, nil, "", true},
		{"23", args{"$domain[].abc", "$domain", 2}, nil, "abc", true},
		{"24", args{"$service[aaa]", "$domain", 2}, nil, "$service[aaa]", true},

		{"err-1", args{"$domain ", "$domain", 2}, nil, "", false},
		{"err-2", args{"$domain    abcd .a ", "$domain", 2}, nil, "", false},
		{"err-3", args{"$domain[ab cd .a ]", "$domain", 2}, []string{"ab"}, "", false},
		{"err-4", args{"$ domain[ab cd .a ]", "$domain", 2}, nil, "$ domain[ab cd .a ]", true},
		{"err-5", args{"domain[ab cd .a ]", "$domain", 2}, nil, "domain[ab cd .a ]", true},

		{"s-1", args{"$d[a]", "$domain", 2}, []string{"a"}, "", true},
		{"s-2", args{"$d[ a ]", "$domain", 2}, []string{"a"}, "", true},
		{"s-3", args{"$d[ a]", "$domain", 2}, []string{"a"}, "", true},
		{"s-4", args{"$d[a ]", "$domain", 2}, []string{"a"}, "", true},
		{"s-5", args{"$d [a]", "$domain", 2}, []string{"a"}, "", true},
		{"s-6", args{"$d [ a]", "$domain", 2}, []string{"a"}, "", true},
		{"s-7", args{"$d [ a ]", "$domain", 2}, []string{"a"}, "", true},
		{"s-8", args{"$d [a ]", "$domain", 2}, []string{"a"}, "", true},
		{"s-9", args{"$d [a ].$abc", "$domain", 2}, []string{"a"}, "$abc", true},
		{"s-10", args{"$d [a ].", "$domain", 2}, []string{"a"}, "", true},
		{"s-11", args{"$d [a ] .", "$domain", 2}, []string{"a"}, "", true},
		{"s-12", args{"$d [a   ]   .", "$domain", 2}, []string{"a"}, "", true},
		{"s-13", args{"$d [abcd   ]   .", "$domain", 2}, []string{"abcd"}, "", true},
		{"s-14", args{"$d [abcd . 123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"s-15", args{"$d [abcd.123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"s-16", args{"$d [   abcd .   123   ]   .", "$domain", 2}, []string{"abcd", "123"}, "", true},
		{"s-17", args{"$d [   abcd .   a   ]   .", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"s-18", args{"$d [   abcd .   a   ]   .", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"s-19", args{"$d [   abcd .a].", "$domain", 2}, []string{"abcd", "a"}, "", true},
		{"s-20", args{"$d [  ].", "$domain", 2}, nil, "", true},
		{"s-21", args{"$d [].", "$domain", 2}, nil, "", true},
		{"s-22", args{"$d[].", "$domain", 2}, nil, "", true},
		{"s-23", args{"$d[].abc", "$domain", 2}, nil, "abc", true},

		{"s-err-1", args{"$d ", "$domain", 2}, nil, "", false},
		{"s-err-2", args{"$d    abcd .a ", "$domain", 2}, nil, "", false},
		{"s-err-3", args{"$d[ab cd .a ]", "$domain", 2}, []string{"ab"}, "", false},
		{"s-err-4", args{"$ d[ab cd .a ]", "$domain", 2}, nil, "$ d[ab cd .a ]", true},
		{"s-err-5", args{"d[ab cd .a ]", "$domain", 2}, nil, "d[ab cd .a ]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDomains, gotRemainder, gotOk := parseNamespaces(tt.args.pathstr, tt.args.startflag, tt.args.aliaslen)
			if !reflect.DeepEqual(gotDomains, tt.wantDomains) {
				t.Errorf("parseDomain() gotDomains = %v, want %v", gotDomains, tt.wantDomains)
			}
			if gotRemainder != tt.wantRemainder {
				t.Errorf("parseDomain() gotRemainder = %v, want %v", gotRemainder, tt.wantRemainder)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseDomain() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
