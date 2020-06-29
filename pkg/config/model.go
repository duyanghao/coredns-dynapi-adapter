package config

type Config struct {
	Server  *Server  `yaml:"server"`
	Gin     *Gin     `yaml:"gin"`
	Log     *Log     `yaml:"log"`
	Coredns *Coredns `yaml:"coredns"`
}

type Server struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

type Coredns struct {
	CorefilePath string `yaml:"corefilePath"`
	ZonesDir     string `yaml:"zonesDir"`
}

type Gin struct {
	Mode string `yaml:"mode"`
}

type Log struct {
	Level string `yaml:"level"`
}
