package service

type Config struct {
	Datapath   string
	Staticpath string
	Basepath   string
}

func GetConfig() *Config {
	return &Config{Datapath: "data", Staticpath: "static", Basepath: "http://localhost:8080"}
}
