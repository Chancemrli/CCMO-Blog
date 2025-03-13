package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量
var Conf = new(AppConfig)

// 通用设置
type AppConfig struct {
	MachineID           uint16 `mapstructure:"machine_id"`
	TokenExpireDuration int    `mapstructure:"token_expire_time"`
	Mode                string `mapstructure:"mode"`
	Port                int    `mapstructure:"port"`
	*LogConfig          `mapstructure:"log"`
	*MySQLConfig        `mapstructure:"mysql"`
	*RedisConfig        `mapstructure:"redis"`
}

// 数据库连接设置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"db"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// Redis连接设置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// 日志设置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// 初始化
func Init() error {
	// 设置文件路径
	viper.SetConfigFile("./conf/config.yaml")
	// viper.AddConfigPath("./conf")
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config.yaml has been modified...")
		viper.Unmarshal(Conf)
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}
	return err
}
