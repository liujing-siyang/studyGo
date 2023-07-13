package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 创建文件/目录监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
        tic := time.NewTicker(10 * time.Second)
		notifyFlag := false
		const writeOrCreateMask = fsnotify.Write | fsnotify.Create
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 打印监听事件
				log.Println("event:", event)
				// ok = event.Has(fsnotify.Write)
				// if ok {
				// 	fmt.Println("write EOF")
				// }
				if event.Op&writeOrCreateMask != 0 {
					tic.Reset(10 * time.Second)
					notifyFlag = true
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			case <-tic.C:
				if notifyFlag {
					fmt.Println("write EOF")
                    notifyFlag = false
				}
			}
		}
	}()
	// 监听当前目录
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		f, err := os.Create("test.txt")
		if err != nil {
			fmt.Println("create file err")
		}
		for i := 0; i < 10000; i++ {
			_, err = f.WriteString(fmt.Sprintf("add col number %d\n", i))
			if err != nil {
				fmt.Printf("add col err %d", i)
			}
		}
		fmt.Println("write file eof")
		f.Close()
	}()
	<-done
}
