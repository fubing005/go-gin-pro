package request

type PageQuery struct {
	Page     int `form:"page,default=1" json:"page" binding:"omitempty"`
	PageSize int `form:"page_size,default=10" json:"page_size" binding:"omitempty"`
}
