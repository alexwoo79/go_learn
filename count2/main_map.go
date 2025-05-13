package main

import (
	"fmt"
	"log"

	"github.com/alexwoo79/go_learn/datafile"
)

var path string = "datafile/votes.txt"

func main() {
	// 主函数
	// 读取文件
	lines, err := datafile.GetStrings(path)
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
