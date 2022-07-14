package mmysql

const (
	UnBind = iota
	Bind
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	UUID      string `json:"uuid" gorm:"column:uuid"`
	//Wallet    string `json:"wallet"`
	Status   int  `json:"status"`            // 0 未绑定 1 已绑定
	Master   bool `json:"master" gorm:"-"`   // 判断当前用户是否Dao主
	Provider bool `json:"provider" gorm:"-"` // 判断当前用户是否Provider
}

func (User) TableName() string {
	return "zta_user"
}
