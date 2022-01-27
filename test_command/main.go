package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {

	//testCommand()
	testCommandContext()

}

func testCommandContext() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond) //则提供的上下文用于终止进程
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, "./sleep5s.exe") //如果执行外部命令耗时超过100*time.Millisecond将调用上下文终止外部命令启动的进程
	fmt.Println(cmd.String())
	var out bytes.Buffer
	cmd.Stdout = &out //绑定输出

	cmd.Start()
	cmd.Wait()
	fmt.Println(string(out.Bytes()))
	time.Sleep(3 * time.Second)
	fmt.Println("after 3s")

	//执行命令并返回标准输出的切片
	// b, _ := cmd.Output()
	// fmt.Println(string(b))

	//执行命令并返回标准输出和错误输出合并的切片
	// output, _ := cmd.CombinedOutput()
	// fmt.Println(string(output))
}

func testCommand() {
	//str, err := Cmd("./dir/hello.exe", nil)

	//str, err := CmdAndChangeDir("dir", "./hello.exe", nil)

	err := CmdAndChangeDirToShow("dir", "./hello.exe", nil)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(str)
}

//参考链接：https://blog.csdn.net/whatday/article/details/109277998

//直接在当前目录使用并返回结果
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

//在命令位置使用并返回结果
func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDir", dir, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

//在命令位置使用并实时输出每行结果
func CmdAndChangeDirToShow(dir string, commandName string, params []string) error {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDirToFile", dir, cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return err
	}
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	err = cmd.Wait()
	return err
}
