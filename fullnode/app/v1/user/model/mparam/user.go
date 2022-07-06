package mparam

type BindWallet struct {
	Wallet string `json:"wallet" form:"wallet" binding:"required"`
}
