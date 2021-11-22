# 也许是最好的汉字数字(中文数字)-阿拉伯数字转换工具

## The Best Tool of Chinese Number to Digits 

Chinese2digits 是一个将中文数字（大写数字） 转化为阿拉伯数字的工具。这个工具是自然语言处理系统（NLP）中非常重要的组件。常见的应用场景有:聊天机器人，知识图谱，OCR（光学字符识别）系统，等等。

Chinese2digits is tool that transfer the Numbers writen in Capital Chinese to digits number. This tool is a vital component of NLP system such as: chat robot, knowledge graph, and OCR system etc.



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

> 好用记得Fork and Star。
>
> If you need the code in JavaScript,just send me message, fork and star~~



## License
chinese2digits is licensed under the Apache License 2.0.


