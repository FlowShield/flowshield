package schema

// StatusText Define status text
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

// NextServer
type NextServer struct {
	Host string
	Port string
}

// ControCommonResult
type ControCommonResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ControPaginate
type ControPaginate struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
}
