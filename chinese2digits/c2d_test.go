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
			got := chinese2digits.TakeNumberFromString(tc.input, tc.sep)
			if !reflect.DeepEqual(got["replacedText"], tc.replacedText) {
				t.Errorf("excepted:%#v, got:%#v", tc.replacedText, got["replacedText"])
			}
			if !reflect.DeepEqual(got["CHNumberStringList"], tc.CHNumberStringList) {
				t.Errorf("excepted:%#v, got:%#v", tc.CHNumberStringList, got["CHNumberStringList"])
			}
			if !reflect.DeepEqual(got["digitsStringList"], tc.digitsStringList) {
				t.Errorf("excepted:%#v, got:%#v", tc.digitsStringList, got["digitsStringList"])
			}
		})
	}
}
