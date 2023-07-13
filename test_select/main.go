package main

import (
    "fmt"
    "time"
)

func sum(a chan int, b chan int) {
	for {
		// 一次只会选择一个case执行，该case执行完才会进入下一轮循环
		select {
		case <-a:
			time.Sleep(5 * time.Second)
			fmt.Print("aaa")
		case <-b:
			time.Sleep(15 * time.Second)
			fmt.Print("bbb")
		}
	}
}

func main() {
	a := make(chan int)
	b := make(chan int)
	go sum(a, b)
	b <- 15
	time.Sleep(5 * time.Second)
	a <- 10
	time.Sleep(15 * time.Second)
}