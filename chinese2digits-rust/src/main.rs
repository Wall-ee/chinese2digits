// #[warn(unused_assignments)]
// #[macro_use] extern crate lazy_static;
extern crate regex;
use std::collections::HashMap;
use regex::Regex;


// 中文转阿拉伯数字
static CHINESE_CHAR_NUMBER_LIST: [(&str, i32); 29] = [
	("幺" , 1), 
	("零" , 0),
	("一" , 1), 
	("二" , 2), 
	("两" , 2), 
	("三" , 3), 
	("四" , 4), 
	("五" , 5),
	("六" , 6), 
	("七" , 7), 
	("八" , 8), 
	("九" , 9), 
	("十" , 10), 
	("百" , 100),
	("千" , 1000), 
	("万" , 10000), 
	("亿" , 100000000), 
	("壹" , 1), 
	("贰" , 2), 
	("叁" , 3), 
	("肆" , 4), 
	("伍" , 5), 
	("陆" , 6), 
	("柒" , 7), 
	("捌" , 8), 
	("玖" , 9), 
	("拾" , 10),
	("佰" , 100), 
	("仟" , 1000)

];

static CHINESE_PURE_COUNTING_UNIT_LIST:[&str;5] = ["十", "百", "千", "万", "亿"];

static CHINESE_PURE_NUMBER_LIST :[&str;13]= ["幺", "一", "二", "两", "三", "四", "五", "六", "七", "八", "九", "十", "零"];

// CoreCHToDigits 是核心转化函数
fn core_ch_to_digits(chinese_chars_to_trans:String) -> String{

	let chinese_chars = chinese_chars_to_trans;
	let total:String;
	let mut temp_val = "".to_string();                      //#用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
	let mut counting_unit:i32 = 1;                  //#表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
	let mut counting_unit_from_string: Vec<i32>= vec![1]; //#原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万
	let mut temp_val_int:i32;
	let  chinese_char_number_dict: HashMap<&str, i32> = CHINESE_CHAR_NUMBER_LIST.into_iter().collect();

	let mut temp_total:i32 = 0;
	// // counting_unit := 1
	// 表示单位：个十百千...
	for i in (0..chinese_chars.chars().count()).rev() {
		// println!("{}",i);
		let char_to_get = chinese_chars.chars().nth(i).unwrap().to_string();
		let val_from_hash_map = chinese_char_number_dict.get(char_to_get.as_str());
		
		match val_from_hash_map{
			Some(val_from_hash_map) => {
				let val = *val_from_hash_map;
				if val >= 10 &&  i == 0_usize{
					// 应对 十三 十四 十*之类
					if val > counting_unit {
						counting_unit = val;
						temp_total = temp_total + val;
						counting_unit_from_string.push(val);
					} else {
						counting_unit_from_string.push(val);
						let max_value_option = counting_unit_from_string.iter().max();
						let max_value = max_value_option.unwrap();
						counting_unit = max_value * val;

					}
				}else if val >= 10 {
					if val > counting_unit {
						counting_unit = val;
						counting_unit_from_string.push(val);
					} else {
						counting_unit_from_string.push(val);
						// let max_value = counting_unit_from_string.iter().max();
						let max_value_option = counting_unit_from_string.iter().max();
						let max_value = max_value_option.unwrap();
						counting_unit = max_value * val;
					}

				}else{
					if i > 0 {
						// #如果下一个不是单位 则本次也是拼接
						let pre_val_temp_option = chinese_char_number_dict.get(chinese_chars.chars().nth(i-1).unwrap().to_string().as_str());
						match pre_val_temp_option{
							Some(pre_val_temp_option) => {
								let pre_val_temp = *pre_val_temp_option;
								if pre_val_temp < 10 {
									// temp_val = strconv.Itoa(val) + temp_val;
									temp_val = val.to_string() + &temp_val;
								} else {
									// #说明已经有大于10的单位插入 要数学计算了
									// #先拼接再计算
									// #如果取值不大于10 说明是0-9 则继续取值 直到取到最近一个大于10 的单位   应对这种30万20千 这样子
									temp_val_int = (val.to_string() + &temp_val).parse::<i32>().unwrap();
									temp_total = temp_total + counting_unit * temp_val_int;
									// #计算后 把临时字符串置位空
									temp_val = "".to_string();
								}
							}
							None => {

							}
						}


					} else {
						// #那就是无论如何要收尾了
						//如果counting unit 等于1  说明所有字符串都是直接拼接的，不用计算，不然会丢失前半部分的零
						if counting_unit == 1 {
							temp_val = val.to_string() + &temp_val;
						}else{
							temp_val_int = (val.to_string() + &temp_val).parse::<i32>().unwrap();
							temp_total = temp_total + counting_unit * temp_val_int;
						}
					}
				}			
			}
			None => println!( "No string in number" ),
		}
	}
	//如果 total 为0  但是 counting_unit 不为0  说明结果是 十万这种  最终直接取结果 十万
	if temp_total == 0 {
		if counting_unit > 10 {
			// 转化为字符串
			total = counting_unit.to_string();
		}else{
			//counting Unit 为1  且tempval 不为空 说明是没有单位的纯数字拼接
			if temp_val != ""{
				total = temp_val;
			}else{
				// 转化为字符串
				total = temp_total.to_string();
			}
		}
	}else{
		// 转化为字符串
		total = temp_total.to_string();
	}

	return total
}


