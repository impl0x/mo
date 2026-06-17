package config


type ServerConfig struct {
	Host string
	Port uint16
	ReadTimeout int
	WriteTimeout int
}
