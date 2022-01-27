package main

import "fmt"

//返回值未命名，将x赋值给返回值，defer改变x的值对返回值已经没有影响
func f1() int {
	x := 5
	defer func() {
		x++
	}()
	return x
}

//先将5赋值给返回值变量x，然后执行defer函数将返回值x加1,返回值x=6
func f2() (x int) {
	defer func() {
		x++
	}()
	return 5
}

//先将x的值赋值给返回值y，然后执行defer函数将返回值x加1，返回值y=5
func f3() (y int) {
	x := 5
	defer func() {
		x++
	}()
	return x
}

//先将5赋值给返回值变量x,defer通过传参的方式将x传递进去（值传递），内部改变x（局部变量）的值对返回值没有影响
func f4() (x int) {
	defer func(x int) {
		x++
	}(x)
	return 5
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func test1() {
	fmt.Println(f1()) //5
	fmt.Println(f2()) //6
	fmt.Println(f3()) //5
	fmt.Println(f4()) //5

	x := 1
	y := 2
	defer calc("AA", x, calc("A", x, y))
	x = 10
	defer calc("BB", x, calc("B", x, y))
	y = 20
}

func main() {
	test1()
}
