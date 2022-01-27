package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"

	"github.com/robfig/cron/v3"
)

var cht chan bool
var chf chan bool

type NacosConfig struct {
	Addr      string `mapstructure:"addr"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataId"`
	DataGroup string `mapstructure:"dataGroup"`
	LogLevel  string `mapstructure:"logLevel"`
}

type ServerConfig struct {
	BI80MysqlUrl      string `json:"BI80MysqlUrl"`
	CrmMysqlUrl       string `json:"CrmMysqlUrl"`
	PrivilegeMysqlUrl string `json:"PrivilegeMysqlUrl"`
	PubilcMysqlUrl    string `json:"PubilcMysqlUrl"`
	SqlLog            bool   `json:"SqlLog"`

	// 微信机器人
	// webhook
	RobotWebhookPrefix string `json:"RobotWebhookPrefix"`

	RKDev string `json:"RKDev"` // 机器人报警

	RKKjBranchVisitRanking      string `json:"RKKjBranchVisitRanking"`      // 科技拜访量排行
	RKRcBranchVisitRanking      string `json:"RKRcBranchVisitRanking"`      // 人才拜访量排行
	RKKjBranchNewDemandsRanking string `json:"RKKjBranchNewDemandsRanking"` // 新增需求排行榜
	RKPersonalDeliveryRanking   string `json:"RKPersonalDeliveryRanking"`   // 交付线个人排行榜
	RKAreaDeliveryRanking       string `json:"RKAreaDeliveryRanking"`       // 交付线大区排行榜
	RKZqxDaily                  string `json:"RKZqxDaily"`                  // 战区线日报

	ZqxDailyUrl string `json:"ZqxDailyUrl"` // 战区线日报跳转地址

	CronKjZqxNotifyAt22 string `json:"CronKjZqxNotifyAt22"` // 战区线机器人通知 22点
	CronRcZqxNotifyAt22 string `json:"CronRcZqxNotifyAt22"` // 战区线日报跳转地址

	IsPushKj bool `json:"IsPushKj"` //科技推送服务开启标志
	IsPushRc bool `json:"IsPushRc"` //人才推送服务开启标志

	CmdWkhtmltoimage string `json:"CmdWkhtmltoimage"` //截图工具
}

var (
	nacosConfig    NacosConfig
	WxserverConfig ServerConfig
	group          sync.WaitGroup
)

func chushi() {
	cht = make(chan bool, 1)
	chf = make(chan bool, 1)
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
	fmt.Println(WxserverConfig)
}
func selecttest() {
	for {
		select {
		case t := <-cht:
			if t {
				fmt.Println("true")
			}
		case f := <-chf:
			if !f {
				fmt.Println("false")
			}
		}
	}
	group.Done()
}

func TestCron() {
	c := cron.New()
	i := 1
	c.AddFunc("*/1 * * * *", func() {
		if WxserverConfig.IsPushKj {
			fmt.Println("每一分钟执行一次", i)
		}
		i++
	})
	c.Start()
}
func main() {
	chushi()
	TestCron()
	group.Add(1)
	go selecttest()
	group.Wait()

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

	err = json.Unmarshal([]byte(content), &WxserverConfig)
	if err != nil {
		return err
	}
	if WxserverConfig.IsPushKj {
		cht <- WxserverConfig.IsPushKj
	} else {
		chf <- WxserverConfig.IsPushKj
	}
	//Listen config change,key=dataId+group+namespaceId.监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.DataGroup,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			content, err = client.GetConfig(vo.ConfigParam{
				DataId: nacosConfig.DataId,
				Group:  nacosConfig.DataGroup,
			})
			if err != nil {
				fmt.Println("重新获取配置失败")
				return
			}
			err = json.Unmarshal([]byte(content), &WxserverConfig)
			if WxserverConfig.IsPushKj {
				cht <- WxserverConfig.IsPushKj
			} else {
				chf <- WxserverConfig.IsPushKj
			}
			if err != nil {
				return
			}
		},
	})
	return nil
}
