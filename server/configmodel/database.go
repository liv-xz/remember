package configmodel

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
	Source   string `yaml:"source"`
}
