package dto

// Pagination 分页请求参数
type Pagination struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// DefaultPagination 默认分页
func DefaultPagination() Pagination {
	return Pagination{Page: 1, PageSize: 20}
}

// Offset 计算偏移量
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
