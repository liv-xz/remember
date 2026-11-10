package configmodel

type RateLimitConfig struct {
	Enable       bool `yaml:"enable"`
	AuthRate     int  `yaml:"auth_rate"`
	AuthBurst    int  `yaml:"auth_burst"`
	GeneralRate  int  `yaml:"general_rate"`
	GeneralBurst int  `yaml:"general_burst"`
}
