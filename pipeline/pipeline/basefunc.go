package pipeline

import (
	"reflect"
)

func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// 平方
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 过滤奇数
func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 != 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

// 求和
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}

type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
	ch := echo(nums)
	for i := range pipeFns {
		ch = pipeFns[i](ch)
	}
	return ch
}

func mysq(in <-chan int, out *chan int) {
	for n := range in {
		*out <- n * n
	}
}

func myodd(in <-chan int, out *chan int) {
	for n := range in {
		if n%2 != 0 {
			*out <- n
		}
	}
}

func mysum(in <-chan int, out *chan int) {
	var sum = 0
	for n := range in {
		sum += n
	}
	*out <- sum
}

// 编程范式中的修饰器，参考：https://coolshell.cn/articles/17929.html
func Do(f func(in <-chan int, out *chan int)) func(in <-chan int) <-chan int {
	return func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			f(in, &out)
			close(out)
		}()
		return out
	}
}

// map语义对数据进行规整
func transform(datachan, function interface{}) interface{} {

	//check the <code data-enlighter-language="raw" class="EnlighterJSRAW">slice</code> type is Slice
	chanInType := reflect.ValueOf(datachan)
	if chanInType.Kind() != reflect.Chan {
		panic("transform: not chan")
	}

	//check the function signature
	fn := reflect.ValueOf(function)
	elemType := chanInType.Type()
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("trasform: function must be of type func(" + chanInType.Type().Elem().String() + ") outputElemType")
	}

	chanOutType := reflect.ValueOf(&datachan)
	res := fn.Call([]reflect.Value{chanInType})
	if chanOutType.Elem().CanSet() {
		chanOutType.Elem().Set(res[0])
	} else {
		panic("transform: canset false")
	}
	return chanOutType.Elem().Interface()

}

// 检查函数签名
func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	//Check it is a funciton
	if fn.Kind() != reflect.Func {
		return false
	}
	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}
	// In() - returns the type of a function type's i'th input parameter.
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - returns the type of a function type's i'th output parameter.
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}
