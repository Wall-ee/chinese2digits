//////////////////////////////
//Author: Li Xiaoran
//First Commit Date: 2018/07/07
//License: GPL
/////////////////////////////
package main

import (
	"fmt"

	"./chinese2digits"
)

func main() {
	fmt.Println("输入：负百分之点二八你好啊百分之三五是不是点伍零")
	fmt.Println("输出：TakeChineseNumberFromString 方法")
	fmt.Println(chinese2digits.TakeChineseNumberFromString("负百分之点二八你好啊百分之三五是不是点伍零", nil, true))
}
