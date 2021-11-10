
// CoreCHToDigits 是核心转化函数
fn CoreCHToDigits(chineseCharsToTrans:String, dotPosition:bool) String {
	chineseChars := []rune(chineseCharsToTrans)
	total := ""
	tempVal := ""                      //#用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
	countingUnit := 1                  //#表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
	countingUnitFromString := []int{1} //#原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万

	tempTotal := 0
	// countingUnit := 1
	// 表示单位：个十百千...
	if dotPosition == 0 {
		//如果是小数点左边 正常执行 考虑各种单位等等
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
				if i > 0 {
					// #如果下一个不是单位 则本次也是拼接
					preValTemp, _ := chineseCharNumberDict[string(chineseChars[i-1])]
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
					tempValInt, err := strconv.Atoi(strconv.Itoa(val) + tempVal)
					if err != nil {
						panic(err)
					} else {
						tempTotal = tempTotal + countingUnit*tempValInt
					}
					// #计算后 把临时字符串置位空
					tempVal = ""
				}

			}
			//如果 total 为0  但是 countingUnit 不为0  说明结果是 十万这种  最终直接取结果 十万
			if (tempTotal == 0) && (countingUnit) > 10 {
				// 转化为字符串
				total = strconv.Itoa(countingUnit)

			} else {
				// 转化为字符串
				total = strconv.Itoa(tempTotal)
			}
		}
	} else {
		//小数点右边，便捷执行，考虑 零零五六这种情况
		for i := len(chineseChars) - 1; i >= 0; i = i - 1 {
			charToGet := string(chineseChars[i])
			val, _ := chineseCharNumberDict[charToGet]
			total = strconv.Itoa(val) + total
		}
	}
	newTotalTemp := []rune(total)
	newTotal := ""
	if strings.HasSuffix(total, ".0") {
		newTotal = string((newTotalTemp[0 : len(newTotalTemp)-2]))
	} else {
		newTotal = total
	}
	return newTotal
}



fn main() {
    println!("Hello, world!");
}
