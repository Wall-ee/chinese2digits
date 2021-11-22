package chinese2digits

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	// got := chinese2digits.TakeNumberFromString("百分之5负千分之15") // 程序输出的结果
	// want := []string{"a", "b", "c"}                          // 期望的结果
	// if !reflect.DeepEqual(want, got) {                       // 因为slice不能比较直接，借助反射包中的方法比较
	// 	t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	// }

	type test struct { // 定义test结构体
		input string
		// sep   string
		replacedText       string
		CHNumberStringList []string
		digitsStringList   []string
		percentConvert     bool
	}
	tests := map[string]test{ // 测试用例使用map存储
		"simple": {input: "百分之5负千分之15",
			percentConvert:     true,
			replacedText:       "0.05-0.015",
			CHNumberStringList: []string{"百分之5", "负千分之15"},
			digitsStringList:   []string{"0.05", "-0.015"}},
		"1": {input: "三零万二零千拉阿拉啦啦30万20千嚯嚯或百四嚯嚯嚯四百三十二分之2345啦啦啦啦",
			percentConvert:     false,
			replacedText:       "320000拉阿拉啦啦30000020000嚯嚯或4%嚯嚯嚯2345/432啦啦啦啦",
			CHNumberStringList: []string{"三零万二零千", "30万", "20千", "百四", "四百三十二分之2345"},
			digitsStringList:   []string{"320000", "300000", "20000", "4%", "2345/432"}},
		"2": {input: "啊啦啦啦300十万你好我20万.3%万你好啊300咯咯咯-.34%啦啦啦300万",
			percentConvert:     true,
			replacedText:       "啊啦啦啦30000000你好我20000030你好啊300咯咯咯-0.0034啦啦啦3000000",
			CHNumberStringList: []string{"300十万", "20万", ".3%万", "300", "-.34%", "300万"},
			digitsStringList:   []string{"30000000", "200000", "30", "300", "-0.0034", "3000000"}},
		"3": {input: "aaaa.3%万啦啦啦啦0.03万",
			percentConvert:     true,
			replacedText:       "aaaa30啦啦啦啦300",
			CHNumberStringList: []string{".3%万", "0.03万"},
			digitsStringList:   []string{"30", "300"}},
		"4": {input: "十分之一",
			percentConvert:     true,
			replacedText:       "0.1",
			CHNumberStringList: []string{"十分之一"},
			digitsStringList:   []string{"0.1"}},
		"5": {input: "四分之三啦啦五百分之二",
			percentConvert:     false,
			replacedText:       "3/4啦啦2/500",
			CHNumberStringList: []string{"四分之三", "五百分之二"},
			digitsStringList:   []string{"3/4", "2/500"}},
		"6": {input: "4分之3负五分之6咿呀呀 四百分之16ooo千千万万",
			percentConvert:     true,
			replacedText:       "0.75-1.2咿呀呀 4.16ooo千千万万",
			CHNumberStringList: []string{"4分之3", "负五分之6", "四百分之16"},
			digitsStringList:   []string{"0.75", "-1.2", "4.16"}},
		"7": {input: "百分之四百三十二万分之四三千分之五今天天气不错三百四十点零零三四",
			percentConvert:     true,
			replacedText:       "4.320.00430.005今天天气不错340.0034",
			CHNumberStringList: []string{"百分之四百三十二", "万分之四三", "千分之五", "三百四十点零零三四"},
			digitsStringList:   []string{"4.32", "0.0043", "0.005", "340.0034"}},
		"8": {input: "四千三",
			percentConvert:     true,
			replacedText:       "4300",
			CHNumberStringList: []string{"四千三"},
			digitsStringList:   []string{"4300"}},
		"9": {input: "伍亿柒仟万拾柒今天天气不错百分之三亿二百万五啦啦啦啦负百分之点二八你好啊三万二",
			percentConvert:     true,
			replacedText:       "570000017今天天气不错3020050啦啦啦啦-0.0028你好啊32000",
			CHNumberStringList: []string{"五亿七千万十七", "百分之三亿二百万五", "负百分之点二八", "三万二"},
			digitsStringList:   []string{"570000017", "3020050", "-0.0028", "32000"}},
		"10": {input: "llalala万三威风威风千四五",
			percentConvert:     true,
			replacedText:       "llalala0.0003威风威风0.045",
			CHNumberStringList: []string{"万三", "千四五"},
			digitsStringList:   []string{"0.0003", "0.045"}},
		"11": {input: "伍亿柒仟万拾柒百分之",
			percentConvert:     true,
			replacedText:       "570001700分之",
			CHNumberStringList: []string{"五亿七千万十七百"},
			digitsStringList:   []string{"570001700"}},
		"12": {input: "负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八",
			percentConvert: true,
			//注意这里可能有 float 的浮点问题要处理
			replacedText:       "-0.0028你好啊0.35是不是0.5-0.6528",
			CHNumberStringList: []string{"负百分之点二八", "百分之三五", "点五零", "百分之负六十五点二八"},
			digitsStringList:   []string{"-0.0028", "0.35", "0.5", "-0.6528"}},
		"13": {input: "2.55万nonono3.1千万",
			percentConvert:     true,
			replacedText:       "25500nonono31000000",
			CHNumberStringList: []string{"2.55万", "3.1千万"},
			digitsStringList:   []string{"25500", "31000000"}},
		"14": {input: "拾",
			percentConvert:     true,
			replacedText:       "拾",
			CHNumberStringList: []string{},
			digitsStringList:   []string{}},
		"15": {input: "零零零三四二啦啦啦啦12.550万啦啦啦啦啦零点零零三四二万",
			percentConvert:     true,
			replacedText:       "000342啦啦啦啦125500啦啦啦啦啦34.2",
			CHNumberStringList: []string{"零零零三四二", "12.550万", "零点零零三四二万"},
			digitsStringList:   []string{"000342", "125500", "34.2"}},

		"16": {input: "10000000000000000000000000000000000000000000连",
			percentConvert:     true,
			replacedText:       "10000000000000000000000000000000000000000000连",
			CHNumberStringList: []string{"10000000000000000000000000000000000000000000"},
			digitsStringList:   []string{"10000000000000000000000000000000000000000000"}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) { // 使用t.Run()执行子测试
			// fmt.Println(tc.input)
			got := TakeNumberFromString(tc.input, tc.percentConvert, true)

			resultReplacedText := got.(map[string]interface{})["replacedText"]
			resultCHNumberStringList := got.(map[string]interface{})["CHNumberStringList"]
			resultDigitsStringList := got.(map[string]interface{})["digitsStringList"]
			if !reflect.DeepEqual(resultReplacedText, tc.replacedText) {
				t.Errorf("excepted:%v, got:%v", tc.replacedText, resultReplacedText)
			}
			if !reflect.DeepEqual(resultCHNumberStringList, tc.CHNumberStringList) {
				t.Errorf("excepted:%v, got:%v", tc.CHNumberStringList, resultCHNumberStringList)
			}
			if !reflect.DeepEqual(resultDigitsStringList, tc.digitsStringList) {
				t.Errorf("excepted:%v, got:%v", tc.digitsStringList, resultDigitsStringList)
			}
		})
	}
}

