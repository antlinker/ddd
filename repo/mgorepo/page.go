package mgorepo

// PageInfo 分页查询参数
type PageInfo struct {
	PageSize int         // 分页条数
	PageID   interface{} // 分页id
	Direct   int         // 分页方向 0 下一页，1上一页
	Desc     bool        // true id 降序 ，false  id 升序
	Current  int         // 当前页
	Mode     int         // 分页模式 mode = 0 带页码分页 = 1 带id分页
}

// PageResult 返回分页结果
type PageResult struct {
	Current  int // 当前页
	Total    int // 总数量
	PageSize int // 分页条数
	End      int // 1 最后一页
}
