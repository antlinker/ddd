package ddd

import (
	"os"
	"testing"

	"github.com/antlinker/ddd/path"
)

func creDomain(parent DomainNode, id string) Domain {
	d := &BaseDomain{}
	d.Init(d, parent, id)

	return d
}

type tlog interface {
	Fatalf(fmtstr string, v ...interface{})
}

func checkDomainPath(p path.Path, wantpath string, t tlog) {
	if p == nil {
		t.Fatalf("RegDomain DomainPath 想要一个 %s 实际得到 nil", wantpath)
		return
	}
	if p.Path() != wantpath {
		t.Fatalf("RegDomain DomainPath 想要一个 %s 实际得到%s", wantpath, p.Path())
		return
	}
}
func creAggregateRoot(parent DomainNode, id string) AggregateRoot {
	d := &BaseAggregateRoot{}
	d.Init(d, parent, id)

	return d
}
func creRepository(parent DomainNode, id string) Repository {
	d := &BaseRepository{}
	d.Init(d, parent, id)

	return d
}
func creService(parent DomainNode, id string) DomainService {
	d := &BaseService{}
	d.Init(d, parent, id)

	return d
}

var (
	rootDomainID = "abc"
	d            = creDomain(nil, rootDomainID)
)

func TestMain(t *testing.M) {

	RegDomain(d)

	os.Exit(t.Run())
}

// func TestExport(t *testing.T) {
// 	p := d.DomainPath()
// 	checkDomainPath(p, "$d[abc]", t)
// 	d1 := creDomain(nil, "bcd")
// 	RegSubDomain(d, d1)
// 	checkDomainPath(d1.DomainPath(), "$d[abc.bcd]", t)
// 	a := creAggregateRoot(nil, "eee")
// 	RegAggregateRoot(d, a)
// 	checkDomainPath(a.DomainPath(), "$d[abc].$ar[eee]", t)
// 	r := creRepository(nil, "info")
// 	RegRepositoryByDomain(d, r)
// 	checkDomainPath(r.DomainPath(), "$d[abc].$repo[info]", t)

// 	d2 := creDomain(nil, "bcd2")
// 	checkDomainPath(d2.DomainPath(), "$d[bcd2]", t)
// 	r2 := creRepository(nil, "repo2")
// 	RegRepositoryByDomain(d2, r2)
// 	checkDomainPath(r2.DomainPath(), "$d[bcd2].$repo[repo2]", t)

// 	RegSubDomain(d, d2)
// 	checkDomainPath(d2.DomainPath(), "$d[abc.bcd2]", t)
// 	checkDomainPath(r2.DomainPath(), "$d[abc.bcd2].$repo[repo2]", t)

// 	r3 := creRepository(nil, "repo3")
// 	SetRepoForARoot(a, r3)
// 	checkDomainPath(r3.DomainPath(), "$d[abc].$ar[eee].$repo[repo3]", t)
// 	svc1 := creService(nil, "svc1")
// 	RegService(d2, svc1)
// 	checkDomainPath(svc1.DomainPath(), "$d[abc.bcd2].$s[svc1]", t)
// 	testPanic("testsv", t, func() {
// 		RegService(d, svc1)
// 	})
// }
func TestInit(t *testing.T) {
	p := d.DomainPath()
	checkDomainPath(p, "$d[abc]", t)
	d1 := creDomain(d, "bcd")
	checkDomainPath(d1.DomainPath(), "$d[abc.bcd]", t)
	a := creAggregateRoot(d, "eee")
	checkDomainPath(a.DomainPath(), "$d[abc].$ar[eee]", t)
	r := creRepository(d, "info")
	checkDomainPath(r.DomainPath(), "$d[abc].$repo[info]", t)

	d2 := creDomain(d, "bcd2")
	checkDomainPath(d2.DomainPath(), "$d[abc.bcd2]", t)
	r2 := creRepository(d2, "repo2")
	checkDomainPath(r2.DomainPath(), "$d[abc.bcd2].$repo[repo2]", t)

	r3 := creRepository(a, "repo3")
	checkDomainPath(r3.DomainPath(), "$d[abc].$ar[eee].$repo[repo3]", t)
	svc1 := creService(d2, "svc1")
	checkDomainPath(svc1.DomainPath(), "$d[abc.bcd2].$s[svc1]", t)
	// testPanic("testsv", t, func() {
	// 	RegService(d, svc1)
	// })
}
func testPanic(name string, t *testing.T, handler func()) {
	t.Run(name, func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatalf("%s 此处期望获取到panic", name)
			}
		}()
		handler()
	})
}