//汉字切分是否正确  分之 切分
fn check_number_seg(chinese_number_list:Vec<String>, origin_text:String) -> Vec<String> {
	let mut new_chinese_number_list :Vec<String> = [].to_vec();
	// #用来控制是否前一个已经合并过  防止多重合并
	let mut temp_pre_text :String = "".to_string();
	let mut temp_mixed_string:String;
	let seg_len = chinese_number_list.len();
	if seg_len > 0 {
		// #加入唯一的一个 或者第一个
		if chinese_number_list[0].starts_with("分之") {
			// #如果以分之开头 记录本次 防止后面要用 是否出现连续的 分之
			new_chinese_number_list.push((&chinese_number_list[0][2..(chinese_number_list[0].chars().count())]).to_string());
		} else {
			new_chinese_number_list.push(chinese_number_list.get(0).unwrap().to_string())
		}

		if seg_len > 1 {
			for i in 1..seg_len {
				// #判断本字符是不是以  分之  开头
				if (chinese_number_list[i]).starts_with("分之") {
					// #如果是以 分之 开头 那么检查他和他见面的汉子数字是不是连续的 即 是否在原始字符串出现
					temp_mixed_string = chinese_number_list.get(i-1).unwrap().to_string() + &chinese_number_list[i];
				
					if origin_text.contains(&temp_mixed_string.to_string()){
						// #如果连续的上一个字段是以分之开头的  本字段又以分之开头
						if temp_pre_text != "" {
							// #检查上一个字段的末尾是不是 以 百 十 万 的单位结尾
							if CHINESE_PURE_COUNTING_UNIT_LIST.iter().any(|&x| x == (&temp_pre_text.chars().last().unwrap().to_string())){
								//上一个字段最后一个数据
								let temp_last_chinese_number = new_chinese_number_list.last().unwrap().to_string();
								// #先把上一个记录进去的最后一位去掉
								let temp_new_chinese_number_list_len = new_chinese_number_list.len();
								new_chinese_number_list[temp_new_chinese_number_list_len-1] = (&temp_last_chinese_number[0..(temp_last_chinese_number.chars().count()-1)]).to_string();
								// #如果结果是确定的，那么本次的字段应当加上上一个字段的最后一个字
								new_chinese_number_list.push(temp_pre_text.chars().nth(temp_pre_text.chars().count()-1).unwrap().to_string() + &chinese_number_list[i]);
								
							} else {
								// #如果上一个字段不是以单位结尾  同时他又是以分之开头，那么 本次把分之去掉
								new_chinese_number_list.push((&chinese_number_list[i][2..]).to_string())
							}
						} else {
							// #上一个字段不以分之开头，那么把两个字段合并记录
							if new_chinese_number_list.len() > 0 {
								//把最后一位赋值
								// newChineseNumberList[len(newChineseNumberList)-1] = tempMixedString
								let temp_new_chinese_number_list_len = new_chinese_number_list.len();
								new_chinese_number_list[temp_new_chinese_number_list_len-1] = temp_mixed_string;
							} else {
								new_chinese_number_list.push(temp_mixed_string);
							}
						}
					} else {
						// #说明前一个数字 和本数字不是连续的
						// #本数字去掉分之二字
						new_chinese_number_list.push((&chinese_number_list[i][2..]).to_string())
					}
					// #记录以 分之 开头的字段  用以下一个汉字字段判别
					temp_pre_text = chinese_number_list.get(i).unwrap().to_string();
				} else {
					// #不是  分之 开头 那么把本数字加入序列
					new_chinese_number_list.push(chinese_number_list.get(i).unwrap().to_string());
					// #记录把不是 分之 开头的字段  临时变量记为空
					temp_pre_text = "".to_string();
				}
			}
		}
	}
	return new_chinese_number_list
}


