package test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var m atomic.Value
var mu sync.Mutex

func ValueCopy() {
	var wg sync.WaitGroup
	ctx := context.Background()
	// m.Store(make(Map))
	wg.Add(1)
	go f1(ctx, &wg)
	wg.Wait()
}

func f1(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// m.Load()
	m := make(map[string]string)
	m["k1"] = "v1"
	a := context.WithValue(ctx, "key1", m)
	wg.Add(1)
	go f2(a, wg)
	//模拟业务时间
	time.Sleep(100 * time.Millisecond)
	s := a.Value("key1").(map[string]string)
	fmt.Printf("f1:%v\n", s)
}

func f2(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	m1 := ctx.Value("key1").(map[string]string)
	m2 := make(map[string]string)       
	for k, v := range m1 {
		m2[k] = v  
	}
	m2["k1"] = "v2"
	b := context.WithValue(ctx, "key1", m2)
	wg.Add(1)
	go f3(b, wg)
	time.Sleep(100 * time.Millisecond)
	s := b.Value("key1").(map[string]string)
	fmt.Printf("f2:%v\n", s)
}

func f3(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// m1 := ctx.Value("key1").(map[string]string)
	// m2 := make(map[string]string)      
	// for k, v := range m1 {
	// 	m2[k] = v  
	// }
	// m2["k1"] = "v3"
	// c := context.WithValue(ctx, "key1", m2)
	time.Sleep(100 * time.Millisecond)
	s := ctx.Value("key1").(map[string]string)
	fmt.Printf("f3:%v\n", s)
}

type Map map[string]string

func (Map) Read(key string) (val string) {
	m1 := m.Load().(Map)
	return m1[key]
}

func (Map) Insert(key, val string) {
	mu.Lock() // synchronize with other potential writers
	defer mu.Unlock()
	m1 := m.Load().(Map) // load current value of the data structure
	m2 := make(Map)      // create a new value
	for k, v := range m1 {
		m2[k] = v // copy all data from the current object to the new one
	}
	m2[key] = val // do the update that we need
	m.Store(m2)   // atomically replace the current object with the new one
	// At this point all new readers start working with the new version.
	// The old version will be garbage collected once the existing readers
	// (if any) are done with it.
}
