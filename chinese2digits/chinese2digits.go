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

// 中文转阿拉伯数字
var chineseCharNumberDict = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5,
	"六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "百": 100,
	"千": 1000, "万": 10000, "亿": 100000000, "壹": 1, "贰": 2, "叁": 3, "肆": 4, "伍": 5, "陆": 6, "柒": 7, "捌": 8, "玖": 9, "拾": 10,
	"佰": 100, "仟": 1000}

// var chinesePercentString = "百分之"
var chineseSignDict = map[string]string{"负": "-", "正": "+", "-": "-", "+": "+"}
var chineseConnectingSignDict = map[string]string{".": ".", "点": ".", "·": "."}

// var chinesePureNumberList = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10}

// """
// 阿拉伯数字转中文
// """
var digitsCharChineseDict = map[string]string{"0": "零", "1": "一", "2": "二", "3": "三", "4": "四", "5": "五", "6": "六", "7": "七", "8": "八", "9": "九", "%": "百分之", "‰": "千分之", "‱": "万分之", ".": "点"}

// var takingChineseNumberRERules, regError1 = regexp.Compile(`(?:(?:[正负]){0,1}(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+)))` +
// `(?:(?:分之)(?:[正负]){0,1}(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|` +
// `(?:点[一二三四五六七八九幺零]+))){0,1}`)

//数字汉字混合提取的正则引擎

// mac 系统对于下面的正则会报错  还在找原因
//runtime error: invalid memory address or nil pointer dereference [signal SIGSEGV: segmentation violation]
//Unable to propogate EXC_BAD_ACCESS signal to target process and panic (see https://github.com/go-delve/delve/issues/852)
//go 语言没有 ‰\‱  会报错
var takingChineseDigitsMixRERules, regError2 = regexp.Compile(`(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})` +
	`(?:(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1}){0,1}` +
	`(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))))` +
	`|(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1})` +
	`(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))`)

var PURE_DIGITS_RE, _ = regexp.Compile(`[0-9]`)

// var takingDigitsRERule, _ = regexp.Compile(`(?:\+|\-){0,1}\d+(?:\.\d+){0,1}(?:\%){0,1}|(?:\+|\-){0,1}\.\d+(?:\%){0,1}`)

