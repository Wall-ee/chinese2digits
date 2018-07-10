# 也许是最好的汉字数字(中文数字)-阿拉伯数字转换工具

## The Best Tool of Chinese Number to Digits 

Chinese2digits 是一个将中文数字（大写数字） 转化为阿拉伯数字的工具。这个工具是自然语言处理系统（NLP）中非常重要的组件。常见的应用场景有:聊天机器人，知识图谱，OCR（光学字符识别）系统，等等。

Chinese2digits is tool that transfer the Numbers writen in Capital Chinese to digits number. This tool is a vital component of NLP system such as: chat robot, knowledge graph, and OCR system etc.




用例：
exanple:
```
输入：负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八
输出：TakeChineseNumberFromString 方法
map[inputText:负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八 replacedText:-0.0028你好啊0.35是不是0.50-0.6528 CHNumberStringList:[负百分之点二八 百分之三五 点伍零 百分之负六十五点二八] digitsStringList:[-0.0028 0.35 0.50 -0.6528]]

```


``` go
package main

import (
	"fmt"
	"./chinese2digits"
)

func main() {
	fmt.Println("输入：负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八")
	fmt.Println("输出：TakeChineseNumberFromString 方法")
	fmt.Println(chinese2digits.TakeChineseNumberFromString("负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八", nil, true))
}

```

> 需要Python 或 JavaScript 源代码  可以给我发站内信。  好用记得Fork and Star。
>
> If you need the code in Python and JavaScript,just send me message, fork and star~~