fn check_chinese_number_reasonable(ch_number:String) ->bool {
	if ch_number.chars().count() > 0 {
		// #由于在上个检查点 已经把阿拉伯数字转为中文 因此不用检查阿拉伯数字部分
		// """
		// 如果汉字长度大于0 则判断是不是 万  千  单字这种
		// """
		for i in CHINESE_PURE_NUMBER_LIST{
			if ch_number.contains(&i.to_string()) {
				return true
			}
		}
	}
	return false
}


// """
// 阿拉伯数字转中文
// """
static DIGITS_CHAR_CHINESE_DICT_KEY:[&str;14] = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "%", "‰", "‱", "."];
static DIGITS_CHAR_CHINESE_DICT: [(&str, &str); 14] = [("0", "零"), ("1", "一"), ("2", "二"), ("3", "三"),( "4", "四"),( "5", "五"), ("6", "六"), ("7", "七"), ("8", "八"), ("9", "九"), ("%", "百分之"), ("‰", "千分之"), ("‱", "万分之"), (".", "点")];

static CHINESE_PER_COUNTING_STRING_LIST:[&str;3] = ["百分之", "千分之", "万分之"];

static TRADITIONAl_CONVERT_DICT: [(&str, &str); 9] = [("壹", "一"), ("贰", "二"), ("叁", "三"), ("肆", "四"), ("伍", "五"), ("陆", "六"), ("柒", "七"),("捌", "八"), ("玖", "九")];
static SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT :[(&str, &str);5]= [("拾", "十"), ("佰", "百"), ("仟", "千"), ("萬", "万"), ("億", "亿")];
static SPECIAL_NUMBER_CHAR_DICT: [(&str, &str); 2] = [("两", "二"), ("俩", "二")];
// static CHINESE_PURE_NUMBER_LIST: [&str; 13] = ["幺", "一", "二", "两", "三", "四", "五", "六", "七", "八", "九", "十", "零"];
static CHINESE_SIGN_LIST: [&str; 4] = ["正", "负", "+", "-"];



//阿拉伯数字转中文
fn digits_to_ch_chars(mixedStringList:Vec<String>) -> Vec<String> {

	let mut resultList :Vec<String>= vec![];
	let  digitsCharChineseDict: HashMap<&str, &str> = DIGITS_CHAR_CHINESE_DICT.into_iter().collect();
	// for i := 0; i < len(mixedStringList); i++ {
	for i in mixedStringList.iter() {
		let mut mixedString = i.to_string();
		if mixedString.starts_with(".") {
			mixedString = "0".to_string() + &mixedString;
		}
		for key in DIGITS_CHAR_CHINESE_DICT_KEY.iter() {
			if mixedString.contains(&key.to_string()) {
				mixedString = mixedString.replace(key, digitsCharChineseDict.get(key).unwrap().to_string().as_str());
				// #应当是只要有百分号 就挪到前面 阿拉伯数字没有四百分之的说法
				// #防止这种 3%万 这种问题
				for kk in CHINESE_PER_COUNTING_STRING_LIST.iter(){
					if mixedString.contains(&kk.to_string()){
						mixedString = kk.to_string() + &mixedString.replace(kk,"");
					}

				}

			}
		}
		resultList.push(mixedString);
	}

	return resultList
}

