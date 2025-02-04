package service

type Config struct {
	datapath   string
	staticpath string
	basepath   string
}

func GetConfig() *Config {
	return &Config{datapath: "data", staticpath: "static", basepath: "http://localhost:8080"}
}
