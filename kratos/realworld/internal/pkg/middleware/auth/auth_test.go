package auth

import (
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	s := GenerateToken("hello world","lihua")
	fmt.Println(s)
}
