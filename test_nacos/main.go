package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

type NacosConfig struct {
	Addr      string `mapstructure:"addr"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataId"`
	DataGroup string `mapstructure:"dataGroup"`
	LogLevel  string `mapstructure:"logLevel"`
}

type ServerConfig struct {
	Company string `json:"Company"`
	Name    string `json:"Name"`
	Score   int64  `json:"Score"`
}

var (
	nacosConfig  NacosConfig
	serverConfig ServerConfig
)

func main() {
	//初始化Nacos
	err := InitNacos()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(nacosConfig)
	err = InitConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(serverConfig)
}

func InitNacos() error {
	viper.SetConfigFile("./nacos-test.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&nacosConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig() error {
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			nacosConfig.Addr,
			nacosConfig.Port,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}

	//Create client for ACM
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(nacosConfig.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("tmp/config/log"),
		constant.WithCacheDir("tmp/config/cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel(nacosConfig.LogLevel),
	)

	//创建动态配置客户端的另一种方式
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		log.Fatalln(err.Error())
	}
	//动态配置
	//get config获取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.DataGroup,
	})
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		return err
	}

	//Listen config change,key=dataId+group+namespaceId.监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.DataGroup,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			err = json.Unmarshal([]byte(content), &serverConfig)
			if err != nil {
				return
			}
		},
	})

	return nil
}