func checkChineseNumberReasonable(chNumber string) bool {
	if len(chNumber) > 0 {
		// #由于在上个检查点 已经把阿拉伯数字转为中文 因此不用检查阿拉伯数字部分
		// """
		// 如果汉字长度大于0 则判断是不是 万  千  单字这种
		// """
		for i := 0; i < len(CHINESE_PURE_NUMBER_LIST); i++ {
			if strings.Contains(chNumber, CHINESE_PURE_NUMBER_LIST[i]) {
				// chNumberReasonable = true
				// break
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
func CoreCHToDigits(chineseCharsToTrans string) string {
	chineseChars := []rune(chineseCharsToTrans)
	total := ""
	tempVal := ""                      //#用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
	countingUnit := 1                  //#表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
	countingUnitFromString := []int{1} //#原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万

	tempTotal := 0
	// countingUnit := 1
	// 表示单位：个十百千...
	for i := len(chineseChars) - 1; i >= 0; i = i - 1 {
		charToGet := string(chineseChars[i])
		val := chineseCharNumberDict[charToGet]
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
			if i > 0 {
				// #如果下一个不是单位 则本次也是拼接
				preValTemp := chineseCharNumberDict[string(chineseChars[i-1])]
				if preValTemp < 10 {
					tempVal = strconv.Itoa(val) + tempVal
				} else {
					// #说明已经有大于10的单位插入 要数学计算了
					// #先拼接再计算
					// #如果取值不大于10 说明是0-9 则继续取值 直到取到最近一个大于10 的单位   应对这种30万20千 这样子
					tempValInt, err := strconv.Atoi(strconv.Itoa(val) + tempVal)
					if err != nil {
						panic(err)
					} else {
						tempTotal = tempTotal + countingUnit*tempValInt
					}
					// #计算后 把临时字符串置位空
					tempVal = ""
				}

			} else {
				// #那就是无论如何要收尾了
				//如果counting unit 等于1  说明所有字符串都是直接拼接的，不用计算，不然会丢失前半部分的零
				if countingUnit == 1 {
					tempVal = strconv.Itoa(val) + tempVal
				} else {
					tempValInt, err := strconv.Atoi(strconv.Itoa(val) + tempVal)
					if err != nil {
						panic(err)
					} else {
						tempTotal = tempTotal + countingUnit*tempValInt
					}
				}
			}

		}

	}

	//如果 total 为0  但是 countingUnit 不为0  说明结果是 十万这种  最终直接取结果 十万
	if tempTotal == 0 {
		if countingUnit > 10 {
			// 转化为字符串
			total = strconv.Itoa(countingUnit)
		} else {
			//counting Unit 为1  且tempval 不为空 说明是没有单位的纯数字拼接
			if tempVal != "" {
				total = tempVal
			} else {
				total = strconv.Itoa(tempTotal)
			}
		}
	} else {
		// 转化为字符串
		total = strconv.Itoa(tempTotal)

	}

	return total
}

// func convertDigitsStringToFloat(digitsString string) float64 {
// 	perCountingString := ""
// 	convertResult := digitsString
// 	for i := 0; i < len(CHINESE_PER_COUNTING_STRING_LIST); i++ {
// 		if strings.Contains(digitsString, CHINESE_PER_COUNTING_STRING_LIST[i]) {
// 			value, exists := CHINESE_PER_COUNTING_DICT[string(CHINESE_PER_COUNTING_STRING_LIST[i])]
// 			if exists {
// 				perCountingString = value
// 			} else {
// 				perCountingString = ""
// 			}
// 			convertResult = strings.Replace(digitsString, CHINESE_PER_COUNTING_STRING_LIST[i], "", -1)
// 		}
// 	}

// 	var finalTotal float64
// 	floatResult, err := strconv.ParseFloat(convertResult, 32)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		switch perCountingString {
// 		case "%":
// 			finalTotal = floatResult / 100
// 		case "‰":
// 			finalTotal = floatResult / 1000
// 		case "‱":
// 			finalTotal = floatResult / 10000
// 		default:
// 			finalTotal = floatResult
// 		}
// 	}
// 	return finalTotal
// }

//汉字切分是否正确  分之 切分
func checkNumberSeg(chineseNumberList []string, originText string) []string {
	var newChineseNumberList []string
	// #用来控制是否前一个已经合并过  防止多重合并
	tempPreText := ""
	tempMixedString := ""
	segLen := len(chineseNumberList)
	if segLen > 0 {
		// #加入唯一的一个 或者第一个
		// _, exists := chineseNumberList[0][:2]

		if strings.HasPrefix(chineseNumberList[0], "分之") {
			// #如果以分之开头 记录本次 防止后面要用 是否出现连续的 分之
			newChineseNumberList = append(newChineseNumberList, chineseNumberList[0][2:])
		} else {
			newChineseNumberList = append(newChineseNumberList, chineseNumberList[0])
		}

		if segLen > 1 {
			for i := 1; i < segLen; i++ {
				// #判断本字符是不是以  分之  开头
				if strings.HasPrefix(chineseNumberList[i], "分之") {
					// #如果是以 分之 开头 那么检查他和他见面的汉子数字是不是连续的 即 是否在原始字符串出现
					tempMixedString = string(chineseNumberList[i-1]) + string(chineseNumberList[i])
					runeTempPreText := []rune(tempPreText)
					if strings.Contains(originText, tempMixedString) {
						// #如果连续的上一个字段是以分之开头的  本字段又以分之开头
						if tempPreText != "" {
							// #检查上一个字段的末尾是不是 以 百 十 万 的单位结尾
							if isExistItem(string(runeTempPreText[len(runeTempPreText)-1]), CHINESE_PURE_COUNTING_UNIT_LIST) > -1 {
								//上一个字段最后一个数据
								tempLastChineseNumber := []rune(newChineseNumberList[len(newChineseNumberList)-1])
								// #先把上一个记录进去的最后一位去掉
								newChineseNumberList[len(newChineseNumberList)-1] = string(tempLastChineseNumber[:len(tempLastChineseNumber)-1])
								// #如果结果是确定的，那么本次的字段应当加上上一个字段的最后一个字
								newChineseNumberList = append(newChineseNumberList, string(runeTempPreText[len(runeTempPreText)-1])+chineseNumberList[i])
							} else {
								// #如果上一个字段不是以单位结尾  同时他又是以分之开头，那么 本次把分之去掉
								newChineseNumberList = append(newChineseNumberList, chineseNumberList[i][2:])
							}
						} else {
							// #上一个字段不以分之开头，那么把两个字段合并记录
							if len(newChineseNumberList) > 0 {
								//把最后一位赋值
								newChineseNumberList[len(newChineseNumberList)-1] = tempMixedString
							} else {
								newChineseNumberList = append(newChineseNumberList, tempMixedString)
							}
						}
					} else {
						// #说明前一个数字 和本数字不是连续的
						// #本数字去掉分之二字
						newChineseNumberList = append(newChineseNumberList, chineseNumberList[i][2:])
					}
					// #记录以 分之 开头的字段  用以下一个汉字字段判别
					tempPreText = chineseNumberList[i]
				} else {
					// #不是  分之 开头 那么把本数字加入序列
					newChineseNumberList = append(newChineseNumberList, chineseNumberList[i])
					// #记录把不是 分之 开头的字段  临时变量记为空
					tempPreText = ""
				}
			}
		}
	}
	return newChineseNumberList
}

// var dotRightPartReplaceRule, regError5 = regexp.Compile("0+$")

// ChineseToDigits 是可以识别包含百分号，正负号的函数，并控制是否将百分之10转化为0.1
func ChineseToDigits(chineseCharsToTrans string, percentConvert bool) string {
	// """
	// 分之  分号切割  要注意
	// """
	finalTotal := ""
	var convertResultList []string
	var chineseCharsListByDiv []string
	// if regError5 != nil {
	// 	panic(regError5)
	// }

	if strings.Contains(chineseCharsToTrans, "分之") {
		chineseCharsListByDiv = strings.Split(chineseCharsToTrans, "分之")

	} else {
		chineseCharsListByDiv = append(chineseCharsListByDiv, chineseCharsToTrans)
	}

	for k := 0; k < len(chineseCharsListByDiv); k++ {

		tempChineseChars := chineseCharsListByDiv[k]

		// chineseChars := tempChineseChars
		chineseChars := []rune(tempChineseChars)
		// """
		// 看有没有符号
		// """
		sign := ""
		for i := 0; i < len(chineseChars); i++ {
			charToGet := string(chineseChars[i])
			value, exists := chineseSignDict[charToGet]
			if exists {
				sign = value
				// chineseCharsToTrans = strings.Replace(chineseCharsToTrans, charToGet, "", -1)
				tempChineseChars = strings.Replace(tempChineseChars, charToGet, "", -1)
			}

		}
		// chineseChars = []rune(chineseCharsToTrans)

		// """
		// 小数点切割，看看是不是有小数点
		// """
		stringContainDot := false
		leftOfDotString := ""
		rightOfDotString := ""
		for key := range chineseConnectingSignDict {
			if strings.Contains(tempChineseChars, key) {
				chineseCharsDotSplitList := strings.Split(tempChineseChars, key)
				leftOfDotString = string(chineseCharsDotSplitList[0])
				rightOfDotString = string(chineseCharsDotSplitList[1])
				stringContainDot = true
				break
			}
		}
		convertResult := ""
		if !stringContainDot {
			convertResult = CoreCHToDigits(tempChineseChars)
		} else {
			convertResult = ""
			tempBuf := bytes.Buffer{}
			tempRightDigits := ""
			// #如果小数点右侧有 单位 比如 2.55万  4.3百万 的处理方式
			// #先把小数点右侧单位去掉
			tempCountString := ""
			listOfRight := []rune(rightOfDotString)
			for ii := len(listOfRight) - 1; ii >= 0; ii-- {
				if isExistItem(string(listOfRight[ii]), CHINESE_PURE_COUNTING_UNIT_LIST) > -1 {
					tempCountString = string(listOfRight[ii]) + tempCountString
				} else {
					rightOfDotString = string(listOfRight[0:(ii + 1)])
					break
				}
			}

			tempCountNum := 1.0
			if tempCountString != "" {
				tempNum, errTemp := strconv.ParseFloat(CoreCHToDigits(tempCountString), 32)
				if errTemp != nil {
					panic(errTemp)
				} else {
					tempCountNum = tempNum
				}
			}

			if leftOfDotString == "" {
				// """
				// .01234 这种开头  用0 补位
				// """
				tempBuf.WriteString("0.")
				tempRightDigits = CoreCHToDigits(rightOfDotString)
				// tempRightDigits = dotRightPartReplaceRule.ReplaceAllString(tempRightDigits, "")
				tempBuf.WriteString(tempRightDigits)
				convertResult = tempBuf.String()
			} else {
				tempBuf.WriteString(CoreCHToDigits(leftOfDotString))
				tempBuf.WriteString(".")
				tempRightDigits = CoreCHToDigits(rightOfDotString)
				// tempRightDigits = dotRightPartReplaceRule.ReplaceAllString(tempRightDigits, "")
				tempBuf.WriteString(tempRightDigits)
				convertResult = tempBuf.String()
			}

			tempStrToFloat, errTemp1 := strconv.ParseFloat(convertResult, 32)
			if errTemp1 != nil {
				panic(errTemp1)
			} else {
				convertResult = strconv.FormatFloat(tempStrToFloat*tempCountNum, 'f', -1, 32)
			}

		}
		//如果转换结果为空字符串 则为百分之10 这种
		if convertResult == "" {
			convertResult = "1"
		}
		convertResult = sign + convertResult
		// #最后在双向转换一下 防止出现 0.3000 或者 00.300的情况

		newConvertResultTemp := []rune(convertResult)
		newBuf := ""
		if strings.HasSuffix(convertResult, ".0") {
			newBuf = string(newConvertResultTemp[0 : len(newConvertResultTemp)-2])
		} else {
			newBuf = convertResult
		}
		convertResultList = append(convertResultList, newBuf)

	}
	if len(convertResultList) > 1 {
		// #是否转换分号及百分比
		if percentConvert {
			tempFloat1, err1 := strconv.ParseFloat(convertResultList[1], 32/64)
			if err1 != nil {
				panic(err1)
			} else {
				tempFloat0, err0 := strconv.ParseFloat(convertResultList[0], 32/64)
				if err0 != nil {
					panic(err0)
				} else {
					// fmt.Println(tempFloat1 / tempFloat0)
					finalTotal = strconv.FormatFloat(tempFloat1/tempFloat0, 'f', -1, 32)
				}
			}

		} else {
			if convertResultList[0] == "100" {
				finalTotal = convertResultList[1] + "%"
			} else if convertResultList[0] == "1000" {
				finalTotal = convertResultList[1] + "‰"
			} else {
				finalTotal = convertResultList[1] + "/" + convertResultList[0]
			}

		}

	} else {
		finalTotal = convertResultList[0]
		//最后再转换一下 防止出现 .50 的问题  不能转换了 否则  超出精度了………… 服了  5亿的话
		// tempFinalTotal, err3 := strconv.ParseFloat(finalTotal, 32)
		// if err3 != nil {
		// 	panic(err3)
		// } else {
		// 	finalTotal = strconv.FormatFloat(tempFinalTotal, 'f', -1, 32)
		// }
	}

	//删除最后的0
	if strings.Contains(finalTotal, ".") {
		finalTotal = strings.TrimSuffix(finalTotal, "0")
		finalTotal = strings.TrimSuffix(finalTotal, ".")
		// if strings.HasSuffix(finalTotal, ".") {

		// }
	}

	return finalTotal
}

//阿拉伯数字转中文
func digitsToCHChars(mixedStringList []string) []string {
	var resultList []string
	for i := 0; i < len(mixedStringList); i++ {
		mixedString := mixedStringList[i]
		if strings.HasPrefix(mixedString, ".") {
			mixedString = "0" + mixedString
		}
		for key := range digitsCharChineseDict {
			if strings.Contains(mixedString, key) {
				mixedString = strings.Replace(mixedString, key, digitsCharChineseDict[key], -1)
				// #应当是只要有百分号 就挪到前面 阿拉伯数字没有四百分之的说法
				// #防止这种 3%万 这种问题
				for kk := 0; kk < len(CHINESE_PER_COUNTING_STRING_LIST); kk++ {
					k := CHINESE_PER_COUNTING_STRING_LIST[kk]
					if strings.Contains(mixedString, k) {
						temp := k + strings.Replace(mixedString, k, "", -1)
						mixedString = temp
					}
				}
			}
		}
		resultList = append(resultList, mixedString)
	}

	return resultList
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
	if simplifConvertSwitch {
		for i := 0; i < stringLength; i++ {
			// #繁体中文数字转简体中文数字
			charToGet = string(chStringList[i])
			value, exists := TRADITIONAl_CONVERT_DICT[charToGet]
			if exists {
				chStringList[i] = []rune(value)[0]
			}
		}

	}
	if stringLength > 1 {
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
		fistCharCheckResult := isExistItem(string(tempNewChNumberStringList[0]), []string{"千", "万", "百"})
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
			if perCountSwitch {
				newChNumberStringList = string(tempNewChNumberStringList[:1]) + "分之" + string(tempNewChNumberStringList[1:])
			}
		}
	}

	return string(newChNumberStringList)
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
	if len(opt) > 3 {
		panic("too many arguments")
	}

	var percentConvert bool
	var traditionalConvert bool
	// var digitsNumberSwitch bool
	// digitsNumberSwitch := false

	switch len(opt) {
	case 1:
		percentConvert = opt[0].(bool)
		traditionalConvert = true
		// digitsNumberSwitch = false
	case 2:
		percentConvert = opt[0].(bool)
		traditionalConvert = opt[1].(bool)
		// digitsNumberSwitch = false
	// case 3:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = opt[1].(bool)
	// digitsNumberSwitch = opt[2].(bool)
	default:
		percentConvert = true
		traditionalConvert = true
	}

	// fmt.Println(digitsNumberSwitch)

	//"""
	//简体转换开关
	//"""
	// originText := chTextString

	convertedString := traditionalTextConvertFunc(chTextString, traditionalConvert)

	//正则引擎
	if regError2 != nil {
		fmt.Println(regError2)
		panic(regError2)
	}
	regMatchResult := takingChineseDigitsMixRERules.FindAllStringSubmatch(convertedString, -1)

	tempText := ""
	CHNumberStringListTemp := []string{}
	for i := 0; i < len(regMatchResult); i++ {
		tempText = regMatchResult[i][0]
		CHNumberStringListTemp = append(CHNumberStringListTemp, tempText)
	}

	// ##检查是不是  分之 切割不完整问题
	CHNumberStringListTemp = checkNumberSeg(CHNumberStringListTemp, convertedString)

	// 检查最后是正负号的问题
	CHNumberStringListTemp = checkSignSeg(CHNumberStringListTemp)

	// #备份一个原始的提取，后期处结果的时候显示用
	OriginCHNumberTake := CHNumberStringListTemp

	// #将阿拉伯数字变成汉字  不然合理性检查 以及后期 如果不是300万这种乘法  而是 四分之345  这种 就出错了
	CHNumberStringListTemp = digitsToCHChars(CHNumberStringListTemp)

	//检查合理性 是否是单纯的单位  等
	// var CHNumberStringList []string
	// var OriginCHNumberForOutput []string
	OriginCHNumberForOutput := []string{}
	for i := 0; i < len(CHNumberStringListTemp); i++ {
		// fmt.Println(aa[i])
		tempText = CHNumberStringListTemp[i]
		if checkChineseNumberReasonable(tempText) {
			// #如果合理  则添加进被转换列表
			CHNumberStringList = append(CHNumberStringList, tempText)
			// #则添加把原始提取的添加进来
			OriginCHNumberForOutput = append(OriginCHNumberForOutput, OriginCHNumberTake[i])
		}
		// CHNumberStringList = append(CHNumberStringList, regMatchResult[i][0])
	}

	// """
	// 进行标准汉字字符串转换 例如 二千二  转换成二千零二
	// """
	CHNumberStringListTemp = []string{}
	for i := 0; i < len(CHNumberStringList); i++ {
		CHNumberStringListTemp = append(CHNumberStringListTemp, standardChNumberConvert(CHNumberStringList[i]))

	}

	//"""
	//将中文转换为数字
	//"""
	digitsStringList := []string{}
	replacedText := convertedString
	tempCHToDigitsResult := ""
	CHNumberStringLenList := []int{}
	structCHAndDigitSlice := []structCHAndDigit{}
	if len(CHNumberStringListTemp) > 0 {
		for i := 0; i < len(CHNumberStringListTemp); i++ {
			tempCHToDigitsResult = ChineseToDigits(CHNumberStringListTemp[i], percentConvert)
			digitsStringList = append(digitsStringList, tempCHToDigitsResult)
			// CHNumberStringLenList = append(CHNumberStringLenList, len(CHNumberStringListTemp[i]))
			CHNumberStringLenList = append(CHNumberStringLenList, len(OriginCHNumberForOutput[i]))
			//将每次的新结构体附加至准备排序的
			structCHAndDigitSlice = append(structCHAndDigitSlice, structCHAndDigit{OriginCHNumberForOutput[i], digitsStringList[i], CHNumberStringLenList[i]})
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
		"inputText":          chTextString,
		"replacedText":       replacedText,
		"CHNumberStringList": OriginCHNumberForOutput,
		"digitsStringList":   digitsStringList,
	}
	return finalResult

}

// TakeNumberFromString will extract the chinese and digits number together from string. and return the convert result
// :param chText: chinese string
// :param percentConvert: convert percent simple. Default is True.  3% will be 0.03 in the result
// :param traditionalConvert: Switch to convert the Traditional Chinese character to Simplified chinese
// :return: Dict like result. 'inputText',replacedText','CHNumberStringList':CHNumberStringList,'digitsStringList'
func TakeNumberFromString(chTextString string, opt ...interface{}) interface{} {

	//默认参数设置
	if len(opt) > 2 {
		panic("too many arguments")
	}

	var percentConvert bool
	var traditionalConvert bool
	// digitsNumberSwitch := false

	switch len(opt) {
	case 1:
		percentConvert = opt[0].(bool)
		traditionalConvert = true
		// digitsNumberSwitch = false
	case 2:
		percentConvert = opt[0].(bool)
		traditionalConvert = opt[1].(bool)
		// digitsNumberSwitch = false
	// case 3:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = opt[1].(bool)
	// digitsNumberSwitch = opt[2].(bool)
	default:
		percentConvert = true
		traditionalConvert = true
	}
	finalResult := TakeChineseNumberFromString(chTextString, percentConvert, traditionalConvert)
	return finalResult
}

// func main() {
// 	fmt.Println(TakeChineseNumberFromString("三十万", true, true, true))
// 	fmt.Println("这个函数被调用了")
// }

// func main() {
// 	fmt.Println(TakeChineseNumberFromString("三十万", true, true, true))
// 	fmt.Println("这个函数被调用了")
// }

// func main() {
// 	reg1, err := regexp.Compile(`(?:(?:(?:百分之[正负]{0,1})|(?:[正负](?:百分之){0,1}))(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+)))|(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九幺零]+){0,1})|(?:点[一二三四五六七八九幺零]+))`)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(reg1.FindAllStringSubmatch("负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八百分之四十啦啦啦啊四万三千四百二", -1))

// }
