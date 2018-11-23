package ddd

import (
	"fmt"

	"github.com/antlinker/ddd/path"
)

var (
	_ DomainNode = &node{}
)

type node struct {
	path     path.Path
	domainID string
	pathItem path.Item
	parent   DomainNode
	children map[path.ItemKind]map[string]DomainNode
}

func (n *node) isInit() bool {
	return n.pathItem != nil
}
func (n *node) Init(self DomainNode, parent DomainNode, DomainID string) {
	panic(fmt.Sprintf("不能直接调用该方法初始化节点"))
}
func (n *node) init(self DomainNode, parent DomainNode, kind path.ItemKind, domainID string, noisparent bool) {
	n.children = make(map[path.ItemKind]map[string]DomainNode)
	n.domainID = domainID
	n.pathItem = path.NewItem(kind, domainID)
	n.domainID = domainID
	if parent != nil {
		n.setParent(parent)
		if !noisparent {
			parent.appendChildren(self)
		}
	}
}
func (n *node) ItemKind() path.ItemKind {
	return n.pathItem.Kind()
}
func (n *node) getChildren() (dn map[path.ItemKind]map[string]DomainNode) {
	return n.children
}
func (n *node) getNodes(k path.ItemKind) (dn map[string]DomainNode, exists bool) {
	dn, exists = n.children[k]
	return
}
func (n *node) getNode(k path.ItemKind, id string) (dn DomainNode, exists bool) {
	if cs, ok := n.children[k]; ok {
		dn, exists = cs[id]
	}
	return
}
func (n *node) appendChildren(c DomainNode) {
	n.appNode(c.ItemKind(), c)
}
func (n *node) appNode(key path.ItemKind, c DomainNode) {
	cs, ok := n.children[key]
	if !ok {
		cs = make(map[string]DomainNode)
		n.children[key] = cs
	}
	cs[c.DomainID()] = c

}

// Domain 获取所在域
func (n *node) Domain() Domain {
	for p := n.parent; p != nil; p = p.Parent() {
		switch d := p.(type) {
		case Domain:
			return d
		}
	}
	return nil
}

func (n node) DomainPath() path.Path {
	return n.path
}
func (n node) Parent() DomainNode {
	return n.parent
}

// func (n *node) setPath(p path.Path) {
// 	n.path = p
// }
func (n node) DomainID() string {
	return n.domainID
}
func (n *node) setParent(parent DomainNode) {
	if n.path != nil && !n.path.IsInvalid() {
		panic(fmt.Sprintf("已经设置过上级的领域节点，不能被再次设置上级节点：%v", n.domainID))
	}
	n.resetParent(parent)
}
func (n *node) resetParent(parent DomainNode) {
	// if n.path != nil && !n.path.IsInvalid() {
	// 	panic(fmt.Sprintf("已经设置过上级的领域节点，不能被再次设置上级节点：%v", n.domainID))
	// }
	if parent.DomainPath() == nil || parent.DomainPath().IsInvalid() {
		panic(fmt.Sprintf("无效的上级领域节点：%v", parent.DomainID()))
	}
	n.parent = parent

	n.updatePath()
}
func (n *node) updatePath() {
	if n.parent == nil {
		return
	}
	n.path = n.parent.DomainPath().Clone()
	n.path.Append(n.pathItem)
}
func (n *node) Trigger(c Context, etype, action string, data interface{}) {
	c.Trigger(etype, action, n.path.Path(), data)
}
