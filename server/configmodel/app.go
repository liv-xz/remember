package configmodel

type AppConfig struct {
	Name           string   `yaml:"name"`
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	Env            string   `yaml:"env"`
	GinMode        string   `yaml:"gin_mode"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	WxAppID        string   `yaml:"wx_app_id"`
	WxAppSecret    string   `yaml:"wx_app_secret"`
}
