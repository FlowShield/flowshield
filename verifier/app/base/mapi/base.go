package mapi

type Paginate struct {
	Total       int `json:"total"`
	TotalPage   int `json:"total_page"`
	CurrentPage int `json:"current_page"`
	PrePage     int `json:"pre_page"`
}

type AdminPaginate struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}
