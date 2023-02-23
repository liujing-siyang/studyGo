package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

// TODO 参考文章链接：https://xujiahua.github.io/posts/20200723-golang-http-reuse/

// TODO 计算客户端host(IP+Port)的数量
var m = make(map[string]int)

var ch = make(chan string, 10)

// TODO 计算链接数量
func count() {
	for s := range ch {
		m[s]++
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	logrus.Info(r.RemoteAddr) // TODO 最后打印的是 remoteAddr
	ch <- r.RemoteAddr
	// time.Sleep(time.Second)
	w.Write([]byte("helloworld"))
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func graceClose() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(ch)
		time.Sleep(time.Second)
		spew.Dump(m)
		os.Exit(0)
	}()
}

func main() {
	graceClose()
	go count()
	port := flag.Int("port", 8087, "")
	flag.Parse()

	logrus.Println("Listen port:", *port)

	http.HandleFunc("/", home)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
