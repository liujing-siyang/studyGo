package pipeline

import (
	"fmt"
	"testing"
)

func TestDecorator(t *testing.T) {
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	odd := Do(myodd)
	sq := Do(mysq)
	sum := Do(mysum)
	for n := range sum(sq(odd(echo(nums)))){
		fmt.Println(n)
	}
}


func TestPipeline(t *testing.T) {
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Map语义
	p := transform(echo(nums),Do(mysum)).(<-chan int)
	res := <- p
	fmt.Println(res)
	for n := range sum(sq(odd(echo(nums)))) {
		fmt.Println(n)
	}
	for n := range pipeline(nums, echo, odd, sq, sum) {
		fmt.Println(n)
	}
}


func TestFan(t *testing.T) {
	nums := makeRange(1, 10000)
	in := echo(nums)

	const nProcess = 5
	var chans [nProcess]<-chan int
	// 为什么五个计算协程分段计算出来的结果每次都是相同的
	for i := range chans {
		chans[i] = sum(prime(in))
	}

	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}
}
