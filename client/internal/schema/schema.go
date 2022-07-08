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

// ControlCommonResult
type ControlCommonResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ControlPaginate
type ControlPaginate struct {
	Total    int `json:"total"`
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
}
