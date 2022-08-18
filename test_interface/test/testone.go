package test

import (
	"fmt"
)

// WashingMachine 洗衣机
type WashingMachine interface {
	Wash()
	Dry()
}

// 甩干器
type Dryer struct{}

// 实现WashingMachine接口的dry()方法
func (d Dryer) Dry() {
	fmt.Println("甩一甩")
}

// 海尔洗衣机
type Haier struct {
	Dryer //嵌入甩干器
}

// 实现WashingMachine接口的wash()方法
func (h Haier) Wash() {
	fmt.Println("洗刷刷")
}

func (h Haier) Dry() {
	fmt.Println("重写甩一甩")
}
