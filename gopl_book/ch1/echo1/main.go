// Echo1 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
)

//	func main() {
//		var s, sep string
//		for i := 1; i < len(os.Args); i++ {
//			s += sep + os.Args[i]
//			sep = " "
//		}
//		fmt.Println(s)
//	}
func main() {
	// s, sep := "", ""

	for _, str := range os.Args[1:] {
		fmt.Println(str)
	}

}

// os.Args[0] 是命令本身二进制可执行文件的路径被打印出来
// os.Args[1:] 是命令行参数
// func main() {
// 	// fmt.Println(os.Args[1:])
// 	fmt.Println(strings.Join(os.Args[1:], " "))
// }
