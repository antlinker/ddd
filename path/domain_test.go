package path

import (
	"testing"
)

func TestNewDomainPath(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want Path
	}{
		{"1", args{id: "abc"}, FromString("$d[abc]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDomainPath(tt.args.id); !got.Equals(tt.want) {
				t.Errorf("NewDomainPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDomain(t *testing.T) {
	p := NewDomainPath("abc")
	testPath(p, "$d[abc]", t)

	di := NewItem(Domain, "abcd")
	p.Append(di)
	testPath(p, "$d[abc.abcd]", t)

	ari := NewItem(AggregateRoot, "users")
	p.Append(ari)

	testPath(p, "$d[abc.abcd].$ar[users]", t)
	ar := NewItem(Aggregate, "AB0001234")
	p.Append(ar)
	testPath(p, "$d[abc.abcd].$ar[users].$a[AB0001234]", t)

}

func TestNewDomain1(t *testing.T) {
	p := NewDomainPath("abc")
	testPath(p, "$d[abc]", t)

	di := NewItem(Domain, "abcd")
	p.Append(di)
	testPath(p, "$d[abc.abcd]", t)

	ari := NewItem(AggregateRoot, "users")
	p.Append(ari)

	testPath(p, "$d[abc.abcd].$ar[users]", t)
	ar := NewItem(Entity, "AB0001234")
	p.Append(ar)
	testPath(p, "$d[abc.abcd].$ar[users].$e[AB0001234]", t)

}

func TestNewDomain2(t *testing.T) {
	p := NewDomainPath("abc")
	testPath(p, "$d[abc]", t)

	di := NewItem(Domain, "abcd")
	p.Append(di)
	testPath(p, "$d[abc.abcd]", t)

	svc := NewItem(Service, "users")
	p.Append(svc)

	testPath(p, "$d[abc.abcd].$s[users]", t)

}

func TestNewDomain3(t *testing.T) {
	p := NewDomainPath("abc")
	testPath(p, "$d[abc]", t)

	di := NewItem(Domain, "abcd")
	p.Append(di)
	testPath(p, "$d[abc.abcd]", t)

	ari := NewItem(AggregateRoot, "users")
	p.Append(ari)

	testPath(p, "$d[abc.abcd].$ar[users]", t)
	ar := NewItem(Repository, "repo")
	p.Append(ar)
	testPath(p, "$d[abc.abcd].$ar[users].$repo[repo]", t)

}
func TestNewDomain4(t *testing.T) {
	p := NewDomainPath("abc")
	testPath(p, "$d[abc]", t)

	ar := NewItem(Repository, "repo")
	p.Append(ar)
	testPath(p, "$d[abc].$repo[repo]", t)

}
func testPath(p Path, want string, t *testing.T) bool {
	if p.Path() != want {
		t.Errorf("NewDomainItem() = %v, want %v", p.Path(), want)
		return false
	}
	return true
}
