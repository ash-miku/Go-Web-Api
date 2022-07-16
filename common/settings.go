package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

var CONFIG *App
var vpConfig *viper.Viper

//DbConfig  数据库配置
type DbConfig struct {
	Project   string `yaml:"Project"`
	Driver    string `yaml:"Driver"`
	Name      string `yaml:"Name"`
	Host      string `yaml:"Host"`
	Port      int    `yaml:"Port"`
	User      string `yaml:"User"`
	Password  string `yaml:"Password"`
	Charset   string `yaml:"Charset"`
	Sslmode   string `yaml:"Sslmode"`
	DmsPublic string `yaml:"DmsPublic"` // 与业务不存在任何关系表的存放数据库(errors tenant) 默认为dms_public
	ParseTime string `yaml:"ParseTime"` // ParseTime=true 将 DATE 和 DATETIME 值的输出类型更改为 time
	Loc       string `yaml:"Loc"`       // 设置 time.Time 值的位置（使用时parseTime=true）。“本地”设置系统的位置。有关详细信息，请参阅 https://golang.org/pkg/time/#LoadLocation
}
type ServerConfig struct {
	ControlServer ControlServer `yaml:"ControlServer"` // 内部组件的中间件
	IcarusServer  IcarusServer  `yaml:"IcarusServer"`  // 关系型数据库的中间件
	NosqlServer   NoSqlServer   `yaml:"NosqlServer"`   // 非关系型数据库的中间件
}
type IcarusServer map[string]string
type NoSqlServer map[string]string
type ControlServer map[string]string

//WindowsAdConfig 域账号登录配置
type WindowsAdConfig struct {
	Status   bool   `yaml:"Status"`
	Server   string `yaml:"Server"`
	Port     int    `yaml:"Port"`
	BaseOn   string `yaml:"BaseOn"`
	Security bool   `yaml:"Security"`
}

//RedisServer redis配置
type RedisServer struct {
	Address  string `yaml:"Address"`
	Password string `yaml:"Password"`
}

// 应用初始信息
type App struct {
	Name              string      `yaml:"Name"`
	Version           string      `yaml:"Version"`
	DB                DbConfig    `yaml:"DB"`
	Redis             RedisServer `yaml:"Redis"`
	SessionExpireTime int         `yaml:"SessionExpireTime"`
}

//func DBParse() {
//	config, err := ioutil.ReadFile("./config/config.yaml")
//
//	if err != nil {
//		log.Fatalf("read config file error: %v", err)
//	}
//
//	var res App
//	if err := yaml.Unmarshal(config, &res); err != nil {
//		log.Fatalf("unmarshal config file error: %v", err)
//	}
//	CONFIG = &res
//	ConfigDe := res
//	ConfigDe.DB.Password = "******"
//	LogDebugf("load config from file", logrus.Fields{"CONFIG": ConfigDe})
//	RedisInit()
//	LogInfo(fmt.Sprintf("load redis config success, redis address is %s", res.Redis))
//}

func ConfigParse() {
	vpConfig = viper.New()
	vpConfig.SetConfigName("config")
	vpConfig.SetConfigType("yaml")
	vpConfig.AddConfigPath("./config/")
	if err := vpConfig.ReadInConfig(); err != nil {
		log.Fatalf("read config file error: %v", err)
	}
	var res App
	if err := vpConfig.Unmarshal(&res); err != nil {
		log.Fatalf("unmarshal config file error: %v", err)
	}
	CONFIG = &res
	ConfigDe := res
	ConfigDe.DB.Password = "******"
	LogDebugf("load config from file", logrus.Fields{"CONFIG": ConfigDe})
	RedisInit()
	LogInfo(fmt.Sprintf("load redis config success, redis address is %s", res.Redis))
}
