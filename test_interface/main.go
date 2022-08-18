package main

import (
	"fmt"
	"test_interface/test"
)

func main() {
	test2()
}

//接口方法重载
func test1() {
	h := test.Haier{}
	h.Dry()
}

//接口作为结构体的一个字段
func test2() {
	flag := test.Reverse.Less(test.Reverse{}, 3, 4)
	fmt.Println(flag)

}
