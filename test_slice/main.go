package main

import "fmt"

var mapSlice = make(map[string][]int)

func test1() {
	mapSlice["lihua"] = []int{89, 98, 100}
	for _, score := range mapSlice["lihua"] {
		fmt.Println(score)
	}
}

// go build -gcflags "-m -l" main.go   逃逸分析下切片分配的位置，分配在堆上
func generateSlice() (list []int) {
	list = make([]int, 3)
	for i := 0; i < 3; i++ {
		list[i] = i + 1
	}
	fmt.Printf("generateSlice:%p\n", list)
	return
}

// 切片的底层结构
// type slice struct {
// 	array unsafe.Pointer // 指针指向底层数组
// 	len   int  // 切片长度
// 	cap   int  // 底层数组容量
// }

func test2() {
	list := generateSlice()
	// 这里list只是函数返回值的副本，但都引用同一个底层数组；
	// 但底层结构却不一定相同，如果generateSlice return list[1:],那底层结构的指针就指向第二给元素，同理len和cap也会发生变化
	// 参考链接：https://www.cnblogs.com/xiaofua/archive/2022/02/23/15925926.html
	fmt.Printf("generateSlice:%p\n", list)
	for i, v := range list {
		fmt.Printf("index : %d,value : %d\n", i, v)
	}
}

func modifySlice(innerSlice []string) {
	fmt.Println("begin modify")
	innerSlice = append(innerSlice, "a") //由于容量不够，将发生扩容，深拷贝，底层数组地址改变，不再是传入的切片参数的数组
	fmt.Printf("%p, %v\n", innerSlice, &innerSlice[0])
	fmt.Println("innerSlice  len:", len(innerSlice), "cap:", cap(innerSlice))
	innerSlice[0] = "b"
	innerSlice[1] = "b"
	fmt.Println(innerSlice)
	fmt.Println("end modify")
}

func test3() {
	outerSlice := []string{"a", "a"}
	fmt.Printf("%p, %v\n", outerSlice, &outerSlice[0])
	fmt.Println("outerSlice  len:", len(outerSlice), "cap:", cap(outerSlice))
	modifySlice(outerSlice)
	fmt.Println("outerSlice  len:", len(outerSlice), "cap:", cap(outerSlice))
	fmt.Printf("%p, %v\n", outerSlice, &outerSlice[0])
	fmt.Print(outerSlice)
}

func modifySlice1(innerSlice []string) {
	innerSlice = append(innerSlice, "a")
	innerSlice[0] = "b"
	innerSlice[1] = "b"
	fmt.Println(innerSlice)
}

//初始化切片的容量为3，所以在innerSlice不会发生扩容操作，但是由于是值传递，innerSlice只是outerSlice的一个副本，
//当进行append操作的时候，也是对同一个数组进行插入，同时改变innerSlice的长度，但是outerSlice的长度（len字段）并没有发生改变，所以打印出来的还是[b b]
func test4() {
	outerSlice := make([]string, 0, 3)
	outerSlice = append(outerSlice, "a", "a")
	modifySlice1(outerSlice)
	fmt.Println(outerSlice)
}

func main() {
	test2()
}
