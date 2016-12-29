package main

import (
	// "net/http"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// HttpServer
	// http.Handle("/", http.FileServer(http.Dir("./")))
	// http.ListenAndServe(":8123", nil)

	// 当前路径和上一级父路径
	var str1, str2 string
	str1 = getCurrenDirectory() // 获取当前目录
	fmt.Println(str1)
	str2 = getParentDirectory(str1) // 获取上一级目录
	fmt.Println(str2)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrenDirectory() string {
	// 当前路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// 固定路径
	// dir, err := filepath.Abs(filepath.Dir("D:/github/golangman.git/"))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
