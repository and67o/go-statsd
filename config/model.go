package config

type StatsDConfig struct {
	Host string
	Port int
}

type Config struct {
	Stats StatsDConfig
}
