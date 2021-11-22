# 也许是最好的汉字数字(中文数字)-阿拉伯数字转换工具

## The Best Tool of Chinese Number to Digits 

Chinese2digits 是一个将中文数字（大写数字） 转化为阿拉伯数字的工具。这个工具是自然语言处理系统（NLP）中非常重要的组件。常见的应用场景有:聊天机器人，知识图谱，OCR（光学字符识别）系统，等等。

Chinese2digits is tool that transfer the Numbers writen in Capital Chinese to digits number. This tool is a vital component of NLP system such as: chat robot, knowledge graph, and OCR system etc.




用例：
example:
```
输入: 啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万
输出: takeNumberFromString 方法
{
	'inputText': '啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万', 
	'replacedText': '啊啦啦啦300.0.0100000.0你好我200000.030.0你好啊300.0.0咯咯咯-0.0034啦啦啦300.0.00000.0', 
	'CHNumberStringList': ['300', '十万', '20万', '.3%万', '300', '-.34%', '300万'], 
	'digitsStringList': ['300.0', '100000.0', '200000.0', '30.0', '300.0', '-0.0034', '3000000.0']
}

输入：伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二
输出：TakeChineseNumberFromString 方法
{
    'inputText': '伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二',
    'replacedText': '570000017今天天气不错3020050.0啦啦啦啦-0.0028你好啊32000',
    'CHNumberStringList': ['五亿七千万十七', '百分之三亿二百万五', '负百分之点二八', '三万二'],
    'digitsStringList': ['570000017', '3020050.0', '-0.0028', '32000']
}
```

takeNumberFromString 方法和 takeChineseNumberFromString (go 语言第一个字母大写)的主要不同在于：

takeNumberFromString会将"300三十二" 识别成2个数字"300" 和"32"; 但是 takeChineseNumberFromString 会将该数字只提取"32"。 因为纯罗马数字提取。

但是，2个方法"都会"将 300万 提取成 3000000, 因为这个只是中文数字表达的变种模式！


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

#混合提取
print(c2d.takeNumberFromString('啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万'))
print(c2d.takeChineseNumberFromString('负百分之点二八你好啊百分之三五是不是点伍零百分之负六十五点二八'))

#将百分比转为小数
print(c2d.takeDigitsNumberFromString('234%lalalal-%nidaye+2.34%',percentConvert=True))


#案例

{'inputText': '啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万', 'replacedText': '啊啦啦啦300.0.0100000.0你好我200000.030.0你好啊300.0.0咯咯咯-0.0034啦啦啦300.0.00000.0', 'CHNumberStringList': ['300', '十万', '20万', '.3%万', '300', '-.34%', '300万'], 'digitsStringList': ['300.0', '100000.0', '200000.0', '30.0', '300.0', '-0.0034', '3000000.0']}

{'inputText': '234%lalalal-%nidaye+2.34%', 'digitsNumberStringList': ['2.34', '0.0234']}

{'inputText': '百分之四百三十二万分之四三千分之五', 'replacedText': '4.320.00430.005', 'CHNumberStringList': ['百分之四百三十二', '万分之四三', '千分之五'], 'digitsStringList': ['4.32', '0.0043', '0.005']}

{'inputText': '伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二', 'replacedText': '570000017今天天气不错3020050.0啦啦啦啦-0.0028你好啊32000', 'CHNumberStringList': ['五亿七千万十七', '百分之三亿二百万五', '负百分之点二八', '三万二'], 'digitsStringList': ['570000017', '3020050.0', '-0.0028', '32000']}

{'inputText': 'llalala万三威风威风千四五', 'replacedText': 'llalala0.0003威风威风0.045', 'CHNumberStringList': ['万三', '千四五'], 'digitsStringList': ['0.0003', '0.045']}

{'inputText': '哥两好', 'replacedText': '哥两好', 'CHNumberStringList': [], 'digitsStringList': []}

{'inputText': '伍亿柒仟万拾柒百分之', 'replacedText': '570000017百分之', 'CHNumberStringList': ['五亿七千万十七'], 'digitsStringList': ['570000017']}

{'inputText': '负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八', 'replacedText': '-0.0028你好啊0.35是不是0.50-0.6528', 'CHNumberStringList': ['负百分之点二八', '百分之三五', '点五零', '百分之负六十五点二八'], 'digitsStringList': ['-0.0028', '0.35', '0.50', '-0.6528']}

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
	fmt.Println(chinese2digits.TakeNumberFromString("伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二"))
	fmt.Println("输入：伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二")
	fmt.Println("输出：TakeChineseNumberFromString 方法")
}

```

## Rust 语言代码示例
## Rust Example

### 安装
### Setup

添加如下代码到

Add this to your `Cargo.toml`:

```toml
[dependencies]
chinese2digits = "1"
```

### 使用
### Usage
``` rust

use chinese2digits::take_number_from_string;

//混合提取
let string_example = "百分之5负千分之15".to_string();
//第二个开关为 百分比转换开关，是否把百分号转换为小鼠  第三个开关为繁体转简体强制开关，是否强制转换
let test_result = take_number_from_string(&string_example, true, true);
assert_eq!(test_result.replaced_text, "0.05-0.015");
assert_eq!(test_result.ch_number_string_list, vec!["百分之5", "负千分之15"]);
assert_eq!(test_result.digits_string_list, vec!["0.05", "-0.015"]);

```


## 更新时间
## Update Time
>2021-11-22
* RUST 版本上线！！
* 重大优化:重构一些核心算法，提升计算效率

>2021-11-10
* 修复小数点右边出现零零三四的情况，可以正确转换 三百十四点零零三四  为 340.0034
* 增加单元测试用例

>2020-11-19
* 修复Bug:无法处理 1.55万这种

>2020-11-02
* 重大更新及接口更新
* 不再需要simpilify 规则 可以识别更加复杂的混合 四百三十儿分之43    30万2千这种
* 修复bug: 判断当前的汉字 与后一个提取的汉子数字 是否是连续的在原始字符串出现的，如果非连续 就出现了问题 比如  7分之前三四五
* 4分之三 4百分之一 混合识别及减少误识别
* 升级正则， 提供一个混合提取的正则表达式
* 七点十分之前  不会误识别 
* 最后的.0 去掉  精确地展示整数
* Go语言精简接口
* 增加单元测试方法 方便开发



>2020-05-31
* 升级可以识别万分之10 这种汉字数字混合
* Go语言正则升级

> 2019-12-25
* 升级可以识别300万 .3千 .50%万 数字汉字混合识别
* 将所有汉字数字与罗马数字提取
* 新增方法 takeNumberFromString，一次性地将所有罗马，中文，罗马中文混合数字全部提取

> 2019-12-9
* 更新正则引擎，采用RE2规则，不会出现长文本的递归炸弹
* go语言也有正则引擎
* 完成万三，千四五这种简称识别

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

> 2018-07-07
* 首次提交

## TODO
- [x] 两千二,三千十四  非标准表述
- [x] 繁体简体转换
- [x] 万三,千四五词汇的理解和转换
- [x] 正则引擎升级,使用顺序逻辑引擎,使用GO语言默认的正则库
- [x] 数字汉字混合模式:300万
- [x] 将所有汉字数字与罗马数字提取
- [x] 新增方法 takeNumberFromString，一次性地将所有罗马，中文，罗马中文混合数字全部提取
- [x] RUST 版本上线

> 好用记得Fork and Star。
>
> If you need the code in JavaScript,just send me message, fork and star~~



## License
chinese2digits is licensed under the Apache License 2.0.


