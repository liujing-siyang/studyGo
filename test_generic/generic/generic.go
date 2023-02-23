package generic

// 泛型类型
type Slice[T int | float32 | float64] []T

// MyMap类型定义了两个类型形参 KEY 和 VALUE。分别为两个形参指定了不同的类型约束
// 这个泛型类型的名字叫： MyMap[KEY, VALUE]
type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

// 一个泛型类型的结构体。可用 int 或 sring 类型实例化
type MyStruct[T int | string] struct {
	Name string
	Data T
}

// 一个泛型接口
type IPrintData[T int | float32 | string] interface {
	Print(data T)
}

// 一个泛型通道，可用类型实参 int 或 string 实例化
type MyChan[T int | string] chan T

//泛型receiver
type MySlice[T int | float32] []T

func (s MySlice[T]) Sum() T {
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

// 泛型函数
func Add[T int | float32 | float64](a T, b T) T {
	return a + b
}

// 不支持泛型方法，只能迂回地通过receiver使用类型形参
type Score[T int | float32 | float64] struct {
}

// 方法可以使用类型定义中的形参 T
func (receiver Score[T]) Add(a T, b T) T {
	return a + b
}

// 泛型接口
// 基本接口，只有方法，1.18之前
type DataProcessor[T any] interface {
	Process(oriData T) (newData T)
	Save(data T) error
}

// 一般接口，不光只有方法，还有类型。不能用来定义变量，只能用于泛型的类型约束中
type DataProcessor2[T any] interface {
	int | ~struct{ Data interface{} }

	Process(data T) (newData T)
	Save(data T) error
}

type CSVProcessor struct {
}

//  oriData 等的类型是 string,实现接口 DataProcessor[string]
func (c CSVProcessor) Process(oriData string) (newData string) {
	return oriData
}

func (c CSVProcessor) Save(oriData string) error {
	return nil
}

// JsonProcessor 实现了接口 DataProcessor2[string] 的两个方法，同时底层类型是 struct{ Data interface{} }。所以实现了接口 DataProcessor2[string]
type JsonProcessor struct {
	Data interface{}
}

func (c JsonProcessor) Process(oriData string) (newData string) {
	return oriData
}

func (c JsonProcessor) Save(oriData string) error {
	return nil
}

// 实例化之后的 DataProcessor2[string] 可用于泛型的类型约束
type ProcessorList[T DataProcessor2[string]] []T
