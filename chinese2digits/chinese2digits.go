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

func checkChineseNumberReasonable(chNumber string) bool {
	chineseChars := []rune(chNumber)
	if len(chNumber) > 0 {
		//如果汉字长度大于0 则判断是不是 万  千  单字这种
		for i := 0; i < len(chineseChars); i++ {
			charToGet := string(chineseChars[i])
			_, exists := chinesePureNumberList[charToGet]
			if exists {
				return true
			}
		}
	}
	return false
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
		// 转化为字符串
		total = strconv.Itoa(tempTotal)
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

// ChineseToDigits 是可以识别包含百分号，正负号的函数，并控制是否将百分之10转化为0.1
func ChineseToDigits(chineseCharsToTrans string, percentConvert bool, simpilfy interface{}) string {

	chineseCharsToTrans = standardChNumberConvert(chineseCharsToTrans)

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
				perCountingString = "%"
			}
			chineseCharsToTrans = strings.Replace(chineseCharsToTrans, perCountingString, "", -1)
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

	convertResult = sign + convertResult

	finalTotal := ""
	if percentConvert == true {
		if perCountingString != "" {
			floatResult, err := strconv.ParseFloat(convertResult, 32)
			if err != nil {
				panic(err)
			} else {
				//看小数点后面有几位
				convertResultDotSplitList := strings.Split(convertResult, ".")
				if len(convertResultDotSplitList) > 1 {
					rightOfConvertResultDotString := string(convertResultDotSplitList[1])
					finalTotal = strconv.FormatFloat(floatResult/100, 'f', (len(rightOfConvertResultDotString) + 2), 32)
				} else {
					finalTotal = strconv.FormatFloat(floatResult/100, 'f', 2, 32)
				}
			}
		} else {
			finalTotal = convertResult
		}
		return finalTotal

	}
	if perCountingString != "" {
		finalTotal = convertResult + perCountingString
		return finalTotal
	}
	finalTotal = convertResult
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
			charToGet := string(chStringList[i])
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
				if isExistItem(chNumberStringList[(tenNumberIndex-1)], CHINESE_PURE_NUMBER_LIST) == -1 {
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
	if len(newChNumberStringList) > 1 {
		// #十位补一：
		fistCharCheckResult := isExistItem(string(newChNumberStringList[0]), []string{"千", "万"})
		if fistCharCheckResult > -1 {
			for i := 1; i < len(newChNumberStringList); i++ {
				// #其余位数都是纯数字 才能执行
				if isExistItem(newChNumberStringList[i], CHINESE_PURE_NUMBER_LIST) > -1 {
					perCountSwitch = true
				} else {
					perCountSwitch = false
				}
			}
			if perCountSwitch == true {
				newChNumberStringList = string([]rune(newChNumberStringList)[:1]) + "分之" + string([]rune(newChNumberStringList)[1:])
			}
		}
	}

	return string(newChNumberStringList)
}

var regRule, regError = regexp.Compile(`(?:(?:(?:百分之[正负]{0,1})|(?:[正负](?:百分之){0,1}))(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+)))|(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+))`)

// TakeChineseNumberFromString 将句子中的汉子数字提取的整体函数
func TakeChineseNumberFromString(chTextString string, opt ...interface{}) interface{} {

	tempCHNumberChar := ""
	tempCHSignChar := ""
	tempCHConnectChar := ""
	tempCHPercentChar := ""
	CHNumberStringList := []string{}
	tempTotalChar := ""

	//默认参数设置
	if len(opt) > 4 {
		panic("too many arguments")
	}

	var simpilfy interface{}
	var percentConvert bool
	var traditionalConvert bool
	var method interface{}

	switch len(opt) {
	case 1:
		simpilfy = opt[0]
		percentConvert = true
		traditionalConvert = true
		method = "regex"
	case 2:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = true
		method = "regex"
	case 3:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		method = "regex"
	case 4:
		simpilfy = opt[0]
		percentConvert = opt[1].(bool)
		traditionalConvert = opt[2].(bool)
		method = opt[3]
	default:
		simpilfy = "auto"
		percentConvert = true
		traditionalConvert = true
		method = "regex"
	}

	//"""
	//简体转换开关
	//"""
	originText := chTextString

	convertedString := traditionalTextConvertFunc(chTextString, traditionalConvert)

	if method == "regex" {
		//正则引擎
		if regError != nil {
			fmt.Println(regError)
		}
		regMatchResult := regRule.FindAllStringSubmatch(convertedString, -1)
		for i := 0; i < len(regMatchResult); i++ {
			// fmt.Println(aa[i])
			CHNumberStringList = append(CHNumberStringList, regMatchResult[i][0])
		}

	} else {
		//普通引擎
		//"""
		//将字符串中所有中文数字列出来
		//"""
		chText := []rune(convertedString)
		for i := 0; i < len(chText); i++ {
			//"""
			//看是不是符号。如果是，就记录。
			//"""
			charToGet := string(chText[i])
			_, exists := chineseSignDict[charToGet]
			if exists {
				//"""
				//如果 符号前面有数字  则 存到结果里面
				//"""
				if tempCHNumberChar != "" {
					if checkChineseNumberReasonable(tempTotalChar) {
						CHNumberStringList = append(CHNumberStringList, tempTotalChar)
						tempCHPercentChar = ""
						tempCHConnectChar = ""
						tempCHSignChar = ""
						tempCHNumberChar = ""
						tempTotalChar = ""
					} else {
						tempCHPercentChar = ""
						tempCHConnectChar = ""
						tempCHSignChar = ""
						tempCHNumberChar = ""
						tempTotalChar = ""
					}

				}
				//"""
				//如果 前一个符号赋值前，临时符号不为空，则把之前totalchar里面的符号替换为空字符串
				//"""
				if tempCHSignChar != "" {
					tempTotalChar = strings.Replace(tempTotalChar, tempCHSignChar, "", -1)
				}

				tempCHSignChar = string(chText[i])
				tempTotalChar = tempTotalChar + tempCHSignChar
				continue

			}
			//"""
			//不是字符是不是"百分之"。
			//"""
			if (len(chText) - i) >= 3 {
				if isExistItem(string(chText[i:(i+3)]), CHINESE_PER_COUNTING_STRING_LIST) > -1 {
					//"""
					//如果 百分之前面有数字  则 存到结果里面
					//"""
					if tempCHNumberChar != "" {
						if checkChineseNumberReasonable(tempTotalChar) {
							CHNumberStringList = append(CHNumberStringList, tempTotalChar)
							tempCHPercentChar = ""
							tempCHConnectChar = ""
							tempCHSignChar = ""
							tempCHNumberChar = ""
							tempTotalChar = ""
						} else {
							tempCHPercentChar = ""
							tempCHConnectChar = ""
							tempCHSignChar = ""
							tempCHNumberChar = ""
							tempTotalChar = ""
						}
					}
					//"""
					//如果 前一个符号赋值前，临时符号不为空，则把之前totalchar里面的符号替换为空字符串
					//"""
					if tempCHPercentChar != "" {
						tempTotalChar = strings.Replace(tempTotalChar, tempCHPercentChar, "", -1)
					}

					tempCHPercentChar = string(chText[i:(i + 3)])
					tempTotalChar = tempTotalChar + tempCHPercentChar
					i = i + 2 //下次循环会默认加+1 所以要小心 +2
					continue

				}
			}

			//"""
			//看是不是点
			//"""
			charToGet = string(chText[i])
			_, exists = chineseConnectingSignDict[charToGet]
			if exists {
				//"""
				//如果 前一个符号赋值前，临时符号不为空，则把之前totalchar里面的符号替换为空字符串
				//"""
				if tempCHConnectChar != "" {
					tempTotalChar = strings.Replace(tempTotalChar, tempCHConnectChar, "", -1)
				}
				tempCHConnectChar = string(chText[i])
				tempTotalChar = tempTotalChar + tempCHConnectChar
				continue

			}

			//"""
			//看是不是数字
			//"""
			charToGet = string(chText[i])
			_, exists = chineseCharNumberDict[charToGet]
			if exists {
				//"""
				//如果 在字典里找到，则记录该字符串
				//"""
				tempCHNumberChar = string(chText[i])
				tempTotalChar = tempTotalChar + tempCHNumberChar
				continue

			} else {
				//"
				//遇到第一个在字典里找不到的，且最终长度大于符号与连接符的。所有临时记录清空, 最终字符串被记录
				//""
				if len(tempTotalChar) > len(tempCHPercentChar+tempCHConnectChar+tempCHSignChar) {
					if checkChineseNumberReasonable(tempTotalChar) {
						CHNumberStringList = append(CHNumberStringList, tempTotalChar)
						tempCHPercentChar = ""
						tempCHConnectChar = ""
						tempCHSignChar = ""
						tempCHNumberChar = ""
						tempTotalChar = ""
					} else {
						tempCHPercentChar = ""
						tempCHConnectChar = ""
						tempCHSignChar = ""
						tempCHNumberChar = ""
						tempTotalChar = ""
					}
				}

				//"
				//遇到第一个在字典里找不到的，且最终长度小于符号与连接符的。所有临时记录清空,。
				//""
			}

		}
		//"""
		//将temp 清干净
		//"""
		if len(tempTotalChar) > len(tempCHPercentChar+tempCHConnectChar+tempCHSignChar) {
			if checkChineseNumberReasonable(tempTotalChar) {
				CHNumberStringList = append(CHNumberStringList, tempTotalChar)
				tempCHPercentChar = ""
				tempCHConnectChar = ""
				tempCHSignChar = ""
				tempCHNumberChar = ""
				tempTotalChar = ""
			} else {
				tempCHPercentChar = ""
				tempCHConnectChar = ""
				tempCHSignChar = ""
				tempCHNumberChar = ""
				tempTotalChar = ""
			}
		}

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