// func main() {
// 	fmt.Println(TakeChineseNumberFromString("三十万", true, true, true))
// 	fmt.Println("这个函数被调用了")
// }

// func main() {
// 	fmt.Println("输入：负百分之点二八你好啊三五四百分之三五是不是点伍零百分之负六十五点二八")
// 	fmt.Println("输出：TakeChineseNumberFromString 方法")
// 	fmt.Println(chinese2digits.TakeNumberFromString("百分之5负千分之15"))
// 	fmt.Println(chinese2digits.TakeNumberFromString("三十万"))
// 	fmt.Println(chinese2digits.TakeNumberFromString("十万啦啦啦300万nihao400十五点八"))
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("百分之四百三十二万分之四三千分之五"))
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("llalala万三威风威风千四五", true, false))
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("负百分之点二八你好啊三五四百分之三五是不是点五零百分之负六十五点二八百分之四十啦啦啦啊四万三千四百二"))
// 	//自动判断转换方式，使用正则引擎
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("壹亿叁仟伍佰万你好亿万", "auto", true, true))
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("一千八百万啦啦啦四万七,皮皮四千万十七", true, true, true, false))
// 	fmt.Println(chinese2digits.TakeChineseNumberFromString("今天万科怎么样负点三六姹紫嫣红千千万万", nil, true, true, false))
// }
