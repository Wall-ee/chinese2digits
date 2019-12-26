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
	fmt.Println("输入：负百分之点二八你好啊三五四百分之三五是不是点伍零百分之负六十五点二八")
	fmt.Println("输出：TakeChineseNumberFromString 方法")
	fmt.Println(chinese2digits.TakeNumberFromString("三十万"))
	fmt.Println(chinese2digits.TakeNumberFromString("十万啦啦啦300万nihao400十五点八"))
	fmt.Println(chinese2digits.TakeChineseNumberFromString("百分之四百三十二万分之四三千分之五"))
	fmt.Println(chinese2digits.TakeChineseNumberFromString("llalala万三威风威风千四五", true, false))
	fmt.Println(chinese2digits.TakeChineseNumberFromString("负百分之点二八你好啊三五四百分之三五是不是点五零百分之负六十五点二八百分之四十啦啦啦啊四万三千四百二"))
	//自动判断转换方式，使用正则引擎
	fmt.Println(chinese2digits.TakeChineseNumberFromString("壹亿叁仟伍佰万你好亿万", "auto", true, true))
	fmt.Println(chinese2digits.TakeChineseNumberFromString("一千八百万啦啦啦四万七,皮皮四千万十七", true, true, true, false))
	fmt.Println(chinese2digits.TakeChineseNumberFromString("今天万科怎么样负点三六姹紫嫣红千千万万", nil, true, true, false))
}
