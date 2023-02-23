package beforetwo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	f1 := 3.1415926
	f2 := 1.41421356237

	c := NewContainer(reflect.TypeOf(f1), 16)

	if err := c.Put(f1); err != nil {
		panic(err)
	}
	if err := c.Put(f2); err != nil {
		panic(err)
	}

	g := 0.0

	if err := c.Get(&g); err != nil {
		panic(err)
	}
	fmt.Printf("%v (%T)\n", g, g) //3.1415926 (float64)
	fmt.Println(c.s.Index(0))     //1.4142135623
}

func TestType(t *testing.T) {
	intContainer := &TContainer{}
	intContainer.Put(7)
	intContainer.Put(42)
	// assert that the actual type is int
	elem, ok := intContainer.Get().(int)
	if !ok {
		fmt.Println("Unable to read an int from intContainer")
	}

	fmt.Printf("assertExample: %d (%T)\n", elem, elem)
}