// """
// 繁体简体转换 及  单位  特殊字符转换 两千变二千
// """
fn traditionalTextConvertFunc(chString:String, simplifConvertSwitch:bool)->String {
	// chStringList := []rune(chString)
	let mut chStringList:Vec<char> = chString.chars().collect();

	let stringLength = chString.chars().count();
	let mut charToGet:String;
	let traditionalConvertDict: HashMap<&str, &str> = DIGITS_CHAR_CHINESE_DICT.into_iter().collect();
	let specialTraditionalCountingUnitCharDict :HashMap<&str,&str> = SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT.into_iter().collect();
	let specialNumberCharDict:HashMap<&str,&str> = SPECIAL_NUMBER_CHAR_DICT.into_iter().collect();
	
	if simplifConvertSwitch {
		for i in 0..chStringList.len(){
			// #繁体中文数字转简体中文数字
			// charToGet = string(chStringList[i])
			charToGet = chStringList[i].to_string();
			let value= traditionalConvertDict.get(charToGet.as_str());
			match value{
				Some(value) =>{
					chStringList[i] = (*value).chars().nth(0).unwrap();
				}
				None => {
				}
			}
		}

	}
	if stringLength > 1 {
		// #检查繁体单体转换
		for i in 0..stringLength {
			// #如果 前后有 pure 汉字数字 则转换单位为简体
			charToGet = chStringList[i].to_string();
			// value, exists := SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT[charToGet]
			let value = specialTraditionalCountingUnitCharDict.get(charToGet.as_str());
			// # 如果前后有单纯的数字 则进行单位转换
			match value {
				Some(value) =>{
					if i == 0{
						if CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x==chStringList[i+1].to_string()){
							chStringList[i] = (*value).chars().nth(0).unwrap();
						}

					}else if i == (stringLength - 1){
						if CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x==chStringList[i-1].to_string()){
							chStringList[i] = (*value).chars().nth(0).unwrap();
						}

					}else{
						if (CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x==chStringList[i-1].to_string())) ||
						(CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x==chStringList[i+1].to_string())) {
							chStringList[i] = (*value).chars().nth(0).unwrap();
						}

					}

				}
				None =>{

				}

			}
			// #特殊变换 俩变二
			charToGet = chStringList[i].to_string();
			let value = specialNumberCharDict.get(charToGet.as_str());
			// # 如果前后有单纯的数字 则进行单位转换
			match value {
				Some(value) => {
					if i == 0 {
						
						if CHINESE_PURE_COUNTING_UNIT_LIST.iter().any(|&x| x==chStringList[i+1].to_string()) {
							chStringList[i] = (*value).chars().nth(0).unwrap();
						}
					}else if i == (stringLength - 1){
						if CHINESE_PURE_COUNTING_UNIT_LIST.iter().any(|&x| x==chStringList[i-1].to_string()) {
							chStringList[i] = (*value).chars().nth(0).unwrap();
						}
					}else{
						if (CHINESE_PURE_COUNTING_UNIT_LIST.iter().any(|&x| x==chStringList[i-1].to_string())) ||
						(CHINESE_PURE_COUNTING_UNIT_LIST.iter().any(|&x| x==chStringList[i+1].to_string())) {
							chStringList[i] = (*value).chars().nth(0).unwrap();
					}

					}

				}
				None =>{

				}
			}

		}

	}
	let final_ch_string_list:String = chStringList.iter().collect();

	return final_ch_string_list
}

