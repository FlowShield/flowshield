package influxdb

// Config configuration
type Config struct {
	Enable              bool   `yaml:"enable"` //Service switch
	Address             string `yaml:"address"`
	Port                int    `yaml:"port"`
	UDPAddress          string `yaml:"udp_address"` //influxdb UDP address of the database，ip:port
	Database            string `yaml:"database"`    //Database name
	Precision           string `yaml:"precision"`   //Accuracy n, u, ms, s, m or h
	UserName            string `yaml:"username"`
	Password            string `yaml:"password"`
	MaxIdleConns        int    `yaml:"max-idle-conns"`
	MaxIdleConnsPerHost int    `yaml:"max-idle-conns-per-host"`
	IdleConnTimeout     int    `yaml:"idle-conn-timeout"`
}

// CustomConfig Custom configuration
type CustomConfig struct {
	Enabled             bool   `yaml:"enabled"` //Service switch
	Address             string `yaml:"address"`
	Port                int    `yaml:"port"`
	UDPAddress          string `yaml:"udp_address"` //influxdb UDP address of the database，ip:port
	Database            string `yaml:"database"`    //Database name
	Precision           string `yaml:"precision"`   //Accuracy n, u, ms, s, m or h
	UserName            string `yaml:"username"`
	Password            string `yaml:"password"`
	ReadUserName        string `yaml:"read-username"`
	ReadPassword        string `yaml:"read-password"`
	MaxIdleConns        int    `yaml:"max-idle-conns"`
	MaxIdleConnsPerHost int    `yaml:"max-idle-conns-per-host"`
	IdleConnTimeout     int    `yaml:"idle-conn-timeout"`
	FlushSize           int    `yaml:"flush-size"`
	FlushTime           int    `yaml:"flush-time"`
}
