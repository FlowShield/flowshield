package mparam

type MachineOauth struct {
	Machine string `json:"machine"`
}

type MachineLongPoll struct {
	Category string `json:"category" form:"category" binding:"required"`
	Timeout  int    `json:"timeout" form:"timeout" binding:"required"`
}