// """
// 标准表述转换  三千二 变成 三千零二  三千十二变成 三千零一十二
// """
fn standardChNumberConvert(chNumberString:String) -> String{
	let chNumberStringList:Vec<char> = chNumberString.chars().collect();

	let mut newChNumberStringList:String = chNumberString;

	// #大于2的长度字符串才有检测和补位的必要
	if chNumberStringList.len() > 2 {
		// #十位补一：
		let tenNumberIndex = chNumberStringList.iter().position(|&x| x=="十".chars().nth(0).unwrap());
		match tenNumberIndex{
			Some(tenNumberIndex) => {
				if tenNumberIndex == 0_usize{
					newChNumberStringList = "一".to_string() + &newChNumberStringList;
				}else{
					// # 如果没有左边计数数字 插入1
					if !(CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x == chNumberStringList[(tenNumberIndex-1)].to_string().as_str())){

						let tempLeftPart:String = chNumberStringList[0..tenNumberIndex].iter().collect();
						let tempRightPart:String = chNumberStringList[tenNumberIndex..].iter().collect();
						newChNumberStringList =  tempLeftPart + "一" + &tempRightPart;
					}

				}
				
			}
			None =>{

			}
		}
		

		// #差位补零
		// #逻辑 如果最后一个单位 不是十结尾 而是百以上 则数字后面补一个比最后一个出现的单位小一级的单位
		// #从倒数第二位开始看,且必须是倒数第二位就是单位的才符合条件

		let tempNewChNumberStringList:Vec<char> = newChNumberStringList.chars().collect();

		let lastCountingUnit = CHINESE_PURE_COUNTING_UNIT_LIST.iter().position(|&x| x==tempNewChNumberStringList[tempNewChNumberStringList.len()-2].to_string().as_str());
		// # 如果最末位的是百开头
		match lastCountingUnit{
			Some(lastCountingUnit) =>{
				if lastCountingUnit >= 1_usize {
					// # 则字符串最后拼接一个比最后一个单位小一位的单位 例如四万三 变成四万三千
					// # 如果最后一位结束的是亿 则补千万
					if lastCountingUnit == 4_usize {
						newChNumberStringList = newChNumberStringList + &"千万".to_string();
					} else {
						newChNumberStringList = newChNumberStringList + &CHINESE_PURE_COUNTING_UNIT_LIST[lastCountingUnit-1].to_string();
		
					}
		
				}
			}
			None =>{

			}
		}


	}
	//大于一的检查是不是万三，千四五这种
	let mut perCountSwitch = false;
	let tempNewChNumberStringList:Vec<char> = newChNumberStringList.chars().collect();
	if tempNewChNumberStringList.len() > 1 {
		// #十位补一：
		// fistCharCheckResult := isExistItem(string(tempNewChNumberStringList[0]), []string{"千", "万", "百"})

		let firstCharCheckResult = ["千", "万", "百"].iter().position(|&x| x==tempNewChNumberStringList[0].to_string().as_str());
		match firstCharCheckResult {
			Some(firstCharCheckResult) =>{
				for i in 1..tempNewChNumberStringList.len(){
					// #其余位数都是纯数字 才能执行
					if CHINESE_PURE_NUMBER_LIST.iter().any(|&x| x == tempNewChNumberStringList[i].to_string().as_str()){
						perCountSwitch = true;
					} else{
						perCountSwitch = false;
						//有一个不合适 退出循环
						break
					}
				}
				if perCountSwitch {
					let tempLeftPartString:String = tempNewChNumberStringList[0..1].iter().collect();
					let tempRightPartString:String = tempNewChNumberStringList[1..].iter().collect();
					newChNumberStringList =  tempLeftPartString + "分之" + &tempRightPartString;
				}

			}
			None =>{

			}
		} 
	}

	return newChNumberStringList
}

//检查初次提取的汉字数字是正负号是否切分正确
fn checkSignSeg(chineseNumberList:Vec<String>) -> Vec<String> {
	let mut newChineseNumberList:Vec<String> = vec![];
	let mut tempSign = "".to_string();
	for i in 0..chineseNumberList.len() {
		// #新字符串 需要加上上一个字符串 最后1位的判断结果
		let forTempSign:String= tempSign;
		let mut newChNumberString = forTempSign + &chineseNumberList[i];
		let tempChineseNumberList:Vec<char> = newChNumberString.chars().collect();
		if tempChineseNumberList.len() > 1 {
			let lastString:String = tempChineseNumberList[(tempChineseNumberList.len()-1)..].iter().collect();
			// #如果最后1位是百分比 那么本字符去掉最后三位  下一个数字加上最后1位
			if CHINESE_SIGN_LIST.iter().any(|&x| x==lastString.as_str()){
				tempSign = lastString;
				// #如果最后1位 是  那么截掉最后1位
				// let tempWithoutLastPartString:String = 
				newChNumberString = tempChineseNumberList[..(tempChineseNumberList.len()-1)].iter().collect();
			} else {
				tempSign = "".to_string();
			}
		}
		newChineseNumberList.push(newChNumberString)
	}
	return newChineseNumberList
}

struct FinalResultStruct {
	inputText:          String,
	replacedText:       String,
	chNumberStringList: Vec<String>,
	digitsStringList:  Vec<String>,
}


static takingChineseDigitsMixRERulesString:&str = r"(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})\
(?:(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1}){0,1}\
(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))))\
|(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1})\
(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))";


// lazy_static! {
// 	static ref init:String="".to_string();
// 	static ref takingChineseDigitsMixRERules:Regex = Regex::new(r"(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})" +
// 	r"(?:(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1}){0,1}" +
// 	r"(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))))" +
// 	r"|(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1})" +
// 	r"(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))").unwrap();
// }




