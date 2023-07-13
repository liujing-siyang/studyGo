package main

// 导入需要的包
import (
	"fmt"
	"hash"
	"hash/fnv"
	"math"
)

// 布隆过滤器结构体
type BloomFilter struct {
	bitArray []bool    // 存储布隆过滤器的位数组
	hashFunc hash.Hash // 哈希函数
	k        uint      // 哈希函数的个数
}

// 创建一个布隆过滤器
func NewBloomFilter(n uint, p float64) *BloomFilter {
	m := math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2  // 计算位数组的长度
	k := uint(math.Ceil(math.Log(2) * m / float64(n )                       // 计算哈希函数的个数
	return &BloomFilter{
		bitArray: make([]bool, int(m ,
		hashFunc: fnv.New64(),
		k:        k,
	}
}

// 添加元素
func (bf *BloomFilter) Add(data []byte) {
	for i := uint(0); i < bf.k; i++ {
		bf.hashFunc.Write(data)
		hashValue := bf.hashFunc.Sum(nil)
		bf.hashFunc.Reset()
		index := hashValue[0] % uint8(len(bf.bitArray 
		bf.bitArray[index] = true
	}
}

// 判断元素是否存在
func (bf *BloomFilter) Contains(data []byte) bool {
	for i := uint(0); i < bf.k; i++ {
		bf.hashFunc.Write(data)
		hashValue := bf.hashFunc.Sum(nil)
		bf.hashFunc.Reset()
		index := hashValue[0] % uint8(len(bf.bitArray 
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

func main() {
	// 创建一个布隆过滤器
	bf := NewBloomFilter(1000000, 0.01)

	// 添加元素
	bf.Add([]byte("hello" 
	bf.Add([]byte("world" 

	// 判断元素是否存在
	fmt.Println(bf.Contains([]byte("hello" ) // true
}


