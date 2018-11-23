package path

import (
	"github.com/antlinker/ddd/log"

	"github.com/pkg/errors"
)

const (
	domainPre             = "$domain"
	domainAliasLen        = 2
	aggregateRootPre      = "$aroot"
	aggregateRootAliasLen = 3
	aggregatePre          = "$aggregate"
	aggregateAliasLen     = 2
	servicePre            = "$service"
	serviceAliasLen       = 2
	repositoryPre         = "$repository"
	repositoryAliasLen    = 5
	entityPre             = "$entity"
	entityAliasLen        = 2

	fmtDomainPath        = "$d[%v]"
	fmtAggregateRootPath = "$ar[%v]"
	fmtAggregatePath     = "$a[%v]"
	fmtRepositoryPath    = "$repo[%v]"
	fmtServicePath       = "$s[%v]"
	fmtEntityPath        = "$e[%v]"
)

// ItemKind 路径组成类型
type ItemKind uint8

const (
	// Domain 路径类型 领域
	Domain ItemKind = iota
	// Service 路径类型 领域服务
	Service
	// AggregateRoot 路径类型 领域聚合根
	AggregateRoot
	// Aggregate 路径类型 领域聚合
	Aggregate
	// Entity 路径类型 领域实例
	Entity
	// Repository 路径类型 领域仓库
	Repository
)

// Name 名字
func (k ItemKind) Name() string {
	switch k {
	case Domain:
		return "domain"
	case Service:
		return "service"
	case AggregateRoot:
		return "aggregateRoot"
	case Aggregate:
		return "aggregate"
	case Repository:
		return "repository"
	case Entity:
		return "entity"
	}
	return "未定义"
}

type parsehandler func(path string) (bool, error)

// Path 路径
type Path interface {
	Parse(path string) error
	IsInvalid() bool
	Path() string
	Next() Item
	Equals(Path) bool
	SetRoot(Item)
	Append(Item)
	Clone() Path
	Last() Item
}

// FromString 从字符串创建路径
func FromString(path string) Path {
	if path == "" {
		return &_path{invalid: true}
	}
	p := &_path{}
	err := p.Parse(path)
	if err != nil {
		log.Warnf("解析错误:%v\n", err)
	}
	return p
}

type _path struct {
	root    Item
	invalid bool
	path    string
}

func (d *_path) Next() Item {
	return d.root
}

func (d *_path) SetRoot(i Item) {
	d.root = i
}
func (d *_path) Last() Item {
	ii := d.root
	for ii.Next() != nil {
		ii = ii.Next()
	}
	return ii
}
func (d *_path) Append(i Item) {
	ii := d.root
	for ii.Next() != nil {
		ii = ii.Next()
	}
	ii.Append(i)
	d.path = i.Path()
}
func (d *_path) Path() string {
	if d.root == nil || d.invalid {
		return ""
	}

	i := d.root
	for i.Next() != nil {
		i = i.Next()
	}
	d.path = i.Path()
	return d.path
}

func (d *_path) Equals(p Path) bool {
	return d.Path() == p.Path()
}

func (d *_path) IsInvalid() bool {
	return d.invalid
}
func (d *_path) parse(path string, hs ...parsehandler) error {
	for _, h := range hs {
		if ok, err := h(path); err != nil || ok {
			if err != nil {
				d.invalid = true
				return err
			}
			return nil
		}
	}
	d.invalid = true
	return errors.Errorf("不能解析路径（%v）", path)
}
func (d *_path) Parse(path string) error {
	d.root = nil

	return d.parse(path,
		d.parseDomain,
		d.parseAggregateRoot,
		d.parseRepository,
		d.parseService,
		d.parseEntity,
		d.parseAggregate,
	)

}
func (d *_path) Clone() Path {
	p := d.Path()
	return FromString(p)
}
func (d *_path) parseEntity(path string) (bool, error) {

	entity, err := parseEntity(path)
	if err != nil {
		return true, err
	}
	if entity != nil {
		entity.parent = d.root
		d.root = entity
		_ = entity.Path()

		return true, nil
	}

	return false, nil

}

func (d *_path) parseService(path string) (bool, error) {

	service, err := parseService(path)
	if err != nil {
		return true, err
	}
	if service != nil {
		service.parent = d.root
		d.root = service
		_ = service.Path()

		return true, nil
	}

	return false, nil

}

func (d *_path) parseAggregate(path string) (bool, error) {

	agg, err := parseAggregate(path)
	if err != nil {
		return true, err
	}
	if agg != nil {
		agg.parent = d.root
		d.root = agg
		_ = agg.Path()

		return true, nil
	}

	return false, nil

}

func (d *_path) parseRepository(path string) (bool, error) {

	repo, err := parseRepository(path)
	if err != nil {
		return true, err
	}
	if repo != nil {
		repo.parent = d.root
		d.root = repo
		_ = repo.Path()
		return true, nil
	}

	return false, nil

}
func (d *_path) parseAggregateRoot(path string) (bool, error) {

	ars, r, err := parseAggregateRoot(path)
	if err != nil {
		return true, err
	}

	if ars != nil {
		ars.parent = d.root
		d.root = ars

		ars.Path()
		if r != "" {
			return true, ars.parse(r)
		}
		return true, nil
	}

	return false, nil

}
func (d *_path) parseDomain(path string) (bool, error) {
	f, e, r, err := parseDomain(path)
	if err != nil {
		return true, err
	}
	if f != nil {
		d.root = f
		if r != "" {
			return true, e.parse(r)
		}
		return true, nil
	}
	return false, nil
}

// $domian[a]
// $d[a.b.c]
