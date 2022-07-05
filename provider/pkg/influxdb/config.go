package influxdb

type Config struct {
	Enable     bool   `yaml:"enable"`
	Address    string `yaml:"address"`
	Port       int    `yaml:"port"`
	UDPAddress string `yaml:"udp_address"`
	Database   string `yaml:"database"`
	// precision: n, u, ms, s, m or h
	Precision           string `yaml:"precision"`
	UserName            string `yaml:"username"`
	Password            string `yaml:"password"`
	MaxIdleConns        int    `yaml:"max-idle-conns"`
	MaxIdleConnsPerHost int    `yaml:"max-idle-conns-per-host"`
	IdleConnTimeout     int    `yaml:"idle-conn-timeout"`
}

// CustomConfig Custom Configuration
type CustomConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Address    string `yaml:"address"`
	Port       int    `yaml:"port"`
	UDPAddress string `yaml:"udp_address"`
	Database   string `yaml:"database"`
	// precision: n, u, ms, s, m or h
	Precision           string `yaml:"precision"`
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
