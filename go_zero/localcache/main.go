package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
)

func main() {
	// test_localcache()
	test_timingWheel()
}


func test_localcache(){
	localCache, err := collection.NewCache(time.Minute)
	if err != nil {
		fmt.Println("create local cache fail")
	}
	localCache.SetWithExpire("k1", "v1", 10*time.Second)
	v, ok := localCache.Get("k1")
	if ok {
		fmt.Println(v)
	}
	time.Sleep(10 *time.Second)
	v, ok = localCache.Get("k1")
	if ok {
		fmt.Println(v)
	}else{
		fmt.Println("no exite or expire")
	}
}

func test_timingWheel(){
	timingWheel, err := collection.NewTimingWheel(time.Second, 30, func(k, v any) {
		key, ok := k.(string)
		if !ok {
			return
		}

		fmt.Printf("key is %s,v is %s\n",key,v)
	})
	if err != nil {
		fmt.Println("create timingWheel fail")
	}
	timingWheel.SetTimer("k1","v1",15 *time.Second)
	time.Sleep(10 *time.Second)
	timingWheel.SetTimer("k1","v2",25 *time.Second)
	select{}
}