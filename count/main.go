package main

import (
	"fmt"
	"log"

	"github.com/alexwoo79/go_learn/datafile"
)

func main() {
	var path string = "datafile/votes.txt"
	// 主函数
	// 读取文件
	lines, err := datafile.GetStrings(path)
	if err != nil {
		log.Fatal(err)
	}
	var names []string
	var counts []int
	for _, line := range lines {
		matched := false
		for i, name := range names {
			if name == line {
				counts[i] += 1
				matched = true
				break
			}
		}
		if !matched { //!matched means mathed == false
			// If the name is not found, add it to the list
			// and set the count to 1
			names = append(names, line)
			counts = append(counts, 1)
		}
	}
	for i, name := range names {
		fmt.Printf("%s: %d\n", name, counts[i])
	}
}
