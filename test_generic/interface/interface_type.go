package interface_type

// ~ : 指定底层类型
// ~后面的类型不能为接口
// ~后面的类型必须为基本类型

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32
}

// 接口类型 Float 代表了一个类型集合，所有以 float32 或 float64 为底层类型的类型，都在这一类型集之中
type Float interface {
	~float32 | ~float64
}

// 使用接口作为类型约束
type Slice[T Int | Uint | Float] []T // 使用 '|' 将多个接口类型组合

var s Slice[int] 

// 接口定义的变化：从方法集变为类型集
// ReaderWriter 接口看成代表了一个类型的集合，所有实现了 Read() Writer() 这两个方法的类型都在接口代表的类型集合当中
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}

type MyIo[T ReadWriter] []T

type File struct{}
func (File)Read(p []byte) (n int, err error){return}
func (File)Write(p []byte) (n int, err error){return}


type Net struct{}
func (Net)Read(p []byte) (n int, err error){return}
func (Net)Write(p []byte) (n int, err error){return}

var myFile MyIo[File] = []File{}
var myNet MyIo[Net] = []Net{}

// 接口实现定义的变化
// 类型 T 实现了接口 I 
// T 不是接口时：类型 T 是接口 I 代表的类型集中的一个成员 (T is an element of the type set of I)
// 如上,File类型是接口ReadWriter代表的类型集中的一个成员
// T 是接口时： T 接口代表的类型集是 I 代表的类型集的子集 (Type set of T is a subset of the type set of I),不太好理解
type IreadWrite interface{
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

type FileOrNet struct{}
func (FileOrNet)Read(p []byte) (n int, err error){return}
func (FileOrNet)Write(p []byte) (n int, err error){return}

func TestIread()ReadWriter{
	var file IreadWrite = FileOrNet{}//IreadWrite类型集是ReadWriter类型集的子集
	return file
}

