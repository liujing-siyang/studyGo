package test

// Interface 定义通过索引对元素排序的接口类型
type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

// reverse 结构体中嵌入了Interface接口
type Reverse struct {
	Interface
}

// Less 为reverse类型添加Less方法，重写原Interface接口类型的Less方法
func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

type SortBy []int

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i] < a[j] }


func Getinstance()SortBy{
	var slice SortBy
	slice = append(slice, 1,3,5,2,4,6)
	return slice
}