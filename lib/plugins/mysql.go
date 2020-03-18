package plugins

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/juetun/study/app-dashboard/lib/common"
	"github.com/spf13/viper"
)

func PluginMysql() (err error) {
	loadMysqlConfig()
	return
}
func loadMysqlConfig() (err error) {
	var io = common.NewSystemOut().SetInfoType(common.LogLevelInfo)
	io.SystemOutPrintln("Database load  start")
	viper.SetConfigName("database") // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	dir := common.GetConfigFileDirectory()
	io.SystemOutPrintf("Database config directory is : '%s' ", dir)

	viper.AddConfigPath(dir)   // path to look for the config file in
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		io.SetInfoType(common.LogLevelError).SystemOutPrintf(fmt.Sprintf("Fatal error database file: %v \n", err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) { //热加载
		fmt.Println("Database config file changed:", e.Name)
	})
	io.SetInfoType(common.LogLevelInfo).SystemOutPrintf(fmt.Sprintf("Database load config finished \n"))
	return
}
