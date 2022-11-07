package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestDoRequest(t *testing.T) {
	var reqData []byte
	queryReq, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9908/v1/li/count", bytes.NewReader(reqData))
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	client := &http.Client{}
	for i := 0; i < 10; i++ {
		res, err := client.Do(queryReq)
		if err != nil {
			log.Fatalf("request err %s", err.Error())
			return
		}
		defer res.Body.Close()
		// 读取响应
		respData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		fmt.Println(string(respData))
	}
}
