package sharedcalls_demo

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/syncx"
)

func shared_calls() {
	const round = 5
	var wg sync.WaitGroup
	barrier := syncx.NewSingleFlight()

	wg.Add(round)
	for i := 0; i < round; i++ {
		// 多个线程同时执行
		go func() {
			defer wg.Done()
			// 可以看到，多个线程在同一个key上去请求资源，获取资源的实际函数只会被调用一次
			val, err := barrier.Do("once", func() (interface{}, error) {
				// sleep 1秒，为了让多个线程同时取once这个key上的数据
				time.Sleep(time.Second)
				// 生成了一个随机的id
				return stringx.RandId(), nil
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		}()
	}

	wg.Wait()
}
