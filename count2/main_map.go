package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/alexwoo79/go_learn/datafile"
)

var (
	// basePath is the base path of the project
	basePath string = "/Users/alex/go/src/github.com/alexwoo79/go_learn"
	// relativePath is the relative path of the file
	relativePath string = "datafile/votes.txt"
	// filePath is the absolute path of the file
	filePath string = filepath.Join(basePath, relativePath)
	// lines is the slice of strings read from the file
)

func main() {
	// 主函数
	// 读取文件
	lines, err := datafile.GetStrings(filePath)
	if err != nil {
		log.Fatal(err)
	}
	counts := make(map[string]int)
	for _, line := range lines {
		counts[line] += 1
	}
	// using for range to iterate over the map
	// and print the name and count
	// 迭代map
	for name, count := range counts {
		fmt.Printf("Vote for %s: %d\n", name, count)
	}

}
