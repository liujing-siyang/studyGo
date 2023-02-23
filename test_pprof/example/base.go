package example

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// 用于服务的测试
func NewProfileHttpServer(addr string) {
	go func() {
		log.Fatalln(http.ListenAndServe(addr, nil))
	}()
}

func Web() {
	mux := http.NewServeMux()
	mux.HandleFunc("/fib/", fibHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	NewProfileHttpServer(":9999")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func fibHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Path[len("/fib/"):])
	if err != nil {
		responseError(w, err)
		return
	}

	var result int
	for i := 0; i < 1000; i++ {
		result = Fib(n)
	}
	response(w, result)
}

func Fib(n int) int {
	if n <= 1 {
		return 1
	}
	return Fib(n-1) + Fib(n-2)
}

func response(w http.ResponseWriter, result int) {
	w.Write([]byte(fmt.Sprintf("result:%d", result)))
}

func responseError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
}
