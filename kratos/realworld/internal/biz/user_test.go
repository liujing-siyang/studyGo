package biz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPwd(t *testing.T) {
	s := hashPassword("abc")
	fmt.Println(s)
}

func TestVerifyPwd(t *testing.T) {
	a := assert.New(t)
	a.True(verifyPassword("$2a$10$pZqvYz8bZm88rD6FlqHUG.9bmOvBgItdJ9tNzR3esiFcK8NFm3bHq", "abc"))
}
