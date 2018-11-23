package path

import (
	"fmt"

	"github.com/antlinker/ddd/log"
)

// Item 路径组成元素
type Item interface {
	CurName() string
	Next() Item
	Parent() Item
	Kind() ItemKind
	Path() string
	Append(Item)
	Data() interface{}
	SetData(data interface{})
	Equal(item Item) bool
	setNext(Item)
	setParent(Item)
	resetPath() string
}

// NewItem 创建Item
func NewItem(kind ItemKind, id string) Item {
	switch kind {
	case Domain:
		return creDomainItem(id)
	case AggregateRoot:
		return creAggregateRootItem(id)
	case Aggregate:
		return creAggregateItem(id)
	case Service:
		return creServiceItem(id)
	case Entity:
		return creEntityItem(id)
	case Repository:
		return creRepositoryItem(id)
	}
	return nil
}

type pathItem struct {
	next    Item
	parent  Item
	curName string
	kind    ItemKind
	path    string
	fmtPath string
	data    interface{}
}

func (i *pathItem) Equal(item Item) bool {
	if i.kind == item.Kind() &&
		i.curName == item.CurName() &&
		i.parent.Equal(item.Parent()) {
		return true
	}
	return false
}
func (i *pathItem) CurName() string {
	return i.curName
}
func (i *pathItem) Kind() ItemKind {
	return i.kind
}
func (i *pathItem) Next() Item {
	return i.next
}
func (i *pathItem) Parent() Item {
	return i.parent
}
func (i *pathItem) Data() interface{} {
	return i.data
}
func (i *pathItem) SetData(data interface{}) {
	i.data = data
}
func (i *pathItem) setNext(in Item) {
	i.next = in
}
func (i *pathItem) setParent(in Item) {
	i.parent = in
}

func (i *pathItem) Path() string {
	if i.path != "" {
		return i.path
	}
	return i.resetPath()
}
func (i *pathItem) resetPath() string {
	i.path = fmt.Sprintf(i.fmtPath, i.curName)
	if i.parent != nil {
		i.path = i.parent.Path() + string(pathSpe) + i.path
	}
	return i.path
}
func (i *pathItem) Append(in Item) {
	log.Error("pathItem禁止添加item")
}
