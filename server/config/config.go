package config

import (
	"bufio"
	"os"
	"strings"

	"server/configmodel"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  configmodel.DatabaseConfig  `yaml:"database"`
	App       configmodel.AppConfig       `yaml:"app"`
	JWT       configmodel.JwtConfig       `yaml:"jwt"`
	Upload    configmodel.UploadConfig    `yaml:"upload"`
	RateLimit configmodel.RateLimitConfig `yaml:"rate_limit"`
	Redis     configmodel.RedisConfig     `yaml:"redis"`
}

var Global *Config

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func applyEnvOverrides() {
	if v := os.Getenv("DB_HOST"); v != "" {
		Global.Database.Host = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		Global.Database.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		Global.JWT.Secret = v
	}
	if v := os.Getenv("WX_APP_ID"); v != "" {
		Global.App.WxAppID = v
	}
	if v := os.Getenv("WX_APP_SECRET"); v != "" {
		Global.App.WxAppSecret = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		Global.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		Global.Redis.Password = v
	}
}

func validateCORS(env string, origins []string) {
	if env == "dev" {
		return
	}
	for _, origin := range origins {
		if origin == "*" {
			logrus.Fatalf("配置错误：生产环境不允许 CORS 通配符 '*'")
		}
		if strings.Contains(strings.ToLower(origin), "localhost") {
			logrus.Fatalf("配置错误：生产环境不允许 CORS origin 含 'localhost': %s", origin)
		}
	}
}

func validateProductionConfig() {
	if Global.App.Env == "dev" {
		return
	}

	isProd := Global.App.Env == "prod"
	hasFatal := false

	if Global.JWT.Secret == "" {
		logrus.Errorf("生产环境配置错误: JWT Secret 不能为空")
		hasFatal = true
	}
	if Global.Database.Password == "" {
		logrus.Errorf("生产环境配置错误: 数据库密码不能为空")
		hasFatal = true
	}

	if len(Global.App.AllowedOrigins) == 0 {
		logrus.Warnf("生产环境配置警告: AllowedOrigins 为空，仅允许同源请求")
	}
	for _, origin := range Global.App.AllowedOrigins {
		if strings.HasPrefix(origin, "http://") {
			logrus.Warnf("生产环境配置警告: 建议使用 HTTPS origin，当前: %s", origin)
		}
	}

	if isProd && hasFatal {
		logrus.Fatalf("生产环境配置校验失败，已终止启动")
	}
}

func Init(path string) {
	loadEnvFile(".env")

	if env := os.Getenv("APP_ENV"); env != "" {
		envPath := "config." + env + ".yaml"
		if _, err := os.Stat(envPath); err == nil {
			path = envPath
			logrus.Infof("APP_ENV=%s，使用配置文件: %s", env, path)
		}
	}

	var err error
	Global, err = Load(path)
	if err != nil {
		logrus.Fatalf("加载配置失败: %v", err)
	}

	applyEnvOverrides()
	validateCORS(Global.App.Env, Global.App.AllowedOrigins)
	validateProductionConfig()
}
