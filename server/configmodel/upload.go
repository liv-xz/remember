package configmodel

type UploadConfig struct {
	Dir          string   `yaml:"dir"`
	MaxSize      int64    `yaml:"max_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}
