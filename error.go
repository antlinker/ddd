package ddd

import "errors"

var (
	// ErrorNoFoundPath 路径找不到对应的领域对象
	ErrorNoFoundPath = errors.New("路径找不到对应的领域对象")
	// ErrorNodeKindNotMatch 路径对应的节点类型不匹配
	ErrorNodeKindNotMatch = errors.New("路径对应的节点类型不匹配")
)
