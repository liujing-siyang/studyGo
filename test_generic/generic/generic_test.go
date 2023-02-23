package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySliceSum(t *testing.T) {
	var s MySlice[int] = []int{1, 2, 3, 4}
	assert.Equal(t, 10, s.Sum())
}

func TestAdd(t *testing.T) {
	num1 := Add[int](1, 2)
	assert.Equal(t, 3, num1)

	num2 := Add(1.0, 2.0) //类型实参的自动推导
	assert.Equal(t, 3.0, num2)
}

func TestScoreAdd(t *testing.T) {
	var MyTwoScore Score[int]
	sum := MyTwoScore.Add(86, 94)
	assert.Equal(t, 180, sum)
}

func TestDataProcessor(t *testing.T) {
	// CSVProcessor实现了接口 DataProcessor[string] ，所以可赋值
	var processor DataProcessor[string] = CSVProcessor{}
	str := processor.Process("name,age\nbob,12\njack,30")
	assert.Equal(t, "name,age\nbob,12\njack,30", str)
	err := processor.Save("name,age\nbob,13\njack,31")
	assert.Equal(t, nil, err)
}

func TestDataProcessor2(t *testing.T) {
	// JsonProcessor实现了接口 DataProcessor2[string] ，所以可赋值
	var j JsonProcessor = JsonProcessor{}
	var processor2 ProcessorList[JsonProcessor] = []JsonProcessor{j}
	p := processor2[0]
	str := p.Process("name,age\nbob,12\njack,30")
	assert.Equal(t, "name,age\nbob,12\njack,30", str)
	err := p.Save("name,age\nbob,13\njack,31")
	assert.Equal(t, nil, err)
}
