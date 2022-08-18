package main

import (
	"encoding/base64"
	"fmt"
	"test_aes/util"
)

func main() {
	password := "71MKCQ9X1kb"

	// data := util.AESEncryptECB([]byte(password), []byte("H6Jrh5qUte2YU2z8u8ouR6VdY3iIdBrx"))
	// cre := util.AESDecryptECB(data, []byte("H6Jrh5qUte2YU2z8u8ouR6VdY3iIdBrx"))
	// fmt.Println(string(cre))
	// fmt.Println(base64.StdEncoding.EncodeToString(data))

	data, _ := util.AESEncryptCBC([]byte(password), []byte("H6Jrh5qUte2YU2z8u8ouR6VdY3iIdBrx"))
	cre, _ := util.AESDecryptCBC(data, []byte("H6Jrh5qUte2YU2z8u8ouR6VdY3iIdBrx"))
	fmt.Println(string(cre))
	fmt.Println(base64.StdEncoding.EncodeToString(data))
}
