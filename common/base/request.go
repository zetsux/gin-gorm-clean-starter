package base

type GetsRequest struct {
	Search string `json:"search" form:"search"`
	Sort   string `json:"sort" form:"sort"`
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
}
