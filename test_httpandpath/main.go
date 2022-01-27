package main

import (
	"fmt"
	"net/http"
	"path"
)
// FileServer返回一个处理程序，该处理程序使用根目录下的文件系统的内容来服务HTTP请求。
// 作为一种特殊情况，返回的文件服务器将任何以“/index.html”结尾的请求重定向到相同的路径，而不以“index.html”结尾。
// 要使用操作系统的文件系统实现，请使用http。Dir

// StripPrefix返回一个处理程序，通过从请求URL的路径中删除给定的前缀(如果设置了RawPath)，并调用处理程序h来服务HTTP请求。
// StripPrefix通过响应一个HTTP 404 not found错误来处理一个不是以前缀开头的路径的请求。
// 前缀必须完全匹配:如果请求中的前缀包含转义字符，那么回复也是一个HTTP 404未找到错误。
func main() {
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
	//Static("assets", "/usr/geektutu/blog/static")
}

func createStaticHandler(relativePath string, fs http.FileSystem) {
	absolutePath := path.Join("prefix", relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	fmt.Println(fileServer)
	//fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

}

// serve static files
func Static(relativePath string, root string) {
	fmt.Println(http.Dir(root))
	createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	fmt.Println(urlPattern)
}
