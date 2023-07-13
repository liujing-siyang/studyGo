package channel_test

import (
	"fmt"
	"time"
)

type Stream struct {
	source <-chan any
}

func Range(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}

func (s Stream) walk() Stream {
	pipe1 := make(chan any, 2)
	ss := Range(pipe1)
	go func() {
		time.Sleep(time.Second * 3)
		for item := range s.source {
			val := item.(int)
			fmt.Println(val)
			pipe1 <- val * val
		}
		close(pipe1)
	}()

	return ss
}
func ChCopy() {
	pipe := make(chan any, 2)
	s := Range(pipe)
	s1 := s.walk()
	go func() {
		// 发送为何一定要异步
		for i := 0; i < 10; i++ {
			pipe <- i
		}
		close(pipe)
	}()

	for item := range s1.source {
		val := item.(int)
		fmt.Printf("--%d--\n", val)
	}
}
