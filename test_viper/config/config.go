package config

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

//当你使用如下方式读取配置时，viper会从./conf目录下查找任何以config为文件名的配置文件，
//如果同时存在./conf/config.json和./conf/config.yaml两个配置文件的话，viper会从哪个配置文件加载配置呢？ ——读取的还是json文件
//在上面两个语句下搭配使用viper.SetConfigType("yaml")指定配置文件类型可以实现预期的效果吗？——不能
func Initconfig() (err error) {
	//viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")

	//监控并重新读取配置文件
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("Config file changed:", e.Name)
	// })
	// viper.WatchConfig()

	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		panic(fmt.Sprintf("Fatal error config file:%s", err))
	}
	var stu Student
	err = viper.Unmarshal(&stu)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	fmt.Println(stu)
	// err = viper.WriteConfigAs("./conf/.configwrite")
	// if err != nil {
	// 	err = viper.WriteConfigAs("./conf/configwrite.yaml")
	// 	fmt.Printf("WriteConfigAs :%s", err)
	// }
	viper.Set("Score", 78)
	age := viper.GetInt("Age")
	fmt.Println(age)
	err = viper.Unmarshal(&stu)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	fmt.Println(stu)

	//写入配置文件
	//viper.WriteConfig()//将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径
	//viper.SafeWriteConfigAs("/conf/config_copy.json") //将当前的viper配置写入给定的文件路径。不会覆盖给定的文件(如果它存在的话)。
	return
}

func Initconfig1() (err error) {
	viper.SetConfigType("yaml") // 或者 viper.SetConfigType("YAML")

	// 任何需要将此配置添加到程序中的方法。
	var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	res := viper.Get("name") // 这里会得到 "steve"
	fmt.Printf("%v\n", res)
	return
}

//反序列化
func Initconfig2() (err error) {
	//如果你想要解析那些键本身就包含.(默认的键分隔符）的配置，你需要修改分隔符
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	v.SetDefault("chart::values", map[string]interface{}{
		"ingress": map[string]interface{}{
			"annotations": map[string]interface{}{
				"traefik.frontend.rule.type":                 "PathPrefix",
				"traefik.ingress.kubernetes.io/ssl-redirect": "true",
			},
		},
	})

	type config struct {
		Chart struct {
			Values map[string]interface{}
		}
	}

	var C config

	err = v.Unmarshal(&C)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	fmt.Println(C)
	return
}
