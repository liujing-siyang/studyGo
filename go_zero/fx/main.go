package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/zeromicro/go-zero/core/fx"
	"github.com/zeromicro/go-zero/core/mr"
)

func main() {
	// fx_demo()
	mr_demo()
}

func fx_demo() {
	mapReduce()
	group()
	reserve()
	parallel()
}

func mapReduce() {
	result, err := fx.From(func(source chan<- interface{}) {
		for i := 0; i < 10; i++ {
			source <- i
		}
	}).Map(func(item interface{}) interface{} {
		i := item.(int)
		return i * i
	}).Filter(func(item interface{}) bool {
		i := item.(int)
		return i%2 == 0
	}).Distinct(func(item interface{}) interface{} {
		return item
	}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		var result int
		for item := range pipe {
			i := item.(int)
			result += i
		}
		return result, nil
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func group() {
	// 例子 按照首字符"g"或者"p"分组，没有则分到另一组
	ss := []string{"golang", "google", "php", "python", "java", "c++"}
	fx.From(func(source chan<- interface{}) {
		for _, s := range ss {
			source <- s
		}
	}).Group(func(item interface{}) interface{} {
		if strings.HasPrefix(item.(string), "g") {
			return "g"
		} else if strings.HasPrefix(item.(string), "p") {
			return "p"
		}
		return ""
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func reserve() {
	fx.Just(1, 2, 3, 4, 5).Reverse().ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func parallel() {
	fx.Just(1, 2, 3, 4, 5).Parallel(func(item interface{}) {
		id := item
		fmt.Printf("处理 %v 的日志! \n", id)
	})
}

func mr_demo() {
	// productDetail(1, 2)
	mrFinishVoid(3,4)
	v,err  := checkLegal([]int64{88,99,11,22,33,44,55,66,77})
	if err != nil{

	}
	fmt.Printf("%v\n",v)
}

type ProductDetail struct {
	UserId int64
	Order  int64
	Store  int64
}

func productDetail(uid, pid int64) (*ProductDetail, error) {
	var pd ProductDetail
	err := mr.Finish(func() (err error) {
		fmt.Println("userrpc")
		pd.UserId = uid
		return nil
	}, func() (err error) {
		fmt.Println("storerpc")
		pd.Order = pid
		return nil
	}, func() (err error) {
		fmt.Println("orderrpc")
		pd.Store = pid
		return nil
	})

	if err != nil {
		log.Printf("product detail error: %v", err)
		return nil, err
	}

	return &pd, nil
}

func mrFinishVoid(uid, pid int64) (*ProductDetail, error) {
	var pd ProductDetail
	mr.FinishVoid(func(){
		fmt.Println("userrpc")
		pd.UserId = uid
	}, func() {
		fmt.Println("storerpc")
		pd.Order = pid
	}, func() {
		fmt.Println("orderrpc")
		pd.Store = pid
	})
	return &pd, nil
}

func checkLegal(uids []int64) ([]int64, error) {
	r, err := mr.MapReduce(func(source chan<- any) {
		for _, uid := range uids {
			source <- uid
		}
	}, func(item any, writer mr.Writer[any], cancel func(error)) {
		uid := item.(int64)
		ok, err := check(uid)
		if err != nil {
			cancel(err)
		}
		if ok {
			writer.Write(uid)
		}
	}, func(pipe <-chan any, writer mr.Writer[any], cancel func(error)) {
		var uids []int64
		for p := range pipe {
			uids = append(uids, p.(int64))
		}
		writer.Write(uids)
	})
	if err != nil {
		log.Printf("check error: %v", err)
		return nil, err
	}

	return r.([]int64), nil
}

func check(uid int64) (bool, error) {
	// do something check user legal
	if uid == 11{
		return false,fmt.Errorf("i hate the number %d",uid)
	}
	return true, nil
}
