package test

import (
	"fmt"
	"testing"
	"time"
)

func TestDone(t *testing.T) {
	messages := make(chan int, 10)
	done := make(chan bool)

	defer close(messages)
	// consumer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			select {
			case <-done:
				t.Log("child process interrupt...")
				return
			default:
				t.Logf("send message: %d\n", <-messages)
			}
		}
	}()
	// producer
	for i := 0; i < 10; i++ {
		messages <- i
	}
	time.Sleep(5 * time.Second)
	close(done) //关闭一个通道返回零值
	time.Sleep(1 * time.Second)
	t.Log("main process exit!")
}

func TestChannl(t *testing.T) {
	done := make(chan bool)
	go func() {
		flag := <-done
		fmt.Println(flag)
	}()
	time.Sleep(3 * time.Second)
	close(done)
}
