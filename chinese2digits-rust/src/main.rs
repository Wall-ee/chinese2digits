// 中文转阿拉伯数字
let data = vec![
	("幺": 1), 
	("零": 0),
	("一": 1), 
	("二": 2), 
	("两": 2), 
	("三": 3), 
	("四": 4), 
	("五": 5),
	("六": 6), 
	("七": 7), 
	("八": 8), 
	("九": 9), 
	("十": 10), 
	("百": 100),
	("千": 1000), 
	("万": 10000), 
	("亿": 100000000), 
	("壹": 1), 
	("贰": 2), 
	("叁": 3), 
	("肆": 4), 
	("伍": 5), 
	("陆": 6), 
	("柒": 7), 
	("捌": 8), 
	("玖": 9), 
	("拾": 10),
	("佰": 100), 
	("仟": 1000)

]

// var chineseCharNumberDict = map[string]int{"幺": 1, "零": 0, "一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5,
// 	"六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "百": 100,
// 	"千": 1000, "万": 10000, "亿": 100000000, "壹": 1, "贰": 2, "叁": 3, "肆": 4, "伍": 5, "陆": 6, "柒": 7, "捌": 8, "玖": 9, "拾": 10,
// 	"佰": 100, "仟": 1000}

let chineseCharNumberDict: HashMap<_, _> = data.into_iter().collect();



// CoreCHToDigits 是核心转化函数
fn CoreCHToDigits(chineseCharsToTrans:String, dotPosition:bool) -> String{

	// chineseChars := []rune(chineseCharsToTrans)
	let chineseChars = chineseCharsToTrans;
	let total = "";
	let tempVal = "";                      //#用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
	let countingUnit = 1;                  //#表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
	let mut countingUnitFromString = vec![1]; //#原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万

	tempTotal := 0
	// // countingUnit := 1
	// 表示单位：个十百千...
	if !dotPosition {
		//如果是小数点左边 正常执行 考虑各种单位等等
		for i in ((chineseChars.len() - 1)..0) {
			let charToGet = chineseChars.chars()[i].to_string()
			let val = chineseCharNumberDict.get(&charToGet)
			match val{
				if val >= 10 &&  i = 0 {
					// 应对 十三 十四 十*之类
					if val > countingUnit {
						countingUnit = val;
						tempTotal = tempTotal + val;
						countingUnitFromString.push(val);
					} else {
						countingUnitFromString.push(val);
						let maxValue = countingUnitFromString.iter().max();
						countingUnit = maxValue * val;

					}
				}else if val >=10 {
					if val > countingUnit {
						countingUnit = val;
						countingUnitFromString.push(val);
					} else {
						countingUnitFromString.push(val);
						let maxValue = countingUnitFromString.iter().max();
						countingUnit = maxValue * val;
					}

				}else{
					if i > 0 {
						// #如果下一个不是单位 则本次也是拼接
						preValTemp, _ := chineseCharNumberDict.get(&chineseChars[i-1].to_string());
						match preValTemp{
							if preValTemp < 10 {
								tempVal = strconv.Itoa(val) + tempVal;
								tempVal = val.to_string() + &tempVal;
							} else {
								// #说明已经有大于10的单位插入 要数学计算了
								// #先拼接再计算
								// #如果取值不大于10 说明是0-9 则继续取值 直到取到最近一个大于10 的单位   应对这种30万20千 这样子
								tempValInt = (val.to_string() + &tempVal).parse::<i32>().unwrap();
								tempTotal = tempTotal + countingUnit * tempValInt;
								// #计算后 把临时字符串置位空
								tempVal = "";
							}
						}


					} else {
						// #那就是无论如何要收尾了
						tempValInt = (val.to_string() + &tempVal).parse::<i32>().unwrap();
						tempTotal = tempTotal + countingUnit * tempValInt;
						// #计算后 把临时字符串置位空
						tempVal = ""
					}
				}			
				//如果 total 为0  但是 countingUnit 不为0  说明结果是 十万这种  最终直接取结果 十万
				if tempTotal = 0  && countingUnit > 10{
					total = countingUnit.to_string();
				}else{
					// 转化为字符串
					total = tempTotal.to_string();
				}
				None => println!( "No string in number" ),
			}
		}
	} else {
		//小数点右边，便捷执行，考虑 零零五六这种情况
		for i in ((len(chineseChars)-1) .. 0){
			let charToGet = chineseChars.chars()[i].to_string()
			let val = chineseCharNumberDict.get(&charToGet)
			match val {
				total = val.to_string() + &total;
			}
			
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
    print!(CoreCHToDigits("Hello, world!".to_string(),false));
}-