static CHINESE_SIGN_DICT_STATIC = map[string]string{"负": "-", "正": "+", "-": "-", "+": "+"}
// ChineseToDigits 是可以识别包含百分号，正负号的函数，并控制是否将百分之10转化为0.1
fn chineseToDigits(chineseCharsToTrans:String, percentConvert:bool) -> String{
	// """
	// 分之  分号切割  要注意
	// """
	let mut finalTotal = "".to_string();
	let mut convertResultList:Vec<String> = vec![];
	let mut chineseCharsListByDiv:Vec<String> = vec![];



	let dotRightPartReplaceRule:Regex = Regex::new("0+$").unwrap();

	if chineseCharsToTrans.contains("分之"){
		let tempSplitResult = chineseCharsToTrans.split("分之");
		for s in tempSplitResult { 
			chineseCharsListByDiv.push(s.to_string()); 
		} 
	}else{
		chineseCharsListByDiv.push(chineseCharsToTrans);
	}

	for k in 0..chineseCharsListByDiv.len(){

		let mut tempChineseChars:String = chineseCharsListByDiv[k];

		// chineseChars := tempChineseChars
		let chineseChars:Vec<char> = tempChineseChars.chars().collect();
		// """
		// 看有没有符号
		// """
		let mut sign = "".to_string();
		for i in 0..chineseChars.len(){
			let charToGet = chineseChars[i].to_string().as_str();
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
				tempRightDigits = dotRightPartReplaceRule.ReplaceAllString(tempRightDigits, "")
				tempBuf.WriteString(tempRightDigits)
				convertResult = tempBuf.String()
			} else {
				tempBuf.WriteString(CoreCHToDigits(leftOfDotString))
				tempBuf.WriteString(".")
				tempRightDigits = CoreCHToDigits(rightOfDotString)
				tempRightDigits = dotRightPartReplaceRule.ReplaceAllString(tempRightDigits, "")
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

	return finalTotal
}

// TakeChineseNumberFromString 将句子中的汉子数字提取的整体函数
fn TakeChineseNumberFromString(chTextString:String,percentConvert:bool, traditionalConvert:bool) -> FinalResultStruct {

	let mut chNumberStringList:Vec<String> = vec![];

	//默认参数设置
	// if len(opt) > 3 {
	// 	panic("too many arguments")
	// }

	// var percentConvert bool
	// var traditionalConvert bool
	// // var digitsNumberSwitch bool
	// // digitsNumberSwitch := false

	// switch len(opt) {
	// case 1:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = true
	// 	// digitsNumberSwitch = false
	// case 2:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = opt[1].(bool)
	// 	// digitsNumberSwitch = false
	// // case 3:
	// // 	percentConvert = opt[0].(bool)
	// // 	traditionalConvert = opt[1].(bool)
	// // digitsNumberSwitch = opt[2].(bool)
	// default:
	// 	percentConvert = true
	// 	traditionalConvert = true
	// }

	// fmt.Println(digitsNumberSwitch)

	//"""
	//简体转换开关
	//"""
	// originText := chTextString


	// let traditionalConvert = true;
	let convertedString :String = traditionalTextConvertFunc(chTextString, traditionalConvert);

	//正则引擎
	let mut regMatchResult:Vec<String>;
	let takingChineseDigitsMixRERules:Regex = Regex::new(takingChineseDigitsMixRERulesString).unwrap();
	for cap in  takingChineseDigitsMixRERules.captures_iter(convertedString.as_str()){
		regMatchResult.push(cap[0].to_string());
	}

	let mut tempText = "".to_string();
	let mut chNumberStringListTemp :Vec<String> = vec![];
	for i in 0..regMatchResult.len(){
		tempText = regMatchResult[i];
		chNumberStringListTemp.push(tempText);
	}

	// ##检查是不是  分之 切割不完整问题
	chNumberStringListTemp = check_number_seg(chNumberStringListTemp, convertedString);

	// 检查最后是正负号的问题
	chNumberStringListTemp = checkSignSeg(chNumberStringListTemp);

	// #备份一个原始的提取，后期处结果的时候显示用
	let originCHNumberTake :Vec<String> = chNumberStringListTemp;

	// #将阿拉伯数字变成汉字  不然合理性检查 以及后期 如果不是300万这种乘法  而是 四分之345  这种 就出错了
	chNumberStringListTemp = digits_to_ch_chars(chNumberStringListTemp);

	//检查合理性 是否是单纯的单位  等
	// var chNumberStringList []string
	// var originCHNumberForOutput []string
	let mut originCHNumberForOutput :Vec<String> =  vec![];
	for i in  0..chNumberStringListTemp.len(){
		// fmt.Println(aa[i])
		tempText = chNumberStringListTemp[i];
		if check_chinese_number_reasonable(tempText) {
			// #如果合理  则添加进被转换列表
			chNumberStringList.push(tempText);
			// #则添加把原始提取的添加进来
			originCHNumberForOutput.push(originCHNumberTake[i]);
		}
		// CHNumberStringList = append(CHNumberStringList, regMatchResult[i][0])
	}

	// """
	// 进行标准汉字字符串转换 例如 二千二  转换成二千零二
	// """
	chNumberStringListTemp = vec![];
	for i in 0..chNumberStringList.len(){
		chNumberStringListTemp.push(standardChNumberConvert(chNumberStringList[i]));

	}

	//"""
	//将中文转换为数字
	//"""
	let mut digitsStringList:Vec<String> =  vec![];
	let mut replacedText = convertedString;
	let mut tempCHToDigitsResult = "".to_string();
	// structCHAndDigitSlice := []structCHAndDigit{}
	if chNumberStringListTemp.len()> 0 {
		for i in 0..chNumberStringListTemp.len(){
			tempCHToDigitsResult = ChineseToDigits(chNumberStringListTemp[i], percentConvert);
			digitsStringList.push(tempCHToDigitsResult);

		}
		//fmt.Println(structCHAndDigitSlice)
		// 按照 中文数字字符串长度 的逆序排序
		digitsStringList.sort_by_key(|x| x.len());
		digitsStringList.reverse();
		//"""
		//按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
		//"""
		for i in 0..digitsStringList.len(){
			replacedText = replacedText.replace(originCHNumberForOutput[i],digitsStringList[i]);
			//fmt.Println(replacedText)
		}

	}

	let finalResult = FinalResultStruct {
		inputText:chTextString,
		replacedText:replacedText,
		chNumberStringList:originCHNumberForOutput,
		digitsStringList:digitsStringList
	};
	return finalResult

}

// TakeNumberFromString will extract the chinese and digits number together from string. and return the convert result
// :param chText: chinese string
// :param percentConvert: convert percent simple. Default is True.  3% will be 0.03 in the result
// :param traditionalConvert: Switch to convert the Traditional Chinese character to Simplified chinese
// :return: Dict like result. 'inputText',replacedText','CHNumberStringList':CHNumberStringList,'digitsStringList'
fn TakeNumberFromString(chTextString:String) -> FinalResultStruct{

	//默认参数设置
	// if len(opt) > 2 {
	// 	panic("too many arguments")
	// }

	// var percentConvert bool
	// var traditionalConvert bool
	// // digitsNumberSwitch := false

	// switch len(opt) {
	// case 1:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = true
	// 	// digitsNumberSwitch = false
	// case 2:
	// 	percentConvert = opt[0].(bool)
	// 	traditionalConvert = opt[1].(bool)
	// 	// digitsNumberSwitch = false
	// // case 3:
	// // 	percentConvert = opt[0].(bool)
	// // 	traditionalConvert = opt[1].(bool)
	// // digitsNumberSwitch = opt[2].(bool)
	// default:
	// 	percentConvert = true
	// 	traditionalConvert = true
	// }
	let percentConvert = true;
	let traditionalConvert = true;
	let finalResult:FinalResultStruct = TakeChineseNumberFromString(chTextString, percentConvert, traditionalConvert);
	return finalResult
}


fn main() {
    println!("{}",core_ch_to_digits("三百四十二万".to_string()));
}

#[test]
fn test_core_ch_to_digits() {
    // do test work
	assert_eq!(core_ch_to_digits("三百四十二万".to_string()), "3420000")
}
#[test]
fn test_check_number_seg() {
    
	let a1 = vec!["百".to_string(),"分之5".to_string(),"负千".to_string(),"分之15".to_string()];
	let a2 = "百分之5负千分之15".to_string();
	let a3 = vec!["百分之5".to_string(),"负千分之15".to_string()];
	assert_eq!(check_number_seg(a1, a2), a3)
}

#[test]
fn test_check_chinese_number_reasonable() {
    
	let a2 = "千千万万".to_string();
	assert_eq!(check_chinese_number_reasonable(a2), false)
}