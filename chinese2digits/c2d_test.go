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
	}
	tests := map[string]test{ // 测试用例使用map存储
		"simple": {input: "百分之5负千分之15",
			replacedText:       "0.05-0.015",
			CHNumberStringList: []string{"百分之5", "负千分之15"},
			digitsStringList:   []string{"0.05", "-0.015"}},
		// "wrong sep":   {input: "a:b:c", sep: ",", want: []string{"a:b:c"}},
		// "more sep":    {input: "abcd", sep: "bc", want: []string{"a", "d"}},
		// "leading sep": {input: "沙河有沙又有河", sep: "沙", want: []string{"河有", "又有河"}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) { // 使用t.Run()执行子测试
			got := TakeNumberFromString(tc.input)

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
