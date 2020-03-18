package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 应用基本配置信息
var app = newApplication()

const (
	// 开发环境
	ENVIRONMENT_DEVELOP = "dev"

	// 测试环境
	ENVIRONMENT_TEST = "test"

	// demo
	ENVIRONMENT_DEMO = "demo"

	// 线上环境
	ENVIRONMENT_RELEASE = "release"
)

// 应用基本的配置结构体
type Application struct {
	AppEnv         string `json:"app_env"`
	AppName        string `json:"app_name"`     // 应用名称
	AppVersion     string `json:"app_version"`  // 应用版本以前缀v 开头
	AppPort        int    `json:"app_port"`     // 应用监听的端口
	AppGraceReload bool   `json:"grace_reload"` // 应用是否支持优雅重启
}

func (r *Application) ToString() string {
	v, _ := json.Marshal(r)
	return string(v)
}

// 初始化应用信息
func newApplication() *Application {
	var env = GetEnv()
	if env == "" { // 默认环境为线上环境
		env = ENVIRONMENT_RELEASE
	}
	var io = NewSystemOut().SetInfoType(LogLevelInfo)
	io.SystemOutPrintf("Env is: '%s'", env)
	return &Application{
		AppName:        "app",
		AppVersion:     "v1.0",
		AppPort:        8080,
		AppEnv:         env,
		AppGraceReload: false,
	}
}

func GetEnv() string {
	return os.Getenv("APP_ENV")
}

// 获取当前应用的基本配置
func GetAppConfig() *Application {
	return app
}
func GetConfigFileDirectory() string {
	var env = ""
	if app.AppEnv != "" {
		env = app.AppEnv + "/"
	}
	return "./config/" + env
}

func PluginsApp() (err error) {
	var io = NewSystemOut().SetInfoType(LogLevelInfo)
	io.SystemOutPrintln("Load app config start")
	defer func() {
		io.SetInfoType(LogLevelInfo).SystemOutPrintf("app config is: %v \n", app.ToString())
		io.SystemOutPrintf("load app config finished \n")
	}()
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	dir := "./config/"
	io.SystemOutPrintf("config directory is : %s ", dir)
	viper.AddConfigPath(dir)   // path to look for the config file in
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		io.SetInfoType(LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) { // 热加载
		fmt.Println("Config file changed:", e.Name)
	})

	version := viper.GetString("app.version")
	app.AppVersion = "v" + version
	app.AppPort = viper.GetInt("app.port")
	app.AppName = viper.GetString("app.name")
	app.AppGraceReload = viper.GetBool("app.grace_reload")
	return
}
