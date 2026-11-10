package flags

import (
	"flag"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagOptions = &Options{}

func Parse1() {
	flag.StringVar(&FlagOptions.File, "f", "config.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "db", false, "是否初始化数据库")
	flag.BoolVar(&FlagOptions.Version, "v", false, "显示版本信息")
	flag.Parse()
}

// 只导出参数值，不包含业务逻辑
func ShouldInitDB() bool {
	return FlagOptions.DB
}
