package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	test2()
}

func test1() {
	now := time.Now()
	fmt.Println(now)
	// 加载时区
	_, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 按照指定时区和指定格式解析字符串时间
	//timeObj, err := time.ParseInLocation("2006/01/02 15:04:05", "2019/08/04 14:15:20", loc)
	timeObj, err := time.Parse("2006/01/02 15:04:05", "2019/08/04 14:15:20")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(timeObj)
	fmt.Println(timeObj.Sub(now))
}
func test2() {
	year, month, day := time.Now().Add(-time.Hour * 24).Date()
	monthTostr := fmt.Sprintf("%d", month)
	num, err := strconv.Atoi(monthTostr)
	if err != nil {
		fmt.Println(err)
	}
	if num < 10 {
		monthTostr = "0" + monthTostr
	}
	data := fmt.Sprintf("%d-%s-%d", year, monthTostr, day)
	datas := fmt.Sprint(year, month, day)
	fmt.Println(data)
	fmt.Println(datas)
	//使用String转字符串
	strs := time.Now().Add(-time.Hour * 24).String()
	fmt.Println(strs)
	//使用Format转字符串，可以根据自己定义的格式
	strsw := time.Now().Add(-time.Hour * 24).Format("2006-01-02")
	fmt.Println(strsw)
	// Yearday := time.Now().YearDay()
	// fmt.Println(Yearday)

}
