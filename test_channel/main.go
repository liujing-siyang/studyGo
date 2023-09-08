package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

/*假设有一个超长的切片，切片的元素类型为int，切片中的元素为乱序排列。限时5秒，使用多个goroutine查找切片中是否存在给定值，
在找到目标值或者超时后立刻结束所有goroutine的执行。
比如切片为：[23, 32, 78, 43, 76, 65, 345, 762, ...... 915, 86]，查找的目标值为345，
如果切片中存在目标值程序输出:"Found it!"并且立即取消仍在执行查找任务的goroutine。
如果在超时时间未找到目标值程序输出:"Timeout! Not Found"，同时立即取消仍在执行查找任务的
*/
var slice = []int{23, 32, 78, 43, 76, 65, 345, 762, 915, 86}

func test1() {
	for i := 0; i < 100000; i++ {
		slice = append(slice, i*2)
	}
	timer := time.NewTimer(time.Second * 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan bool)
	dataLen := len(slice)
	size := 10
	target := 188888
	for i := 0; i < dataLen; i += size {
		end := i + size
		if end > dataLen {
			end = dataLen - 1
		}
		go SearchTarget(ctx, slice[i:end], target, resultChan)
	}
	
	select {
	case <-timer.C:
		fmt.Fprintln(os.Stderr, "Timeout! Not Found")
		cancel()
	case <-resultChan:
		fmt.Fprintf(os.Stdout, "Found it!\n")
		cancel()
	}
	time.Sleep(time.Second * 1)
}

func SearchTarget(ctx context.Context, data []int, target int, resultChan chan bool) {
	for _, v := range data {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stdout, "Task cancelded! \n")//该函数可能在cancel发出信号前结束
			return
		default:
		}
		//time.Sleep(time.Millisecond * 1500)
		if v == target {
			resultChan <- true
			return
		}
	}
}

var wg = sync.WaitGroup{}

func main(){
	test2()
}

// 死锁
func test2(){
	ch := make(chan int,0)
	wg.Add(1)
	go func1(ch)
	num := <- ch
	fmt.Println(num)
	wg.Wait()
	fmt.Println("end")
}


func func1(ch chan int){
	defer func(){
		wg.Done()
	}()
	for i := 0;i< 2;i++{
		ch <- i
		if i == 1{
			return
		}
	}
}