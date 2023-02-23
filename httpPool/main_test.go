package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:          2, //总共的最大空闲连接，默认为100
		MaxIdleConnsPerHost:   1, //控制单个Host的连接池大小，默认为2
		MaxConnsPerHost:       2, //控制单个Host的最大连接总数,默认为0不限制
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func get(url string) {
	resp, err := _httpCli.Get(url)
	if err != nil {
		// do nothing
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// do nothing
		return
	}
}

func TestLongShort(t *testing.T) {
	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://127.0.0.1:8087")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://127.0.0.1:8088")
		}
	}()

	time.Sleep(time.Second * 10)
}

func TestLongLong(t *testing.T) {
	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://127.0.0.1:8087")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://127.0.0.1:8089")
		}
	}()

	time.Sleep(time.Second * 10)
}

// go test -v  -run TestLong$ main_test.go ,指定测试函数
func TestLong(t *testing.T) {
	go func() {
		for i := 0; i < 1000; i++ {
			if i%100 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://127.0.0.1:8087")
		}
	}()

	time.Sleep(time.Second * 10)
}
