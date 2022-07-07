package mmysql

type Resource struct {
	ID        uint   `gorm:"primarykey"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
	Name      string `json:"name"`
	UUID      string `json:"uuid" gorm:"column:uuid"`
	UserUUID  string `json:"user_uuid" gorm:"user_uuid"`
	Type      string `json:"type"`
	Host      string `json:"host"`                    // api.github.com
	Port      string `json:"port" gorm:"column:port"` // 80-443;3306;6379
	Cid       string `json:"cid"`
}

func (Resource) TableName() string {
	return "zta_resource"
}
