package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	h := sha256.New()
	io.WriteString(h, "His money is twice tainted: 'taint yours and 'taint mine.")
	fmt.Printf("% x\n", h.Sum(nil))

	v := sha1.New()
	io.WriteString(v, "His money is twice tainted: 'taint yours and 'taint mine.")
	fmt.Printf("% x\n", v.Sum(nil))

	q := md5.New()
	io.WriteString(q, "需要加密的密码")
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))
	salt1 := "@#$%"
	salt2 := "^&*()"

	// salt1 + 用户名 + salt2 + MD5 拼接
	io.WriteString(h, salt1)
	io.WriteString(h, "abc")
	io.WriteString(h, salt2)
	io.WriteString(h, pwmd5)

	last := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("%x\n", last)
}
