package main

import (
	stderr "errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/errors"
	pkgerr "github.com/pkg/errors"
)

func main() {
	// kratosUnwrapErr()
	// stddiffpkgErr()
	// TestWrap()
	ErrIs()
}

func kratosUnwrapErr() {
	err1 := pkgGetErr()
	//调用的标准库
	err2 := errors.Unwrap(err1)
	fmt.Printf("%+v\n", err1)
	fmt.Printf("%+v\n", err2)
	err3 := errors.FromError(err1)
	// err3 = err3.WithCause(fmt.Errorf("qw%s", "er")) //添加root error
	err4 := err3.Unwrap()
	fmt.Printf("%+v\n", err3)
	fmt.Printf("%s\n", err3.Error())
	fmt.Printf("%+v\n", err4)
}

// 区别是pkgerr 携带堆栈信息，%+v 获取
func stddiffpkgErr() {
	err1 := stdGetErr()
	err2 := stderr.Unwrap(err1)
	err3 := pkgGetErr()
	err4 := pkgerr.Unwrap(err3)
	fmt.Printf("%+v\n", err1)
	fmt.Printf("%+v\n", err2)
	fmt.Printf("%+v\n", err3)
	fmt.Printf("%+v\n", err4)
}

func stdGetErr() error {
	err1 := stderr.New("error1")
	err2 := fmt.Errorf("error2: [%w]", err1)
	return err2
}

func pkgGetErr() error {
	//pkgerr 携带堆栈信息
	err1 := pkgerr.New("error1")
	err2 := fmt.Errorf("error2: [%w]", err1)
	return err2
}

// pkgerr使用
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		//wrap 记录堆栈信息
		// return nil, pkgerr.Wrap(err, "open failed")
		return nil,pkgerr.Errorf("%s open failed","xxx")
	}
	defer f.Close()
	var v []byte
	_, err = f.Read(v)
	if err != nil {
		return nil, pkgerr.Wrap(err, "read failed")
	}
	return v, nil
}

func ReadConfig() ([]byte, error) {
	home := os.Getenv("HOME")
	config, err := ReadFile(filepath.Join(home, ".setting.xml"))
	return config, pkgerr.WithMessage(err, "could not read config")
}

func TestWrap() {
	_, err := ReadConfig()
	if err != nil {
		fmt.Printf("original error:[%T] [%v]\n", pkgerr.Cause(err), pkgerr.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		// fmt.Printf("stack trace:\n%+v\n", fmt.Sprintf("stack trace:\n%+v\n",err))
	}
	os.Exit(1)
}



var Myerr error = pkgerr.New("less or equal 0")

func ErrIs() {
	err := Myerr
	if err != nil {
		if pkgerr.Is(err, Myerr) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

}