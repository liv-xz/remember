package configmodel

type JwtConfig struct {
	Secret        string `yaml:"secret"`
	AccessExpire  int    `yaml:"access_expire"`
	RefreshExpire int    `yaml:"refresh_expire"`
}
