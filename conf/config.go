package conf

import (
	"flag"
	"github.com/spf13/viper"
	"my/blog-backend/lib/log"
)

var (
	file string
	C    *Config
)

type Config struct {
	ServerPort string
	DB         *DB
}

type DB struct {
	Name            string
	URL             string
	ConnMaxLifeTime int
	MaxIdleConns    int
}

func init() {
	flag.StringVar(&file, "f", "conf/config.yaml", "config file path")
	flag.Parse()
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
	log.Info("配置初始化完成")
}
