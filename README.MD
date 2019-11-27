# 也许是最好的汉字数字(中文数字)-阿拉伯数字转换工具

## The Best Tool of Chinese Number to Digits 

Chinese2digits 是一个将中文数字（大写数字） 转化为阿拉伯数字的工具。这个工具是自然语言处理系统（NLP）中非常重要的组件。常见的应用场景有:聊天机器人，知识图谱，OCR（光学字符识别）系统，等等。

Chinese2digits is tool that transfer the Numbers writen in Capital Chinese to digits number. This tool is a vital component of NLP system such as: chat robot, knowledge graph, and OCR system etc.




用例：
example:
```
输入：伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二
输出：TakeChineseNumberFromString 方法
{
    'inputText': '伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二',
    'replacedText': '570000017今天天气不错3020050.0啦啦啦啦-0.0028你好啊32000',
    'CHNumberStringList': ['五亿七千万十七', '百分之三亿二百万五', '负百分之点二八', '三万二'],
    'digitsStringList': ['570000017', '3020050.0', '-0.0028', '32000']
}
```

## Pyhon 语言代码示例
## Pyhon Example

### 安装
### setup

``` python
pip install chinese2digits
```

### 使用
### usage
``` python

import chinese2digits as c2d

#将百分比转为小数
print(c2d.takeDigitsNumberFromString('234%lalalal-%nidaye+2.34%',percentConvert=True))

#使用正则表达式
print(c2d.takeChineseNumberFromString('负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八'))

#使用普通顺序逻辑引擎
print(c2d.takeChineseNumberFromString('负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八',method='normal'))

```

## Go 语言代码示例
## GoLang Example

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


## 更新时间
## Update Time
> 2019-11-27
* 修复 四千三  五亿七千万十七 等简写单位的识别
* 增加 繁体简体转换
* 解决 两,俩识别问题

> 2019-10-07
* 解决 一千八百万的识别问题 即多个单位的组合和跨越组合

> 2019-06-02
* 提交到PIP 可以使用pip 安装

> 2019-02-22
* 用Python3 重写了算法 并增加了 正则切割引擎

> 2018-07-26
* 增加千千万万  万科 这种不合法词汇的校验

> 2018-07-26
* 增加千千万万  万科 这种不合法词汇的校验

> 2018-07-07
* 首次提交

## TODO
- [x] 两千二,三千十四  非标准表述
- [x] 繁体简体转换
- [ ] 万三,千四五词汇的理解和转换
- [ ] 正则引擎升级,使用顺序逻辑引擎,使用GO语言默认的正则库


> 好用记得Fork and Star。
>
> If you need the code in JavaScript,just send me message, fork and star~~



## License
chinese2digits is licensed under the Apache License 2.0.


