package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	filename, err := os.OpenFile("./file", os.O_RDWR|os.O_CREATE, 6)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer filename.Close()
	var timeStamp = time.Now().Unix()
	// 构造一个rand，并使用时间戳作为他的随机种子
	r := rand.New(rand.NewSource(timeStamp))
	// 取100以内的随机数
	for i := 0; i < 10; i++ {
		num := r.Intn(100)
		filename.WriteString(fmt.Sprintf("%d\n", num))
	}
	http.HandleFunc("/", hello)
	server := &http.Server{
		Addr: ":8000",
	}
	fmt.Println("server startup...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello liwenzhou.com!"))
}
