package config

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

// **嵌入文件只能在写embed指令的Go文件的同级目录或者子目录中
//
//go:embed *.yaml
var configs embed.FS

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // 默认使用开发环境
	}
	vp := viper.New()
	// 根据环境变量 ENV 决定要读取的应用启动配置
	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(err)
	}
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configFileStream))
	if err != nil {
		// 加载不到应用配置, 阻挡应用的继续启动
		panic(err)
	}

	// 初始化 App 变量
	App = &appConfig{}
	err = vp.UnmarshalKey("app", App)
	if err != nil {
		panic(fmt.Sprintf("解析 app 配置失败: %v", err))
	}

	// 初始化 Database 变量
	Database = &databaseConfig{}
	err = vp.UnmarshalKey("database", Database)
	if err != nil {
		panic(fmt.Sprintf("解析 database 配置失败: %v", err))
	}
	Database.MaxLifeTime *= time.Second

	// 打印配置信息用于调试
	fmt.Printf("当前环境: %s\n", env)
	fmt.Printf("App 配置: %+v\n", App)
	fmt.Printf("Database 配置: %+v\n", Database)
}
