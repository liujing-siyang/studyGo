package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"sync/atomic"

	"github.com/Jeffail/tunny"

	"github.com/panjf2000/ants/v2"
)

func main() {
	// test1()
	// test2()
	test3()
}

// 测试tunny
func test1() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("numCpus：%d\n", numCPUs)
	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		// TODO: Something CPU heavy with payload
		var err error
		f, ok := payload.(func() error)
		if ok {
			err = f()
		}
		return err
	})
	defer pool.Close()
	var wg sync.WaitGroup
	now := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			res := pool.Process(work()) //Process是同步的，协程池成了伪命题？
			fmt.Println(res)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("spend time: %v", time.Since(now))
}

func test2() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("numCpus：%d\n", numCPUs)
	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		// TODO: Something CPU heavy with payload
		var err error
		f, ok := payload.(func() error)
		if ok {
			err = f()
		}
		return err
	})
	now := time.Now()
	for i := 0; i < 100; i++ {
		res := pool.Process(work())
		fmt.Println(res)
	}
	fmt.Printf("spend time: %v", time.Since(now))
	time.Since(now)
}

func work() error {
	fmt.Println("hello world!")
	time.Sleep(time.Millisecond * 100)
	return nil
}

// 测试ants
var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func test3() {
	defer ants.Release()

	runTimes := 1000

	// Use the common pool.
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}


