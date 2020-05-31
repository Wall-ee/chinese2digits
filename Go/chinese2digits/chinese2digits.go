//////////////////////////////
//Author: Li Xiaoran
//First Commit Date: 2018/07/07
//License: GPL
/////////////////////////////

package chinese2digits

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var chineseCountingString = map[string]int{"十": 10, "百": 100, "千": 1000, "万": 10000, "亿": 100000000, "拾": 10,
	"佰": 100, "仟": 1000}

// 中文转阿拉伯数字
var chineseCharNumberDict = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5,
	"六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "百": 100,
	"千": 1000, "万": 10000, "亿": 100000000, "壹": 1, "贰": 2, "叁": 3, "肆": 4, "伍": 5, "陆": 6, "柒": 7, "捌": 8, "玖": 9, "拾": 10,
	"佰": 100, "仟": 1000}
var chinesePercentString = "百分之"
var chineseSignDict = map[string]string{"负": "-", "正": "+", "-": "-", "+": "+"}
var chineseConnectingSignDict = map[string]string{".": ".", "点": ".", "·": "."}

var chinesePureNumberList = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10}

var takingChineseNumberRERules, regError1 = regexp.Compile(`(?:(?:(?:[百千万]分之[正负]{0,1})|(?:[正负](?:[百千万]分之){0,1}))` +
	`(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
	`(?:点[一二三四五六七八九幺零]+)))(?:分之){0,1}|(?:(?:[一二三四五六七八九十千万亿兆幺零百]+` +
	`(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+))(?:分之){0,1}`)

//数字汉字混合提取的正则引擎
// var takingChineseDigitsMixRERules, regError2 = regexp.Compile(`(?:(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:\%){0,1}|(?:\+|\-){0,1}\.\d+(?:\%){0,1}){0,1}` +
// 	`(?:(?:(?:(?:[百千万]分之[正负]{0,1})|(?:[正负](?:[百千万]分之){0,1}))` +
// 	`(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
// 	`(?:点[一二三四五六七八九幺零]+)))(?:分之){0,1}|` +
// 	`(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
// 	`(?:点[一二三四五六七八九幺零]+))(?:分之){0,1})`)
// mac 系统对于下面的正则会报错  还在找原因
//runtime error: invalid memory address or nil pointer dereference [signal SIGSEGV: segmentation violation]
//Unable to propogate EXC_BAD_ACCESS signal to target process and panic (see https://github.com/go-delve/delve/issues/852)
var takingChineseDigitsMixRERules, regError2 = regexp.Compile(`(?:(?:(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:[\%\‰\‱]){0,1}|` +
	`(?:\+|\-){0,1}\.\d+(?:[\%\‰\‱]){0,1})){0,1}` +
	`(?:(?:(?:[百千万]分之[正负]{0,1})|(?:[正负](?:[百千万]分之){0,1})){0,1}` +
	`(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
	`(?:点[一二三四五六七八九幺零]+))(?:分之){0,1})|` +
	`(?:(?:(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:[\%\‰\‱]){0,1}|` +
	`(?:\+|\-){0,1}\.\d+(?:[\%\‰\‱]){0,1}))` +
	`(?:(?:(?:[百千万]分之[正负]{0,1})|(?:[正负](?:[百千万]分之){0,1})){0,1}` +
	`(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
	`(?:点[一二三四五六七八九幺零]+))(?:分之){0,1}){0,1}`)

var PURE_DIGITS_RE, regError3 = regexp.Compile(`[0-9]`)

var takingDigitsRERule, regError4 = regexp.Compile(`(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:\%){0,1}|(?:\+|\-){0,1}\.\d+(?:\%){0,1}`)

func checkChineseNumberReasonable(chNumber string, opt ...bool) []string {

	digitsNumberSwitch := false
	if len(opt) == 1 {
		digitsNumberSwitch = opt[0]
	}
	//#TODO 混合提取函数
	// chineseChars := []rune(chNumber)
	//"""
	//先看数字部分是不是合理
	//"""
	//"""
	//汉字数字切割 然后再进行识别
	//"""
	chPartRegMatchResult := takingChineseNumberRERules.FindAllStringSubmatch(chNumber, -1)

	chNumberPart := ""
	if len(chPartRegMatchResult) > 0 {
		chNumberPart = chPartRegMatchResult[0][0]
	}
	digitsPartRegMatchResult := takingDigitsRERule.FindAllStringSubmatch(chNumber, -1)
	digitsNumberPart := ""
	if len(digitsPartRegMatchResult) > 0 {
		digitsNumberPart = digitsPartRegMatchResult[0][0]
	}

	digitsNumberReasonable := false
	if digitsNumberPart != "" {
		//"""
		//罗马数字合理部分检查
		//"""
		pureDigitsREResult := PURE_DIGITS_RE.FindAllStringSubmatch(digitsNumberPart, -1)
		if len(pureDigitsREResult) > 0 {
			digitsNumberReasonable = true
		} else {
			digitsNumberReasonable = false
		}

	}

	chNumberReasonable := false
	if chNumberPart != "" {
		// """
		// 如果汉字长度大于0 则判断是不是 万  千  单字这种
		// """

		for i := 0; i < len(CHINESE_PURE_NUMBER_LIST); i++ {
			if strings.Contains(chNumberPart, CHINESE_PURE_NUMBER_LIST[i]) {
				chNumberReasonable = true
				break
			}
		}

	}
	result := []string{}
	if chNumberPart != "" {
		//#中文部分合理
		if chNumberReasonable == true {
			//#罗马部分也合理 则为mix双合理模式 300三十万
			if digitsNumberReasonable == true {
				//# 看看结果需不需要纯罗马数字结果
				if digitsNumberSwitch == false {
					//#只返回中文部分
					result = []string{chNumberPart}
				} else {
					//#返回双部分
					result = []string{digitsNumberPart, chNumberPart}
				}

			} else {
				// #罗马部分不合理，中文合理  .三百万这种
				result = []string{chNumberPart}
			}

		} else {
			//#中文部分不合理，说明是单位这种
			//#看看罗马部分是否合理
			if digitsNumberReasonable == true {
				//#罗马部分合理 说明是 mix 合理模式  300万这种
				result = []string{chNumber}
			} else {
				//#罗马部分也不合理  双不合理模式  空结果
				result = []string{}
			}
		}

	} else {
		//#汉字部分啥都没有，看看罗马数字部分
		//# 看看结果需不需要纯罗马数字结果
		if digitsNumberSwitch == false {
			result = []string{}
		} else {
			//#需要纯罗马部分，检查罗马数字，罗马部分合理，返回罗马部分
			if digitsNumberReasonable == true {
				result = []string{digitsNumberPart}
			} else {
				//#罗马部分不合理  返回空
				result = []string{}
			}
		}

	}
	return result
}

//找到数组最大值
func maxValueInArray(arrayToCalc []int) int {
	//获取一个数组里最大值，并且拿到下标

	//假设第一个元素是最大值，下标为0
	maxVal := arrayToCalc[0]

	for i := 1; i < len(arrayToCalc); i++ {
		//从第二个 元素开始循环比较，如果发现有更大的，则交换
		if maxVal < arrayToCalc[i] {
			maxVal = arrayToCalc[i]
		}
	}
	return maxVal
}

// CoreCHToDigits 是核心转化函数
func CoreCHToDigits(chineseCharsToTrans string, simpilfy interface{}) string {
	chineseChars := []rune(chineseCharsToTrans)
	total := ""
	simpilfySign := false

	switch simpilfy.(type) {
	case bool:
		if simpilfy == true {
			simpilfySign = true
		} else {
			simpilfySign = false
		}
	default:
		simpilfy = nil
	}

	if simpilfy == nil {
		if len(chineseChars) > 1 {
			// 如果字符串大于1 且没有单位 ，simplilfy is true. unsimpilify is false
			for i := 0; i < len(chineseChars); i++ {
				charToGet := string(chineseChars[i])
				_, exists := chineseCountingString[charToGet]
				if !exists {
					//如果没有十百千 则采用简单拼接的方法  不采用十进制算法
					// fmt.Printf(strconv.Itoa(value2))
					simpilfySign = true
				} else {
					simpilfySign = false
					break
				}

			}
		}
	}

	if simpilfySign == false {
		tempTotal := 0
		countingUnit := 1
		countingUnitFromString := []int{1}
		// 表示单位：个十百千...
		for i := len(chineseChars) - 1; i >= 0; i = i - 1 {
			charToGet := string(chineseChars[i])
			val, _ := chineseCharNumberDict[charToGet]
			if (val >= 10) && (i == 0) {
				// 应对 十三 十四 十*之类
				if val > countingUnit {
					countingUnit = val
					tempTotal = tempTotal + val
					countingUnitFromString = append(countingUnitFromString, val)
				} else {
					countingUnitFromString = append(countingUnitFromString, val)
					// countingUnit = countingUnit * val
					countingUnit = maxValueInArray(countingUnitFromString) * val
				}
			} else if val >= 10 {
				if val > countingUnit {
					countingUnit = val
					countingUnitFromString = append(countingUnitFromString, val)
				} else {
					// countingUnit = countingUnit * val
					countingUnitFromString = append(countingUnitFromString, val)
					countingUnit = maxValueInArray(countingUnitFromString) * val
				}
			} else {
				tempTotal = tempTotal + countingUnit*val
			}
		}
		//如果 total 为0  但是 countingUnit 不为0  说明结果是 十万这种  最终直接取结果 十万
		if (tempTotal == 0) && (countingUnit) > 0 {
			// 转化为字符串
			total = strconv.Itoa(countingUnit)

		} else {
			// 转化为字符串
			total = strconv.Itoa(tempTotal)
		}
	} else {
		// total:= ""
		tempBuf := bytes.Buffer{}
		for i := 0; i < len(chineseChars); i++ {
			charToGet := string(chineseChars[i])
			val, exisits := chineseCharNumberDict[charToGet]
			if !exisits {
			} else {
				tempBuf.WriteString(strconv.Itoa(val))
			}

		}
		total = tempBuf.String()
	}
	return total
}

func convertDigitsStringToFloat(digitsString string) float64 {
	perCountingString := ""
	convertResult := digitsString
	for i := 0; i < len(CHINESE_PER_COUNTING_STRING_LIST); i++ {
		if strings.Contains(digitsString, CHINESE_PER_COUNTING_STRING_LIST[i]) {
			value, exists := CHINESE_PER_COUNTING_DICT[string(CHINESE_PER_COUNTING_STRING_LIST[i])]
			if exists {
				perCountingString = value
			} else {
				perCountingString = ""
			}
			convertResult = strings.Replace(digitsString, CHINESE_PER_COUNTING_STRING_LIST[i], "", -1)
		}
	}

	var finalTotal float64
	floatResult, err := strconv.ParseFloat(convertResult, 32)
	if err != nil {
		panic(err)
	} else {
		switch perCountingString {
		case "%":
			finalTotal = floatResult / 100
		case "‰":
			finalTotal = floatResult / 1000
		case "‱":
			finalTotal = floatResult / 10000
		default:
			finalTotal = floatResult
		}
	}
	return finalTotal
}

// ChineseToDigits 是可以识别包含百分号，正负号的函数，并控制是否将百分之10转化为0.1
func ChineseToDigits(chineseCharsToTrans string, percentConvert bool, simpilfy interface{}) string {

	//"""
	//汉字数字切割 然后再进行识别
	//"""
	chPartRegMatchResult := takingChineseNumberRERules.FindAllStringSubmatch(chineseCharsToTrans, -1)

	chPartString := ""
	if len(chPartRegMatchResult) > 0 {
		chPartString = chPartRegMatchResult[0][0]
	}
	digitsPartRegMatchResult := takingDigitsRERule.FindAllStringSubmatch(chineseCharsToTrans, -1)
	digitsPartString := ""
	if len(digitsPartRegMatchResult) > 0 {
		digitsPartString = digitsPartRegMatchResult[0][0]
	}

	finalTotal := ""

	digitsPart := 1.0
	if digitsPartString != "" {
		digitsPart = convertDigitsStringToFloat(digitsPartString)
	}

	if chPartString != "" {
		chineseCharsToTrans = standardChNumberConvert(chPartString)
		//chineseCharsToTrans = standardChNumberConvert(chineseCharsToTrans)

		chineseChars := []rune(chineseCharsToTrans)

		// """
		// 看有没有符号
		// """
		sign := ""
		for i := 0; i < len(chineseChars); i++ {
			charToGet := string(chineseChars[i])
			value, exists := chineseSignDict[charToGet]
			if exists {
				sign = value
				chineseCharsToTrans = strings.Replace(chineseCharsToTrans, charToGet, "", -1)
			}

		}

		// """
		// 看有没有百分号 千分号 万分号
		// """
		chineseChars = []rune(chineseCharsToTrans)
		perCountingString := ""

		for i := 0; i < len(CHINESE_PER_COUNTING_STRING_LIST); i++ {
			if strings.Contains(chineseCharsToTrans, CHINESE_PER_COUNTING_STRING_LIST[i]) {
				value, exists := CHINESE_PER_COUNTING_DICT[string(CHINESE_PER_COUNTING_STRING_LIST[i])]
				if exists {
					perCountingString = value
				} else {
					perCountingString = ""
				}
				chineseCharsToTrans = strings.Replace(chineseCharsToTrans, CHINESE_PER_COUNTING_STRING_LIST[i], "", -1)
			}
		}

		chineseChars = []rune(chineseCharsToTrans)

		// """
		// 小数点切割，看看是不是有小数点
		// """
		stringContainDot := false
		leftOfDotString := ""
		rightOfDotString := ""
		for key := range chineseConnectingSignDict {
			if strings.Contains(chineseCharsToTrans, key) {
				chineseCharsDotSplitList := strings.Split(chineseCharsToTrans, key)
				leftOfDotString = string(chineseCharsDotSplitList[0])
				rightOfDotString = string(chineseCharsDotSplitList[1])
				stringContainDot = true
				break
			}
		}

		convertResult := ""
		if !stringContainDot {
			convertResult = CoreCHToDigits(chineseCharsToTrans, simpilfy)
		} else {
			convertResult = ""
			tempBuf := bytes.Buffer{}

			if leftOfDotString == "" {
				// """
				// .01234 这种开头  用0 补位
				// """
				tempBuf.WriteString("0.")
				tempBuf.WriteString(CoreCHToDigits(rightOfDotString, simpilfy))
				convertResult = tempBuf.String()
			} else {
				tempBuf.WriteString(CoreCHToDigits(leftOfDotString, simpilfy))
				tempBuf.WriteString(".")
				tempBuf.WriteString(CoreCHToDigits(rightOfDotString, simpilfy))
				convertResult = tempBuf.String()
			}

		}
		//如果转换结果为空字符串 则为百分之10 这种
		if convertResult == "" {
			convertResult = "1"
		}
		convertResult = sign + convertResult

		floatResult, err := strconv.ParseFloat(convertResult, 32)
		if err != nil {
			panic(err)
		} else {
			if percentConvert == true {
				//看小数点后面有几位，如果小数点右边有数字 则手动保留一定的位数
				convertResultDotSplitList := strings.Split(convertResult, ".")
				rightOfConvertResultDotString := ""
				if len(convertResultDotSplitList) > 1 {
					rightOfConvertResultDotString = string(convertResultDotSplitList[1])
				}
				switch perCountingString {
				case "%":
					finalTotal = strconv.FormatFloat(digitsPart*floatResult/100, 'f', (len(rightOfConvertResultDotString) + 2), 32)
				case "‰":
					finalTotal = strconv.FormatFloat(digitsPart*floatResult/1000, 'f', (len(rightOfConvertResultDotString) + 3), 32)
				case "‱":
					finalTotal = strconv.FormatFloat(digitsPart*floatResult/10000, 'f', (len(rightOfConvertResultDotString) + 4), 32)
				default:
					finalTotal = strconv.FormatFloat(digitsPart*floatResult, 'f', -1, 32)
				}
				// switch perCountingString {
				// case "%":
				// 	finalTotal = strconv.FormatFloat(digitsPart*floatResult/100, 'f', -1, 32)
				// case "‰":
				// 	finalTotal = strconv.FormatFloat(digitsPart*floatResult/1000, 'f', -1, 32)
				// case "‱":
				// 	finalTotal = strconv.FormatFloat(digitsPart*floatResult/10000, 'f', -1, 32)
				// default:
				// 	finalTotal = strconv.FormatFloat(digitsPart*floatResult, 'f', -1, 32)
				// }

			} else {
				finalTotal = strconv.FormatFloat(digitsPart*floatResult, 'f', -1, 32) + perCountingString
			}
			return finalTotal
		}
	} else {
		//"""
		//如果中文部分没有数值 ，取罗马数字部分
		//"""
		if percentConvert == true {
			finalTotal = strconv.FormatFloat(digitsPart, 'f', -1, 32)
		} else {
			finalTotal = digitsPartString
		}
		return finalTotal
	}
	return finalTotal
}

type structCHAndDigit struct {
	CHNumberString    string //中文数字字符串
	digitsString      string //阿拉伯数字字符串
	CHNumberStringLen int    //中文数字字符串长度
}

//结构体排序工具重写
// 按照 structToReplace.CHNumberStringLen 从大到小排序
type structToReplace []structCHAndDigit

// 重写 Len() 方法
func (a structToReplace) Len() int { return len(a) }

// 重写 Swap() 方法
func (a structToReplace) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// 重写 Less() 方法， 从大到小排序
func (a structToReplace) Less(i, j int) bool { return a[j].CHNumberStringLen < a[i].CHNumberStringLen }

var CHINESE_PURE_COUNTING_UNIT_LIST = [5]string{"十", "百", "千", "万", "亿"}

var TRADITIONAl_CONVERT_DICT = map[string]string{"壹": "一", "贰": "二", "叁": "三", "肆": "四", "伍": "五", "陆": "六", "柒": "七",
	"捌": "八", "玖": "九"}
var SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT = map[string]string{"拾": "十", "佰": "百", "仟": "千", "萬": "万", "億": "亿"}

var SPECIAL_NUMBER_CHAR_DICT = map[string]string{"两": "二", "俩": "二"}
var CHINESE_PURE_NUMBER_LIST = []string{"幺", "一", "二", "两", "三", "四", "五", "六", "七", "八", "九", "十", "零"}

var CHINESE_PER_COUNTING_STRING_LIST = []string{"百分之", "千分之", "万分之"}
var CHINESE_SIGN_LIST = []string{"正", "负", "+", "-"}
var CHINESE_PER_COUNTING_DICT = map[string]string{"百分之": "%", "千分之": "‰", "万分之": "‱"}

func isExistItem(value interface{}, array interface{}) int {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return i
			}
		}
	}
	return -1
}

// """
// 繁体简体转换 及  单位  特殊字符转换 两千变二千
// """
func traditionalTextConvertFunc(chString string, simplifConvertSwitch bool) string {
	chStringList := []rune(chString)
	stringLength := len(chStringList)
	charToGet := ""
	if simplifConvertSwitch == true {
		for i := 0; i < stringLength; i++ {
			// #繁体中文数字转简体中文数字
			charToGet = string(chStringList[i])
			value, exists := TRADITIONAl_CONVERT_DICT[charToGet]
			if exists {
				chStringList[i] = []rune(value)[0]
			}
		}

	}

	// #检查繁体单体转换
	for i := 0; i < stringLength; i++ {
		// #如果 前后有 pure 汉字数字 则转换单位为简体
		charToGet = string(chStringList[i])
		value, exists := SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT[charToGet]
		// # 如果前后有单纯的数字 则进行单位转换
		if exists {
			switch i {
			case 0:
				if isExistItem(string(chStringList[i+1]), CHINESE_PURE_NUMBER_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			case stringLength - 1:
				if isExistItem(string(chStringList[i-1]), CHINESE_PURE_NUMBER_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			default:
				if isExistItem(string(chStringList[i-1]), CHINESE_PURE_NUMBER_LIST) != -1 ||
					isExistItem(string(chStringList[i+1]), CHINESE_PURE_NUMBER_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			}
		}
		// #特殊变换 俩变二
		charToGet = string(chStringList[i])
		value, exists = SPECIAL_NUMBER_CHAR_DICT[charToGet]
		// # 如果前后有单纯的数字 则进行单位转换
		if exists {
			switch i {
			case 0:
				if isExistItem(string(chStringList[i+1]), CHINESE_PURE_COUNTING_UNIT_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			case stringLength - 1:
				if isExistItem(string(chStringList[i-1]), CHINESE_PURE_COUNTING_UNIT_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			default:
				if isExistItem(string(chStringList[i-1]), CHINESE_PURE_COUNTING_UNIT_LIST) != -1 ||
					isExistItem(string(chStringList[i+1]), CHINESE_PURE_COUNTING_UNIT_LIST) != -1 {
					chStringList[i] = []rune(value)[0]
				}
			}
		}

	}

	return string(chStringList)
}

// """
// 标准表述转换  三千二 变成 三千零二  三千十二变成 三千零一十二
// """
func standardChNumberConvert(chNumberString string) string {
	chNumberStringList := []rune(chNumberString)
	newChNumberStringList := chNumberString

	// #大于2的长度字符串才有检测和补位的必要
	if len(chNumberStringList) > 2 {
		// #十位补一：
		tenNumberIndex := isExistItem([]rune("十")[0], chNumberStringList)
		if tenNumberIndex > -1 {
			if tenNumberIndex == 0 {
				newChNumberStringList = "一" + string(chNumberStringList)
			} else {
				// # 如果没有左边计数数字 插入1
				if isExistItem(string(chNumberStringList[(tenNumberIndex-1)]), CHINESE_PURE_NUMBER_LIST) == -1 {
					newChNumberStringList = string(chNumberStringList[:tenNumberIndex]) + "一" + string(chNumberStringList[tenNumberIndex:])
				}
			}
		}

		// #差位补零
		// #逻辑 如果最后一个单位 不是十结尾 而是百以上 则数字后面补一个比最后一个出现的单位小一级的单位
		// #从倒数第二位开始看,且必须是倒数第二位就是单位的才符合条件
		lastCountingUnit := isExistItem(string([]rune(newChNumberStringList)[len([]rune(newChNumberStringList))-2]), CHINESE_PURE_COUNTING_UNIT_LIST)
		// # 如果最末位的是百开头
		if lastCountingUnit >= 1 {
			// # 则字符串最后拼接一个比最后一个单位小一位的单位 例如四万三 变成四万三千
			// # 如果最后一位结束的是亿 则补千万
			if lastCountingUnit == 4 {
				newChNumberStringList = newChNumberStringList + "千万"
			} else {
				newChNumberStringList = newChNumberStringList + string(CHINESE_PURE_COUNTING_UNIT_LIST[lastCountingUnit-1])

			}

		}

	}
	//大于一的检查是不是万三，千四五这种
	perCountSwitch := false
	tempNewChNumberStringList := []rune(newChNumberStringList)
	if len(newChNumberStringList) > 1 {
		// #十位补一：
		fistCharCheckResult := isExistItem(string(tempNewChNumberStringList[0]), []string{"千", "万"})
		if fistCharCheckResult > -1 {
			for i := 1; i < len(tempNewChNumberStringList); i++ {
				// #其余位数都是纯数字 才能执行
				if isExistItem(string(tempNewChNumberStringList[i]), CHINESE_PURE_NUMBER_LIST) > -1 {
					perCountSwitch = true
				} else {
					perCountSwitch = false
					//有一个不合适 退出循环
					break
				}
			}
			if perCountSwitch == true {
				newChNumberStringList = string(tempNewChNumberStringList[:1]) + "分之" + string(tempNewChNumberStringList[1:])
			}
		}
	}

	return string(newChNumberStringList)
}

//检查初次提取的汉字数字是否切分正确
func checkNumberSeg(chineseNumberList []string) []string {
	newChineseNumberList := []string{}
	tempPreCounting := ""
	for i := 0; i < len(chineseNumberList); i++ {
		// #新字符串 需要加上上一个字符串 最后3位的判断结果
		newChNumberString := tempPreCounting + chineseNumberList[i]
		tempChineseNumberList := []rune(newChNumberString)
		if len(tempChineseNumberList) > 2 {
			lastString := string(tempChineseNumberList[len(tempChineseNumberList)-3:])
			// #如果最后3位是百分比 那么本字符去掉最后三位  下一个数字加上最后3位
			if isExistItem(lastString, CHINESE_PER_COUNTING_STRING_LIST) > -1 {
				tempPreCounting = lastString
				// #如果最后三位 是  那么截掉最后3位
				newChNumberString = string(tempChineseNumberList[:len(tempChineseNumberList)-3])
			} else {
				tempPreCounting = ""
			}
		}
		newChineseNumberList = append(newChineseNumberList, newChNumberString)
	}
	return newChineseNumberList
}

//检查初次提取的汉字数字是正负号是否切分正确
func checkSignSeg(chineseNumberList []string) []string {
	newChineseNumberList := []string{}
	tempSign := ""
	for i := 0; i < len(chineseNumberList); i++ {
		// #新字符串 需要加上上一个字符串 最后1位的判断结果
		newChNumberString := tempSign + chineseNumberList[i]
		tempChineseNumberList := []rune(newChNumberString)
		if len(tempChineseNumberList) > 1 {
			lastString := string(tempChineseNumberList[len(tempChineseNumberList)-1:])
			// #如果最后1位是百分比 那么本字符去掉最后三位  下一个数字加上最后1位
			if isExistItem(lastString, CHINESE_SIGN_LIST) > -1 {
				tempSign = lastString
				// #如果最后1位 是  那么截掉最后1位
				newChNumberString = string(tempChineseNumberList[:len(tempChineseNumberList)-1])
			} else {
				tempSign = ""
			}
		}
		newChineseNumberList = append(newChineseNumberList, newChNumberString)
	}
	return newChineseNumberList
}

// TakeChineseNumberFromString 将句子中的汉子数字提取的整体函数
func TakeChineseNumberFromString(chTextString string, opt ...interface{}) interface{} {

	CHNumberStringList := []string{}

	//默认参数设置
	if len(opt) > 4 {
		panic("too many arguments")
	}

	var simpilfy interface{}
	var percentConvert bool
	var traditionalConvert bool
	digitsNumberSwitch := false

	switch len(opt) {
	case 1:
		simpilfy = opt[0]
		percentConvert = true
		traditionalConvert = true
		digitsNumberSwitch = false
	case 2:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = true
		digitsNumberSwitch = false
	case 3:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		digitsNumberSwitch = false
	case 4:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		digitsNumberSwitch = opt[3].(bool)
	default:
		simpilfy = "auto"
		percentConvert = true
		traditionalConvert = true
	}

	//"""
	//简体转换开关
	//"""
	originText := chTextString

	convertedString := traditionalTextConvertFunc(chTextString, traditionalConvert)

	//正则引擎
	if regError1 != nil {
		fmt.Println(regError1)
	}
	regMatchResult := takingChineseDigitsMixRERules.FindAllStringSubmatch(convertedString, -1)

	tempText := ""
	CHNumberStringListTemp := []string{}
	for i := 0; i < len(regMatchResult); i++ {
		tempText = regMatchResult[i][0]
		CHNumberStringListTemp = append(CHNumberStringListTemp, tempText)
	}

	// #检查末尾百分之万分之问题
	CHNumberStringListTemp = checkNumberSeg(CHNumberStringListTemp)

	// 检查最后是正负号的问题
	CHNumberStringListTemp = checkSignSeg(CHNumberStringListTemp)

	//检查合理性
	for i := 0; i < len(CHNumberStringListTemp); i++ {
		// fmt.Println(aa[i])

		tempText = CHNumberStringListTemp[i]
		resonableResult := checkChineseNumberReasonable(tempText, digitsNumberSwitch)
		if len(resonableResult) > 0 {
			CHNumberStringList = append(CHNumberStringList, resonableResult...)
		}

		// CHNumberStringList = append(CHNumberStringList, regMatchResult[i][0])
	}

	//"""
	//将中文转换为数字
	//"""
	digitsStringList := []string{}
	replacedText := chTextString
	tempCHToDigitsResult := ""
	CHNumberStringLenList := []int{}
	structCHAndDigitSlice := []structCHAndDigit{}
	if len(CHNumberStringList) > 0 {
		for i := 0; i < len(CHNumberStringList); i++ {
			tempCHToDigitsResult = ChineseToDigits(CHNumberStringList[i], percentConvert, simpilfy)
			digitsStringList = append(digitsStringList, tempCHToDigitsResult)
			CHNumberStringLenList = append(CHNumberStringLenList, len(CHNumberStringList[i]))
			//将每次的新结构体附加至准备排序的
			structCHAndDigitSlice = append(structCHAndDigitSlice, structCHAndDigit{CHNumberStringList[i], digitsStringList[i], CHNumberStringLenList[i]})
		}
		//fmt.Println(structCHAndDigitSlice)

		sort.Sort(structToReplace(structCHAndDigitSlice)) // 按照 中文数字字符串长度 的逆序排序
		//"""
		//按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
		//"""
		for i := 0; i < len(CHNumberStringLenList); i++ {
			replacedText = strings.Replace(replacedText, structCHAndDigitSlice[i].CHNumberString, structCHAndDigitSlice[i].digitsString, -1)
			//fmt.Println(replacedText)
		}

	}
	finalResult := map[string]interface{}{
		"inputText":          originText,
		"replacedText":       replacedText,
		"CHNumberStringList": CHNumberStringList,
		"digitsStringList":   digitsStringList,
	}
	return finalResult

}

func TakeNumberFromString(chTextString string, opt ...interface{}) interface{} {

	//默认参数设置
	if len(opt) > 4 {
		panic("too many arguments")
	}

	var simpilfy interface{}
	var percentConvert bool
	var traditionalConvert bool
	digitsNumberSwitch := true

	switch len(opt) {
	case 1:
		simpilfy = opt[0]
		percentConvert = true
		traditionalConvert = true
		digitsNumberSwitch = true
	case 2:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = true
		digitsNumberSwitch = true
	case 3:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		digitsNumberSwitch = true
	case 4:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		digitsNumberSwitch = opt[3].(bool)
	default:
		simpilfy = "auto"
		percentConvert = true
		traditionalConvert = true
	}

	finalResult := TakeChineseNumberFromString(chTextString, simpilfy, percentConvert, traditionalConvert, digitsNumberSwitch)
	return finalResult
}

// func main() {
// 	fmt.Println(TakeChineseNumberFromString("负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八百分之四十啦啦啦啊四万三千四百二", true, true, true))
// 	fmt.Println("这个函数被调用了")
// }

// func main() {
// 	reg1, err := regexp.Compile(`(?:(?:(?:百分之[正负]{0,1})|(?:[正负](?:百分之){0,1}))(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+)))|(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+))`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(reg1.FindAllStringSubmatch("负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八百分之四十啦啦啦啊四万三千四百二", -1))

// }
