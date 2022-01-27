package main

import "fmt"

var mapSlice = make(map[string][]int)

func test1() {
	mapSlice["lihua"] = []int{89, 98, 100}
	for _, score := range mapSlice["lihua"] {
		fmt.Println(score)
	}
}
func main() {
	test1()
}
