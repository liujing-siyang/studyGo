package main

import "fmt"

func main() {
	test1()
}

func test1() {
	i := 5
	j := 10
	c := &i
	c = &j
	fmt.Println(*c)
}
