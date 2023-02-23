package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type PushedReply struct {
	Flag     bool              `json:"flag"`
	ErrorMsg string            `json:"errorMsg"`
	Data     map[string]string `json:"data"`
}

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		str, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println("read err")
		}
		fmt.Println(string(str))
		var res PushedReply
		res.Flag = false
		res.ErrorMsg = string(str)
		w.Header().Set("content-type", "text/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":12346", nil))
}
