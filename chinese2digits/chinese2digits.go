package chinese2digits

import (
	"strings"
	"fmt"
	"bytes"
	"strconv"
	"sort"
)



var CHINESE_COUNTING_STRING = map[string]int{"十": 10, "百": 100, "千": 1000, "万": 10000, "亿": 100000000}

// 中文转阿拉伯数字

var CHINESE_CHAR_NUMBER_DICT = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "百": 100, "千": 1000, "万": 10000, "亿": 100000000}
var CHINESE_PERCENT_STRING = "百分之"
var CHINESE_SIGN_DICT = map[string]string{"负": "-", "正": "+", "-": "-", "+": "+"}
var CHINESE_CONNECTING_SIGN_DICT = map[string]string{".": ".", "点": ".", "·": "."}

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

	if (simpilfy == nil) {
		if (len(chineseChars) > 1) {
			// 如果字符串大于1 且没有单位 ，simplilfy is true. unsimpilify is false
			for i := 0; i < len(chineseChars); i++ {
				charToGet := string(chineseChars[i])
				_, exists := CHINESE_COUNTING_STRING[charToGet];
				if !exists {
					// fmt.Printf(strconv.Itoa(value2))
					simpilfySign = false;
				} else {
					simpilfySign = true;
					break;
				}

			}
		}
	}

	if (simpilfySign == true) {
		tempTotal := 0
		r := 1
		// 表示单位：个十百千...             
		for i := len(chineseChars) - 1; i >= 0; i = i - 1 {
			charToGet := string(chineseChars[i])
			val, _ := CHINESE_CHAR_NUMBER_DICT[charToGet]
			if ((val >= 10) && (i == 0)) {
				// 应对 十三 十四 十*之类
				if (val > r) {
					r = val
					tempTotal = tempTotal + val
				} else {
					r = r * val
				}
			} else if (val >= 10) {
				if (val > r) {
					r = val
				} else {
					r = r * val
				}
			} else {
				tempTotal = tempTotal + r*val
			}
		}
		// 转化为字符串
		total = strconv.Itoa(tempTotal)
	} else {
		// total:= ""
		tempBuf := bytes.Buffer{}
		for i := 0; i < len(chineseChars); i++ {
			charToGet := string(chineseChars[i])
			val, exisits := CHINESE_CHAR_NUMBER_DICT[charToGet]
			if (!exisits) {
			} else {
				tempBuf.WriteString(strconv.Itoa(val))
			}

		}
		total = tempBuf.String()
	}
	return total
}

