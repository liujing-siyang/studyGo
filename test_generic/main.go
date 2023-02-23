package main

import (
	"fmt"
	"io"
	"test_generic/generic"
)

type MySlice[T io.Reader] []T

type Myint int

func (Myint) Read(p []byte) (n int, err error) {
	return
}

func main() {
	var a generic.Slice[int] = []int{1, 2, 3}
	fmt.Printf("Type Name: %T\n", a)
	var b generic.MyMap[string, float64] = map[string]float64{
		"jack_score": 9.6,
		"bob_score":  8.4,
	}
	fmt.Printf("Type Name: %T\n", b)
	var c MySlice[Myint] = []Myint{1, 2, 3}
	fmt.Printf("Type Name: %T\n", c)
	var d generic.DataProcessor[string] = generic.CSVProcessor{}
	fmt.Printf("Type Name: %T\n", d)
	var e generic.ProcessorList[generic.JsonProcessor] = []generic.JsonProcessor{}
	fmt.Printf("Type Name: %T\n", e)
}

// 参考链接 https://segmentfault.com/a/1190000041634906#item-4-4
