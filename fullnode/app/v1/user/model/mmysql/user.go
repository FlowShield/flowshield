package mmysql

type User struct {
	ID        uint   `gorm:"primarykey"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	UUID      string `json:"uuid" gorm:"column:uuid"`
}

func (User) TableName() string {
	return "zta_user"
}