func ChineseToDigits(chineseCharsToTrans string, percentConvert bool, simpilfy interface{}) string {

	chineseChars := []rune(chineseCharsToTrans)

	// """
	// 看有没有符号
	// """
	sign := ""
	for i := 0; i < len(chineseChars); i++ {
		charToGet := string(chineseChars[i])
		value, exists := CHINESE_SIGN_DICT[charToGet];
		if exists {
			sign = value
			chineseCharsToTrans = strings.Replace(chineseCharsToTrans, charToGet, "", -1)
		}

	}

	// """
	// 看有没有百分号
	// """
	chineseChars = []rune(chineseCharsToTrans)
	percentString := ""

	if strings.Contains(chineseCharsToTrans, CHINESE_PERCENT_STRING) {
		percentString = "%"
		chineseCharsToTrans = strings.Replace(chineseCharsToTrans, CHINESE_PERCENT_STRING, "", -1)
	}

	chineseChars = []rune(chineseCharsToTrans)

	// """
	// 小数点切割，看看是不是有小数点
	// """
	stringContainDot := false
	leftOfDotString := ""
	rightOfDotString := ""
	for key, _ := range CHINESE_CONNECTING_SIGN_DICT {
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

	finalTotal :=""
	if percentConvert == true {
		if percentString == "%" {
			floatResult, err := strconv.ParseFloat(convertResult, 32)
			if err != nil {
				panic(err)
			} else {
				finalTotal = strconv.FormatFloat(floatResult/100, 'f', -1, 32)
			}
		}
		return finalTotal

	} else {
		if percentString == "%" {
			finalTotal = convertResult + percentString
			return finalTotal
		} else {
			return finalTotal
		}

	}
}


type structCHAndDigit struct{
	CHNumberString string //中文数字字符串
	digitsString string //阿拉伯数字字符串
	CHNumberStringLen int //中文数字字符串长度
}


//结构体排序工具重写
// 按照 structToReplace.CHNumberStringLen 从大到小排序
type  structToReplace []structCHAndDigit

// 重写 Len() 方法
func (a structToReplace) Len()int {return len(a)}
// 重写 Swap() 方法
func (a structToReplace) Swap(i,j int) {a[i],a[j] = a[j],a[i] }
// 重写 Less() 方法， 从大到小排序
func (a structToReplace) Less(i,j int) bool{return a[j].CHNumberStringLen < a[i].CHNumberStringLen}



func TakeChineseNumberFromString(chTextString string,simpilfy interface {},percentConvert bool) interface{}{
	tempCHNumberChar :=""
	tempCHSignChar :=""
	tempCHConnectChar :=""
	tempCHPercentChar :=""
	CHNumberStringList := []string {}
	tempTotalChar := ""
	//"""
	//将字符串中所有中文数字列出来
	//"""
	chText := []rune(chTextString)
	for i:=0; i<len(chText);i++ {
		//"""
		//看是不是符号。如果是，就记录。
		//"""
		charToGet := string(chText[i])
		_, exists := CHINESE_SIGN_DICT[charToGet]
		if exists {
			//"""
			//如果 符号前面有数字  则 存到结果里面
			//"""
			if tempCHNumberChar != "" {
				CHNumberStringList = append(CHNumberStringList, tempTotalChar)
				tempCHPercentChar = ""
				tempCHConnectChar = ""
				tempCHSignChar = ""
				tempCHNumberChar = ""
				tempTotalChar = ""
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
		if string(chText[i:(i + 3)]) == CHINESE_PERCENT_STRING {
			//"""
			//如果 百分之前面有数字  则 存到结果里面
			//"""
			if tempCHNumberChar != "" {
				CHNumberStringList = append(CHNumberStringList, tempTotalChar)
				tempCHPercentChar = ""
				tempCHConnectChar = ""
				tempCHSignChar = ""
				tempCHNumberChar = ""
				tempTotalChar = ""
			}
			//"""
			//如果 前一个符号赋值前，临时符号不为空，则把之前totalchar里面的符号替换为空字符串
			//"""
			if tempCHPercentChar != "" {
				tempTotalChar = strings.Replace(tempTotalChar, tempCHPercentChar, "", -1)
			}

			tempCHPercentChar = string(chText[i:(i + 3)])
			tempTotalChar = tempTotalChar + tempCHPercentChar
			i = i + 3
			continue

		}

		//"""
		//看是不是点
		//"""
		charToGet = string(chText[i])
		_, exists = CHINESE_CONNECTING_SIGN_DICT[charToGet]
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
		_, exists = CHINESE_CHAR_NUMBER_DICT[charToGet]
		if exists {
			//"""
			//如果 在字典里找到，则记录该字符串
			//"""
			tempCHNumberChar = string(chText[i])
			tempTotalChar = tempTotalChar + tempCHNumberChar
			continue

		}else {
			//"
			//遇到第一个在字典里找不到的，且最终长度大于符号与连接符的。所有临时记录清空, 最终字符串被记录
			//""
			if len(tempTotalChar) > len(tempCHPercentChar+tempCHConnectChar+tempCHSignChar){
				CHNumberStringList = append(CHNumberStringList,tempTotalChar)
				tempCHPercentChar = ""
				tempCHConnectChar = ""
				tempCHSignChar = ""
				tempCHNumberChar = ""
				tempTotalChar = ""
			}

			//"
			//遇到第一个在字典里找不到的，且最终长度小于符号与连接符的。所有临时记录清空,。
			//""
		}

	}
	//"""
	//将temp 清干净
	//"""
	if len(tempTotalChar) > len(tempCHPercentChar+tempCHConnectChar+tempCHSignChar){
		CHNumberStringList = append(CHNumberStringList,tempTotalChar)
		tempCHPercentChar = ""
		tempCHConnectChar = ""
		tempCHSignChar = ""
		tempCHNumberChar = ""
		tempTotalChar = ""
	}
	//"""
	//将中文转换为数字
	//"""
	digitsStringList := []string{}
	replacedText := chTextString
	tempCHToDigitsResult :=""
	CHNumberStringLenList := []int{}
	structCHAndDigitSlice := []structCHAndDigit{}
	if len(CHNumberStringList)>0{
		for i:=0;i<len(CHNumberStringList);i++{
			tempCHToDigitsResult = ChineseToDigits(CHNumberStringList[i],percentConvert,simpilfy)
			digitsStringList = append(digitsStringList,tempCHToDigitsResult)
			CHNumberStringLenList = append(CHNumberStringLenList, len(CHNumberStringList[i]))
			//将每次的新结构体附加至准备排序的
			structCHAndDigitSlice = append(structCHAndDigitSlice,structCHAndDigit{CHNumberStringList[i],digitsStringList[i],CHNumberStringLenList[i]})
		}
		//fmt.Println(structCHAndDigitSlice)


		sort.Sort(structToReplace(structCHAndDigitSlice))    // 按照 中文数字字符串长度 的逆序排序
		//"""
		//按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
		//"""
		for i:=0;i<len(CHNumberStringLenList);i++{
			replacedText = strings.Replace(replacedText,structCHAndDigitSlice[i].CHNumberString,structCHAndDigitSlice[i].digitsString,-1)
		}

	}
	finalResult := map [string]interface{}{
		"inputText":chTextString,
		"replacedText":replacedText,
		"CHNumberStringList":CHNumberStringList,
		"digitsStringList":digitsStringList,
	}
	return finalResult


	}

func main() {
	fmt.Println(TakeChineseNumberFromString("负百分之点二八百分之三五", true,true))
	fmt.Println("这个函数被调用了")
}