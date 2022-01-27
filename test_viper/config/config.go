package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Student struct {
	Id    int
	Name  string
	Age   int
	Score int
}

func Initconfig() (err error) {
	//viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
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
	return
}

func Initconfig1() (err error) {
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
