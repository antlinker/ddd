package path

import (
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantpath string
		invalid  bool
	}{
		{"1", args{"$domain[ aaa.bb ]"}, "$d[aaa.bb]", false},
		{"2", args{"$d[ aaa.bb ]"}, "$d[aaa.bb]", false},
		{"3", args{"$domain[ aaa.bb ].$service[ ccc ]"}, "$d[aaa.bb].$s[ccc]", false},
		{"4", args{"$domain[ aaa.bb ].$repo[ ccc ]"}, "$d[aaa.bb].$repo[ccc]", false},
		{"5", args{"$domain[ aaa.bb ].$aroot[ ccc ]"}, "$d[aaa.bb].$ar[ccc]", false},
		{"6", args{"$domain[ aaa.bb ].$aroot[ ccc ].$a[ id ]"}, "$d[aaa.bb].$ar[ccc].$a[id]", false},
		{"7", args{"$domain[ aaa.bb ].$aroot[ ccc ].$repo[ id ]"}, "$d[aaa.bb].$ar[ccc].$repo[id]", false},
		{"8", args{"$domain[ aaa.bb ].$entity[ ccc:111 ]"}, "$d[aaa.bb].$e[ccc:111]", false},

		{"10", args{"$service[ ccc ]"}, "$s[ccc]", false},
		{"11", args{"$repo[ ccc ]"}, "$repo[ccc]", false},
		{"12", args{"$aroot[ ccc ]"}, "$ar[ccc]", false},
		{"13", args{"$aroot[ ccc ].$repo[ id ]"}, "$ar[ccc].$repo[id]", false},
		{"14", args{"$aroot[ ccc ].$a[ id ]"}, "$ar[ccc].$a[id]", false},
		{"15", args{"$entity[ ccc:111 ]"}, "$e[ccc:111]", false},
		{"16", args{"$a[ ccc:111 ]"}, "$a[ccc:111]", false},

		{"err-0", args{""}, "", true},

		{"err-1", args{"$domain[ aaa.bb "}, "", true},
		{"err-2", args{"$d[ aaa.b b ]"}, "", true},
		{"err-3", args{"$domain[ aaa.bb ].$service[ ccc "}, "", true},
		{"err-4", args{"$domain[ aaa.bb ].$repo[ ccc "}, "", true},
		{"err-5", args{"$domain[ aaa.bb ].$aroot[ ccc "}, "", true},
		{"err-6", args{"$domain[ aaa.bb ].$aroot[ ccc ].$a[ id "}, "", true},
		{"err-7", args{"$domain[ aaa.bb ].$aroot[ ccc ].$repo[ id "}, "", true},
		{"err-8", args{"$domain[ aaa.bb ].$entity[ ccc: 111 ]"}, "", true},
		{"err-8-1", args{"$domain[ aaa.bb ].$entity[ ccc:111 ].$a[aa]"}, "", true},
		{"err-8-2", args{"$domain[ aaa.bb ].$asdfasd[ ccc:111 ].$a[aa]"}, "", true},

		{"err-10", args{"$service[ ccc "}, "", true},
		{"err-10-1", args{"$service[ ccc ].$a[333]"}, "", true},
		{"err-11", args{"$repo[ ccc "}, "", true},
		{"err-11-1", args{"$repo[ ccc ].$a[11]"}, "", true},
		{"err-12", args{"$aroot[ ccc"}, "", true},
		{"err-13-0", args{"$aroot[ ccc].$abc[123]"}, "", true},
		{"err-13", args{"$aroot[ ccc ].$repo[ id "}, "", true},
		{"err-14", args{"$aroot[ ccc ].$a[ id "}, "", true},
		{"err-15", args{"$entity[ ccc:111 "}, "", true},
		{"err-15", args{"$entity[ ccc:111 ].$ar[22]"}, "", true},
		{"err-16", args{"$a[ ccc:111 "}, "", true},
		{"err-16-1", args{"$a[ ccc:111].$a[11] "}, "", true},
		{"err-17", args{"$eee[ ccc:111]"}, "", true},
		{"err-18", args{"$xxx[ ccc:111]"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromString(tt.args.path)
			if got.Path() != tt.wantpath {
				t.Errorf("FromString().Path = %v, want %v", got.Path(), tt.wantpath)
			}
			if got.IsInvalid() != tt.invalid {
				t.Errorf("FromString().IsInvalid = %v, want %v", got.IsInvalid(), tt.invalid)
			}
		})
	}
}

func TestFromStringItem(t *testing.T) {
	path := "$domain[ aaa.bb ].$service[ ccc ]"
	p := FromString(path)

	testItem(nil, p.Next(), t,
		testItemResult{Domain, "aaa", "$d[aaa]"},
		testItemResult{Domain, "bb", "$d[aaa.bb]"},
		testItemResult{Service, "ccc", "$d[aaa.bb].$s[ccc]"})
}

type testItemResult struct {
	k    ItemKind
	name string
	path string
}

func testItem(parent Item, i Item, t *testing.T, ks ...testItemResult) {
	if i.Kind() != ks[0].k {
		t.Errorf("Item().Kind = %v, want %v", i.Kind(), ks[0])
	}
	if i.CurName() != ks[0].name {
		t.Errorf("Item().CurName = %v, want %v", i.CurName(), ks[0].name)
	}
	if i.Path() != ks[0].path {
		t.Errorf("Item().Path = %v, want %v", i.Path(), ks[0].path)
	}
	if parent == nil && i.Parent() != nil {
		t.Errorf("Item().Parent = %q, want root", i.Parent())
	}
	if parent != nil && i.Parent() != parent {
		t.Errorf("Item().Parent = %q, want %q", i.Parent(), parent)
	}
	if len(ks)-1 > 0 && i.Next() != nil {

		testItem(i, i.Next(), t, ks[1:]...)
	}

}
