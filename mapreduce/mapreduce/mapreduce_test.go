package mapreduce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


// Map/Filter/Reduce只是一种控制逻辑，真正的业务逻辑是在传给他们的数据和那个函数来定义的。
// 这是一个很经典的业务逻辑和控制逻辑分离解耦的编程模式
// 流程就是先用Map对数据集进行规整化映射，然后通过Filter过滤掉不符合要求的部分数据，对剩下的数据使用Reduce进行化简得到最终结果
func TestMap(t *testing.T) {
	list := []string{"1", "2", "3", "4", "5", "6"}
	result := Transform(list, func(a string) string {
		return a + a + a
	})
	want := []string{"111", "222", "333", "444", "555", "666"}
	assert.Equal(t, want, result)
	intlist := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	intresult := TransformInPlace(intlist, func(a int) int {
		return a * 3
	})
	intwant := []int{3, 6, 9, 12, 15, 18, 21, 24, 27}
	assert.Equal(t, intwant, intresult)
}

func mul(a, b int) int {
	return a * b
}

func fac(n int) int {
	if n <= 1 {
		return 1
	}
	return n * fac(n-1)
}

func TestReduce(t *testing.T) {
	const size = 10
	a := make([]int, size)
	for i := range a {
		a[i] = i + 1
	}
	for i := 1; i < 10; i++ {
		// 将切片中的元素全部相乘
		out := Reduce(a[:i], mul, 1).(int)
		expect := fac(i)
		if expect != out {
			t.Fatalf("expected %d got %d", expect, out)
		}
	}
}

func TestFilter(t *testing.T) {
	const size = 10
	a := make([]int, size)
	for i := range a {
		a[i] = i + 1
	}
	out := Filter(a, func(num int) bool {
		return num%2 == 1
	})
	want := []int{1, 3, 5, 7, 9}
	assert.Equal(t, want, out)
	FilterInPlace(&a, func(num int) bool {
		return num%2 == 1
	})
	assert.Equal(t, want, a)
}

