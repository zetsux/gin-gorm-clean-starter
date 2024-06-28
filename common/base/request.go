package base

type GetsRequest struct {
	Search  string `json:"search" form:"search"`
	Sort    string `json:"sort" form:"sort"`
	Page    int    `json:"page" form:"page"`
	PerPage int    `json:"per_page" form:"per_page"`
}
