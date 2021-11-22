// #[warn(unused_assignments)]
// #[macro_use] extern crate lazy_static;
// extern crate regex;
use regex::Regex;
//pcre2 速度快
// extern crate pcre2;
// use pcre2::bytes::Regex;
use std::collections::HashMap;
use std::str;
use rust_decimal::prelude::*;

// 中文转阿拉伯数字
static CHINESE_CHAR_NUMBER_LIST: [(&str, i32); 29] = [
	("幺", 1),
	("零", 0),
	("一", 1),
	("二", 2),
	("两", 2),
	("三", 3),
	("四", 4),
	("五", 5),
	("六", 6),
	("七", 7),
	("八", 8),
	("九", 9),
	("十", 10),
	("百", 100),
	("千", 1000),
	("万", 10000),
	("亿", 100000000),
	("壹", 1),
	("贰", 2),
	("叁", 3),
	("肆", 4),
	("伍", 5),
	("陆", 6),
	("柒", 7),
	("捌", 8),
	("玖", 9),
	("拾", 10),
	("佰", 100),
	("仟", 1000),
];

static CHINESE_PURE_COUNTING_UNIT_LIST: [&str; 5] = ["十", "百", "千", "万", "亿"];

static CHINESE_PURE_NUMBER_LIST: [&str; 13] = [
	"幺", "一", "二", "两", "三", "四", "五", "六", "七", "八", "九", "十", "零",
];

// CoreCHToDigits 是核心转化函数
fn core_ch_to_digits(chinese_chars_to_trans: String) -> String {
	let chinese_chars = chinese_chars_to_trans;
	let total: String;
	let mut temp_val = "".to_string(); //#用以记录临时是否建议数字拼接的字符串 例如 三零万 的三零
	let mut counting_unit: i32 = 1; //#表示单位：个十百千,用以计算单位相乘 例如八百万 百万是相乘的方法，但是如果万前面有 了一千八百万 这种，千和百不能相乘，要相加...
	let mut counting_unit_from_string: Vec<i32> = vec![1]; //#原始字符串提取的单位应该是一个list  在计算的时候，新的单位应该是本次取得的数字乘以已经发现的最大单位，例如 4千三百五十万， 等于 4000万+300万+50万
	let mut temp_val_int: i32;
	let chinese_char_number_dict: HashMap<&str, i32> =
		CHINESE_CHAR_NUMBER_LIST.into_iter().collect();

	let mut temp_total: i32 = 0;
	// // counting_unit := 1
	// 表示单位：个十百千...
	for i in (0..chinese_chars.chars().count()).rev() {
		// println!("{}",i);
		let char_to_get = chinese_chars.chars().nth(i).unwrap().to_string();
		let val_from_hash_map = chinese_char_number_dict.get(char_to_get.as_str());

		match val_from_hash_map {
			Some(val_from_hash_map) => {
				let val = *val_from_hash_map;
				if val >= 10 && i == 0_usize {
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
				} else if val >= 10 {
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
				} else {
					if i > 0 {
						// #如果下一个不是单位 则本次也是拼接
						let pre_val_temp_option = chinese_char_number_dict.get(
							chinese_chars
								.chars()
								.nth(i - 1)
								.unwrap()
								.to_string()
								.as_str(),
						);
						match pre_val_temp_option {
							Some(pre_val_temp_option) => {
								let pre_val_temp = *pre_val_temp_option;
								if pre_val_temp < 10 {
									// temp_val = strconv.Itoa(val) + temp_val;
									temp_val = val.to_string() + &temp_val;
								} else {
									// #说明已经有大于10的单位插入 要数学计算了
									// #先拼接再计算
									// #如果取值不大于10 说明是0-9 则继续取值 直到取到最近一个大于10 的单位   应对这种30万20千 这样子
									temp_val_int =
										(val.to_string() + &temp_val).parse::<i32>().unwrap();
									temp_total = temp_total + counting_unit * temp_val_int;
									// #计算后 把临时字符串置位空
									temp_val = "".to_string();
								}
							}
							None => {}
						}
					} else {
						// #那就是无论如何要收尾了
						//如果counting unit 等于1  说明所有字符串都是直接拼接的，不用计算，不然会丢失前半部分的零
						if counting_unit == 1 {
							temp_val = val.to_string() + &temp_val;
						} else {
							temp_val_int = (val.to_string() + &temp_val).parse::<i32>().unwrap();
							temp_total = temp_total + counting_unit * temp_val_int;
						}
					}
				}
			}
			None => println!("No string in number"),
		}
	}
	//如果 total 为0  但是 counting_unit 不为0  说明结果是 十万这种  最终直接取结果 十万
	if temp_total == 0 {
		if counting_unit > 10 {
			// 转化为字符串
			total = counting_unit.to_string();
		} else {
			//counting Unit 为1  且tempval 不为空 说明是没有单位的纯数字拼接
			if temp_val != "" {
				total = temp_val;
			} else {
				// 转化为字符串
				total = temp_total.to_string();
			}
		}
	} else {
		// 转化为字符串
		total = temp_total.to_string();
	}

	return total;
}

//汉字切分是否正确  分之 切分
fn check_number_seg(chinese_number_list: &Vec<String>, origin_text: &String) -> Vec<String> {
	let mut new_chinese_number_list: Vec<String> = [].to_vec();
	// #用来控制是否前一个已经合并过  防止多重合并
	let mut temp_pre_text: String = "".to_string();
	let mut temp_mixed_string: String;
	let seg_len = chinese_number_list.len();
	if seg_len > 0 {
		// #加入唯一的一个 或者第一个
		if chinese_number_list[0].starts_with("分之") {
			// #如果以分之开头 记录本次 防止后面要用 是否出现连续的 分之
			new_chinese_number_list.push(
				(&chinese_number_list[0][2..(chinese_number_list[0].chars().count())]).to_string(),
			);
		} else {
			new_chinese_number_list.push(chinese_number_list.get(0).unwrap().to_string())
		}

		if seg_len > 1 {
			for i in 1..seg_len {
				// #判断本字符是不是以  分之  开头
				if (chinese_number_list[i]).starts_with("分之") {
					// #如果是以 分之 开头 那么检查他和他见面的汉子数字是不是连续的 即 是否在原始字符串出现
					temp_mixed_string = chinese_number_list.get(i - 1).unwrap().to_string()
						+ &chinese_number_list[i];

					if origin_text.contains(&temp_mixed_string.to_string()) {
						// #如果连续的上一个字段是以分之开头的  本字段又以分之开头
						if temp_pre_text != "" {
							// #检查上一个字段的末尾是不是 以 百 十 万 的单位结尾
							if CHINESE_PURE_COUNTING_UNIT_LIST
								.iter()
								.any(|&x| x == (&temp_pre_text.chars().last().unwrap().to_string()))
							{
								//上一个字段最后一个数据
								let temp_last_chinese_number =
									new_chinese_number_list.last().unwrap().to_string();
								// #先把上一个记录进去的最后一位去掉
								let temp_new_chinese_number_list_len =
									new_chinese_number_list.len();
								let temp_last_chinese_number_char_list: Vec<char> =
									temp_last_chinese_number.chars().collect();
								//把暂存结果的上一个 最后一个单位去掉
								new_chinese_number_list[temp_new_chinese_number_list_len - 1] =
									temp_last_chinese_number_char_list
										[0..(temp_last_chinese_number_char_list.len() - 1)]
										.iter()
										.collect();
								// #如果结果是确定的，那么本次的字段应当加上上一个字段的最后一个字
								new_chinese_number_list.push(
									temp_pre_text
										.chars()
										.nth(temp_pre_text.chars().count() - 1)
										.unwrap()
										.to_string() + &chinese_number_list[i],
								);
							} else {
								// #如果上一个字段不是以单位结尾  同时他又是以分之开头，那么 本次把分之去掉
								new_chinese_number_list
									.push((&chinese_number_list[i][2..]).to_string())
							}
						} else {
							// #上一个字段不以分之开头，那么把两个字段合并记录
							if new_chinese_number_list.len() > 0 {
								//把最后一位赋值
								// newChineseNumberList[len(newChineseNumberList)-1] = tempMixedString
								let temp_new_chinese_number_list_len =
									new_chinese_number_list.len();
								new_chinese_number_list[temp_new_chinese_number_list_len - 1] =
									temp_mixed_string;
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
	return new_chinese_number_list;
}

fn check_chinese_number_reasonable(ch_number: &String) -> bool {
	if ch_number.chars().count() > 0 {
		// #由于在上个检查点 已经把阿拉伯数字转为中文 因此不用检查阿拉伯数字部分
		// """
		// 如果汉字长度大于0 则判断是不是 万  千  单字这种
		// """
		for i in CHINESE_PURE_NUMBER_LIST {
			if ch_number.contains(&i.to_string()) {
				return true;
			}
		}
	}
	return false;
}

// """
// 阿拉伯数字转中文
// """
static DIGITS_CHAR_CHINESE_DICT_KEY: [&str; 14] = [
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "%", "‰", "‱", ".",
];
static DIGITS_CHAR_CHINESE_DICT: [(&str, &str); 14] = [
	("0", "零"),
	("1", "一"),
	("2", "二"),
	("3", "三"),
	("4", "四"),
	("5", "五"),
	("6", "六"),
	("7", "七"),
	("8", "八"),
	("9", "九"),
	("%", "百分之"),
	("‰", "千分之"),
	("‱", "万分之"),
	(".", "点"),
];

static CHINESE_PER_COUNTING_STRING_LIST: [&str; 3] = ["百分之", "千分之", "万分之"];

static TRADITIONAL_CONVERT_DICT_STATIC: [(&str, &str); 9] = [
	("壹", "一"),
	("贰", "二"),
	("叁", "三"),
	("肆", "四"),
	("伍", "五"),
	("陆", "六"),
	("柒", "七"),
	("捌", "八"),
	("玖", "九"),
];
static SPECIAL_TRADITIONAL_COUNTING_UNIT_CHAR_DICT_STATIC: [(&str, &str); 5] = [
	("拾", "十"),
	("佰", "百"),
	("仟", "千"),
	("萬", "万"),
	("億", "亿"),
];
static SPECIAL_NUMBER_CHAR_DICT: [(&str, &str); 2] = [("两", "二"), ("俩", "二")];
// static CHINESE_PURE_NUMBER_LIST: [&str; 13] = ["幺", "一", "二", "两", "三", "四", "五", "六", "七", "八", "九", "十", "零"];
static CHINESE_SIGN_LIST: [&str; 4] = ["正", "负", "+", "-"];

//阿拉伯数字转中文
fn digits_to_ch_chars(mixed_string_list: &Vec<String>) -> Vec<String> {
	let mut result_list: Vec<String> = vec![];
	let digits_char_chinese_dict: HashMap<&str, &str> =
		DIGITS_CHAR_CHINESE_DICT.into_iter().collect();
	// for i := 0; i < len(mixedStringList); i++ {
	for i in mixed_string_list.iter() {
		let mut mixed_string = i.to_string();
		if mixed_string.starts_with(".") {
			mixed_string = "0".to_string() + &mixed_string;
		}
		for key in DIGITS_CHAR_CHINESE_DICT_KEY.iter() {
			if mixed_string.contains(&key.to_string()) {
				mixed_string = mixed_string.replace(
					key,
					digits_char_chinese_dict
						.get(key)
						.unwrap()
						.to_string()
						.as_str(),
				);
				// #应当是只要有百分号 就挪到前面 阿拉伯数字没有四百分之的说法
				// #防止这种 3%万 这种问题
				for kk in CHINESE_PER_COUNTING_STRING_LIST.iter() {
					if mixed_string.contains(&kk.to_string()) {
						mixed_string = kk.to_string() + &mixed_string.replace(kk, "");
					}
				}
			}
		}
		result_list.push(mixed_string);
	}

	return result_list;
}

// """
// 繁体简体转换 及  单位  特殊字符转换 两千变二千
// """
fn traditional_text_convert_func(ch_string: &String, simplif_convert_switch: bool) -> String {
	// chStringList := []rune(chString)
	let mut ch_string_list: Vec<char> = ch_string.chars().collect();

	let string_length = ch_string.chars().count();
	let mut char_to_get: String;
	let traditional_convert_dict: HashMap<&str, &str> =
		TRADITIONAL_CONVERT_DICT_STATIC.into_iter().collect();
	let special_traditional_counting_unit_char_dict: HashMap<&str, &str> =
		SPECIAL_TRADITIONAL_COUNTING_UNIT_CHAR_DICT_STATIC
			.into_iter()
			.collect();
	let special_number_char_dict: HashMap<&str, &str> =
		SPECIAL_NUMBER_CHAR_DICT.into_iter().collect();
	if simplif_convert_switch {
		for i in 0..ch_string_list.len() {
			// #繁体中文数字转简体中文数字
			// charToGet = string(chStringList[i])
			char_to_get = ch_string_list[i].to_string();
			let value = traditional_convert_dict.get(char_to_get.as_str());
			match value {
				Some(value) => {
					ch_string_list[i] = (*value).chars().nth(0).unwrap();
				}
				None => {}
			}
		}
	}
	if string_length > 1 {
		// #检查繁体单体转换
		for i in 0..string_length {
			// #如果 前后有 pure 汉字数字 则转换单位为简体
			char_to_get = ch_string_list[i].to_string();
			// value, exists := SPECIAL_TRADITIONAl_COUNTING_UNIT_CHAR_DICT[charToGet]
			let value = special_traditional_counting_unit_char_dict.get(char_to_get.as_str());
			// # 如果前后有单纯的数字 则进行单位转换
			match value {
				Some(value) => {
					if i == 0 {
						if CHINESE_PURE_NUMBER_LIST
							.iter()
							.any(|&x| x == ch_string_list[i + 1].to_string())
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					} else if i == (string_length - 1) {
						if CHINESE_PURE_NUMBER_LIST
							.iter()
							.any(|&x| x == ch_string_list[i - 1].to_string())
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					} else {
						if (CHINESE_PURE_NUMBER_LIST
							.iter()
							.any(|&x| x == ch_string_list[i - 1].to_string()))
							|| (CHINESE_PURE_NUMBER_LIST
								.iter()
								.any(|&x| x == ch_string_list[i + 1].to_string()))
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					}
				}
				None => {}
			}
			// #特殊变换 俩变二
			char_to_get = ch_string_list[i].to_string();
			let value = special_number_char_dict.get(char_to_get.as_str());
			// # 如果前后有单纯的数字 则进行单位转换
			match value {
				Some(value) => {
					if i == 0 {
						if CHINESE_PURE_COUNTING_UNIT_LIST
							.iter()
							.any(|&x| x == ch_string_list[i + 1].to_string())
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					} else if i == (string_length - 1) {
						if CHINESE_PURE_COUNTING_UNIT_LIST
							.iter()
							.any(|&x| x == ch_string_list[i - 1].to_string())
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					} else {
						if (CHINESE_PURE_COUNTING_UNIT_LIST
							.iter()
							.any(|&x| x == ch_string_list[i - 1].to_string()))
							|| (CHINESE_PURE_COUNTING_UNIT_LIST
								.iter()
								.any(|&x| x == ch_string_list[i + 1].to_string()))
						{
							ch_string_list[i] = (*value).chars().nth(0).unwrap();
						}
					}
				}
				None => {}
			}
		}
	}
	let final_ch_string_list: String = ch_string_list.iter().collect();

	return final_ch_string_list;
}

// """
// 标准表述转换  三千二 变成 三千零二  三千十二变成 三千零一十二
// """
fn standard_ch_number_convert(ch_number_string: String) -> String {
	let ch_number_string_list: Vec<char> = ch_number_string.chars().collect();

	let mut new_ch_number_string_list: String = ch_number_string;

	// #大于2的长度字符串才有检测和补位的必要
	if ch_number_string_list.len() > 2 {
		// #十位补一：
		let ten_number_index = ch_number_string_list
			.iter()
			.position(|&x| x == "十".chars().nth(0).unwrap());
		match ten_number_index {
			Some(ten_number_index) => {
				if ten_number_index == 0_usize {
					new_ch_number_string_list = "一".to_string() + &new_ch_number_string_list;
				} else {
					// # 如果没有左边计数数字 插入1
					if !(CHINESE_PURE_NUMBER_LIST.iter().any(|&x| {
						x == ch_number_string_list[(ten_number_index - 1)]
							.to_string()
							.as_str()
					})) {
						let temp_left_part: String =
							ch_number_string_list[0..ten_number_index].iter().collect();
						let temp_right_part: String =
							ch_number_string_list[ten_number_index..].iter().collect();
						new_ch_number_string_list = temp_left_part + "一" + &temp_right_part;
					}
				}
			}
			None => {}
		}

		// #差位补零
		// #逻辑 如果最后一个单位 不是十结尾 而是百以上 则数字后面补一个比最后一个出现的单位小一级的单位
		// #从倒数第二位开始看,且必须是倒数第二位就是单位的才符合条件

		let temp_new_ch_number_string_list: Vec<char> = new_ch_number_string_list.chars().collect();

		let last_counting_unit = CHINESE_PURE_COUNTING_UNIT_LIST.iter().position(|&x| {
			x == temp_new_ch_number_string_list[temp_new_ch_number_string_list.len() - 2]
				.to_string()
				.as_str()
		});
		// # 如果最末位的是百开头
		match last_counting_unit {
			Some(last_counting_unit) => {
				if last_counting_unit >= 1_usize {
					// # 则字符串最后拼接一个比最后一个单位小一位的单位 例如四万三 变成四万三千
					// # 如果最后一位结束的是亿 则补千万
					if last_counting_unit == 4_usize {
						new_ch_number_string_list = new_ch_number_string_list + &"千万".to_string();
					} else {
						new_ch_number_string_list = new_ch_number_string_list
							+ &CHINESE_PURE_COUNTING_UNIT_LIST[last_counting_unit - 1].to_string();
					}
				}
			}
			None => {}
		}
	}
	//大于一的检查是不是万三，千四五这种
	let mut per_count_switch = false;
	let temp_new_ch_number_string_list: Vec<char> = new_ch_number_string_list.chars().collect();
	if temp_new_ch_number_string_list.len() > 1 {
		// #十位补一：
		// fistCharCheckResult := isExistItem(string(tempNewChNumberStringList[0]), []string{"千", "万", "百"})

		let _first_char_check_result = ["千", "万", "百"]
			.iter()
			.position(|&x| x == temp_new_ch_number_string_list[0].to_string().as_str());
		match _first_char_check_result {
			Some(_first_char_check_result) => {
				for i in 1..temp_new_ch_number_string_list.len() {
					// #其余位数都是纯数字 才能执行
					if CHINESE_PURE_NUMBER_LIST
						.iter()
						.any(|&x| x == temp_new_ch_number_string_list[i].to_string().as_str())
					{
						per_count_switch = true;
					} else {
						per_count_switch = false;
						//有一个不合适 退出循环
						break;
					}
				}
				if per_count_switch {
					let temp_left_part_string: String =
						temp_new_ch_number_string_list[0..1].iter().collect();
					let temp_right_part_string: String =
						temp_new_ch_number_string_list[1..].iter().collect();
					new_ch_number_string_list =
						temp_left_part_string + "分之" + &temp_right_part_string;
				}
			}
			None => {}
		}
	}

	return new_ch_number_string_list;
}

//检查初次提取的汉字数字是正负号是否切分正确
fn check_sign_seg(chinese_number_list: &Vec<String>) -> Vec<String> {
	let mut new_chinese_number_list: Vec<String> = vec![];
	let mut temp_sign = "".to_string();
	for i in 0..chinese_number_list.len() {
		// #新字符串 需要加上上一个字符串 最后1位的判断结果
		let mut new_ch_number_string = temp_sign.clone() + &chinese_number_list[i];
		let temp_chinese_number_list: Vec<char> = new_ch_number_string.chars().collect();
		if temp_chinese_number_list.len() > 1 {
			let last_string: String = temp_chinese_number_list
				[(temp_chinese_number_list.len() - 1)..]
				.iter()
				.collect();
			// #如果最后1位是百分比 那么本字符去掉最后三位  下一个数字加上最后1位
			if CHINESE_SIGN_LIST.iter().any(|&x| x == last_string.as_str()) {
				temp_sign = last_string;
				// #如果最后1位 是  那么截掉最后1位
				// let tempWithoutLastPartString:String =
				new_ch_number_string = temp_chinese_number_list
					[..(temp_chinese_number_list.len() - 1)]
					.iter()
					.collect();
			} else {
				temp_sign = "".to_string();
			}
		}
		new_chinese_number_list.push(new_ch_number_string)
	}
	return new_chinese_number_list;
}

pub struct C2DResultStruct {
	pub input_text: String,
	pub replaced_text: String,
	pub ch_number_string_list: Vec<String>,
	pub digits_string_list: Vec<String>,
}

// static TAKING_CHINESE_DIGITS_MIX_RE_RULES_STRING: &str = r#"(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})(?:(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}
// |\.\d+(?:[\%]){0,1}){0,1}(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))))
// |(?:(?:\d+(?:\.\d+){0,1}(?:[\%]){0,1}|\.\d+(?:[\%]){0,1})(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+
// 	(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))"#;
//rust  % 号不能加转义
static TAKING_CHINESE_DIGITS_MIX_RE_RULES_STRING: &str = concat!(
	r#"(?:(?:分之){0,1}(?:\+|\-){0,1}[正负]{0,1})(?:(?:(?:\d+(?:\.\d+){0,1}(?:[%‰‱]){0,1}"#,
	r#"|\.\d+(?:[%‰‱]){0,1}){0,1}(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+(?:点[一二三四五六七八九万亿兆幺零]+){0,1})"#,
	r#"|(?:点[一二三四五六七八九万亿兆幺零]+))))|(?:(?:\d+(?:\.\d+){0,1}(?:[%‰‱]){0,1}|\.\d+(?:[%‰‱]){0,1})"#,
	r#"(?:(?:(?:[一二三四五六七八九十千万亿兆幺零百]+"#,
	r#"(?:点[一二三四五六七八九万亿兆幺零]+){0,1})|(?:点[一二三四五六七八九万亿兆幺零]+))){0,1}))"#
);

static CHINESE_SIGN_DICT_STATIC: [(&str, &str); 4] =
	[("负", "-"), ("正", "+"), ("-", "-"), ("+", "+")];

static CHINESE_CONNECTING_SIGN_STATIC: [(&str, &str); 3] = [(".", "."), ("点", "."), ("·", ".")];

// ChineseToDigits 是可以识别包含百分号，正负号的函数，并控制是否将百分之10转化为0.1
fn chinese_to_digits(chinese_chars_to_trans: String, percent_convert: bool) -> String {
	// """
	// 分之  分号切割  要注意
	// """
	let mut final_total: String;
	let mut convert_result_list: Vec<String> = vec![];
	let mut chinese_chars_list_by_div: Vec<String> = vec![];

	let chinese_sign_dict: HashMap<&str, &str> = CHINESE_SIGN_DICT_STATIC.into_iter().collect();
	let chinese_connecting_sign_dict: HashMap<&str, &str> =
		CHINESE_CONNECTING_SIGN_STATIC.into_iter().collect();

	if chinese_chars_to_trans.contains("分之") {
		let temp_split_result: Vec<&str> = chinese_chars_to_trans.split("分之").collect();
		for s in temp_split_result {
			chinese_chars_list_by_div.push(s.to_string());
		}
	} else {
		chinese_chars_list_by_div.push(chinese_chars_to_trans);
	}

	for k in 0..chinese_chars_list_by_div.len() {
		let mut temp_chinese_char: String = chinese_chars_list_by_div[k].to_string();

		// chinese_char := temp_chinese_char
		let chinese_char: Vec<char> = temp_chinese_char.chars().collect();
		// """
		// 看有没有符号
		// """
		let mut sign = "".to_string();
		let mut char_to_get: String;
		for i in 0..chinese_char.len() {
			char_to_get = chinese_char[i].to_string();
			let value = chinese_sign_dict.get(char_to_get.as_str());
			match value {
				Some(value) => {
					sign = value.to_string();
					// chineseCharsToTrans = strings.Replace(chineseCharsToTrans, charToGet, "", -1)
					temp_chinese_char = temp_chinese_char.replace(&char_to_get, "");
				}
				None => {}
			}
		}
		// chinese_char = []rune(chineseCharsToTrans)

		// """
		// 小数点切割，看看是不是有小数点
		// """
		let mut string_contain_dot = false;
		let mut left_of_dot_string = "".to_string();
		let mut right_of_dot_string = "".to_string();
		for (key, _) in &chinese_connecting_sign_dict {
			if temp_chinese_char.contains(&key.to_string()) {
				let chinese_chars_dot_split_list: Vec<&str> =
					temp_chinese_char.split(key).collect();
				left_of_dot_string = chinese_chars_dot_split_list[0].to_string();
				right_of_dot_string = chinese_chars_dot_split_list[1].to_string();
				string_contain_dot = true;
				break;
			}
		}
		let mut convert_result: String;
		if !string_contain_dot {
			convert_result = core_ch_to_digits(temp_chinese_char);
		} else {
			// convert_result = "".to_string();
			let mut temp_buf: String;
			let temp_right_digits: String;
			// #如果小数点右侧有 单位 比如 2.55万  4.3百万 的处理方式
			// #先把小数点右侧单位去掉
			let mut temp_count_string = "".to_string();
			let list_of_right: Vec<char> = right_of_dot_string.chars().collect();
			for ii in (0..(list_of_right.len())).rev() {
				let _find_counting_unit_index_result = CHINESE_PURE_COUNTING_UNIT_LIST
					.iter()
					.position(|&x| x == list_of_right[ii].to_string().as_str());
				match _find_counting_unit_index_result {
					Some(_find_counting_unit_index_result) => {
						temp_count_string = list_of_right[ii].to_string() + &temp_count_string;
					}
					None => {
						right_of_dot_string = list_of_right[0..(ii + 1)].iter().collect();
						break;
					}
				}
			}

			let mut temp_count_num = 1.0;
			if temp_count_string != "" {
				let temp_num = core_ch_to_digits(temp_count_string).parse::<f32>().unwrap();
				temp_count_num = temp_num;
			}

			if left_of_dot_string == "" {
				// """
				// .01234 这种开头  用0 补位
				// """
				temp_buf = "0.".to_string();
				temp_right_digits = core_ch_to_digits(right_of_dot_string);
				temp_buf = temp_buf + &temp_right_digits;
				convert_result = temp_buf;
			} else {
				temp_buf = core_ch_to_digits(left_of_dot_string);
				temp_buf = temp_buf + ".";
				temp_right_digits = core_ch_to_digits(right_of_dot_string);
				temp_buf = temp_buf + &temp_right_digits;
				convert_result = temp_buf;
			}

			let temp_str_to_float = convert_result.parse::<f32>().unwrap();
			convert_result = (temp_str_to_float * temp_count_num).to_string();
		}
		//如果转换结果为空字符串 则为百分之10 这种
		if convert_result == "" {
			convert_result = "1".to_string();
		}
		convert_result = sign + &convert_result;
		// #最后在双向转换一下 防止出现 0.3000 或者 00.300的情况

		let new_convert_result_temp: Vec<char> = convert_result.chars().collect();
		let new_buf: String;
		if convert_result.ends_with(".0") {
			new_buf = new_convert_result_temp[0..new_convert_result_temp.len() - 2]
				.iter()
				.collect();
		} else {
			new_buf = convert_result;
		}
		convert_result_list.push(new_buf);
	}
	if convert_result_list.len() > 1 {
		// #是否转换分号及百分比
		if percent_convert {
			// let temp_float1 = convert_result_list[1].parse::<f32>().unwrap();
			// let temp_float0 = convert_result_list[0].parse::<f32>().unwrap();
			// // fmt.Println(tempFloat1 / tempFloat0)
			// final_total = (temp_float1 / temp_float0).to_string();
			let temp_float1 = Decimal::from_str(&convert_result_list[1]).unwrap();
			let temp_float0 = Decimal::from_str(&convert_result_list[0]).unwrap();
			final_total = (temp_float1 / temp_float0).to_string();
		} else {
			if convert_result_list[0] == "100" {
				final_total = convert_result_list[1].to_string() + "%"
			} else if convert_result_list[0] == "1000" {
				final_total = convert_result_list[1].to_string() + "‰"
			} else {
				final_total = convert_result_list[1].to_string() + "/" + &convert_result_list[0]
			}
		}
	} else {
		final_total = convert_result_list[0].to_string()
		//最后再转换一下 防止出现 .50 的问题  不能转换了 否则  超出精度了………… 服了  5亿的话
		// tempFinalTotal, err3 := strconv.ParseFloat(finalTotal, 32)
		// if err3 != nil {
		// 	panic(err3)
		// } else {
		// 	finalTotal = strconv.FormatFloat(tempFinalTotal, 'f', -1, 32)
		// }
	}
	//删除最后的0
	if final_total.contains("."){
		final_total = final_total.trim_end_matches("0").to_string();
		final_total = final_total.trim_end_matches(".").to_string();
	}
	

	return final_total;
}

// 将句子中的汉子数字提取的整体函数
/// Will extract the chinese and digits number together from string. and return the convert result
/// 
/// Returns `C2DResultStruct`  data struct
/// ```
/// pub struct C2DResultStruct {
///		pub input_text: String,
///		pub replaced_text: String,
///		pub ch_number_string_list: Vec<String>,
///		pub digits_string_list: Vec<String>,
/// }
/// ```
///
///
/// [`chText`]: chinese string
/// 
/// [`percentConvert`]: convert percent simple. Default is True.  3% will be 0.03 in the result
/// 
/// [`traditionalConvert`]: Switch to convert the Traditional Chinese character to Simplified chinese
/// 
/// # Examples
///
/// Basic usage:
///
/// ```
/// use chinese2digits::take_number_from_string;
/// 
/// let string_example = "百分之5负千分之15".to_string();
/// let test_result = take_number_from_string(&string_example, true, true);
/// assert_eq!(test_result.replaced_text, "0.05-0.015");
///	assert_eq!(test_result.ch_number_string_list, vec!["百分之5", "负千分之15"]);
///	assert_eq!(test_result.digits_string_list, vec!["0.05", "-0.015"]);
/// ```
/// 
fn take_chinese_number_from_string(
	ch_text_string: &str,
	percent_convert: bool,
	traditional_convert: bool,
) -> C2DResultStruct {
	let mut ch_number_string_list: Vec<String> = vec![];

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
	let converted_string: String =
		traditional_text_convert_func(&ch_text_string.to_string(), traditional_convert);

	//正则引擎
	let mut reg_match_result: Vec<String> = vec![];
	let taking_chinese_digits_mix_re_rules: Regex =
		Regex::new(TAKING_CHINESE_DIGITS_MIX_RE_RULES_STRING).unwrap();
	//PCRE2 引擎，但是正则断句不对导致字符串切割总是失败。
	// for cap in taking_chinese_digits_mix_re_rules.captures_iter(&converted_string.as_bytes()) {
	// 	let caps = cap.unwrap();
	// 	let cap_result = str::from_utf8(&caps[0]);
	// 	match cap_result{
	// 		Ok(cap_result) => {
	// 			reg_match_result.push(cap_result.to_string());
	// 		}
	// 		Err(cap_result) =>{
	// 			println!("{}",cap_result)
	// 		}

	// 	}

	// }

	// let caps=taking_chinese_digits_mix_re_rules.captures(&converted_string.as_bytes()).unwrap();

	// match caps{
	// 	Some(caps)=>{
	// 		for i in 0..caps.len(){
	// 			let cap = caps.get(i).unwrap();
	// 			let cap_result = str::from_utf8(cap.as_bytes());
	// 			match cap_result{
	// 				Ok(cap_result)=>{
	// 					println!("{}",cap_result.to_string());
	// 					reg_match_result.push(cap_result.to_string());
	// 				}
	// 				Err(cap_result) =>{
	// 					println!("error is {}",cap_result);
	// 				}
	// 			}

	// 		}

	// 	}
	// 	None=>{

	// 	}
	// }

	// println!("Values inside vec: {:?}", reg_match_result);
	for cap in taking_chinese_digits_mix_re_rules.captures_iter(&converted_string) {
		reg_match_result.push(cap[0].to_string());
	}

	let mut temp_text: String;
	let mut ch_number_string_list_temp: Vec<String> = vec![];
	for i in 0..reg_match_result.len() {
		temp_text = reg_match_result[i].to_string();
		ch_number_string_list_temp.push(temp_text);
	}
	// println!("Values inside ch_number_string_list_temp: {:?}", ch_number_string_list_temp);
	// ##检查是不是  分之 切割不完整问题
	ch_number_string_list_temp = check_number_seg(&ch_number_string_list_temp, &converted_string);

	// 检查最后是正负号的问题
	ch_number_string_list_temp = check_sign_seg(&ch_number_string_list_temp);

	// #备份一个原始的提取，后期处结果的时候显示用
	let origin_ch_number_take: Vec<String> = ch_number_string_list_temp.clone();

	// #将阿拉伯数字变成汉字  不然合理性检查 以及后期 如果不是300万这种乘法  而是 四分之345  这种 就出错了
	ch_number_string_list_temp = digits_to_ch_chars(&ch_number_string_list_temp);

	//检查合理性 是否是单纯的单位  等
	// var chNumberStringList []string
	// var originCHNumberForOutput []string
	let mut origin_ch_number_for_output: Vec<String> = vec![];
	for i in 0..ch_number_string_list_temp.len() {
		// fmt.Println(aa[i])
		temp_text = ch_number_string_list_temp[i].to_string();
		if check_chinese_number_reasonable(&temp_text) {
			// #如果合理  则添加进被转换列表
			ch_number_string_list.push(temp_text);
			// #则添加把原始提取的添加进来
			origin_ch_number_for_output.push(origin_ch_number_take[i].to_string());
		}
		// CHNumberStringList = append(CHNumberStringList, regMatchResult[i][0])
	}

	// """
	// 进行标准汉字字符串转换 例如 二千二  转换成二千零二
	// """
	ch_number_string_list_temp = vec![];
	for i in 0..ch_number_string_list.len() {
		ch_number_string_list_temp.push(standard_ch_number_convert(
			ch_number_string_list[i].to_string(),
		));
	}

	//"""
	//将中文转换为数字
	//"""
	let mut digits_string_list: Vec<String> = vec![];
	let mut replaced_text = converted_string;
	let mut temp_ch_to_digits_result: String;
	if ch_number_string_list_temp.len() > 0 {
		for i in 0..ch_number_string_list_temp.len() {
			temp_ch_to_digits_result =
				chinese_to_digits(ch_number_string_list_temp[i].to_string(), percent_convert);
			digits_string_list.push(temp_ch_to_digits_result);
		}
		// 按照 中文数字字符串长度 的逆序排序
		// println!("Values inside origin digits_string_list: {:?}", digits_string_list);

		let mut tuple_to_replace: Vec<_> = origin_ch_number_for_output
			.iter()
			.zip(digits_string_list.iter())
			.enumerate()
			.collect();

		tuple_to_replace.sort_by_key(|(_, (_, z))| z.chars().count());

		// digits_string_list.sort_by_key(|x| x.chars().count());
		//自动逆序 不需要 reverse（）
		// digits_string_list.reverse();

		// println!("Values inside digits_string_list: {:?}", digits_string_list);
		// println!("replaced text is {}",replaced_text);
		// println!("Values inside origin_ch_number_for_output: {:?}", origin_ch_number_for_output);
		//"""
		//按照提取出的中文数字字符串长短排序，然后替换。防止百分之二十八 ，二十八，这样的先把短的替换完了的情况
		//"""
		for i in 0..digits_string_list.len() {
			// replaced_text = replaced_text.replace(
			// 	&origin_ch_number_for_output[i].to_string(),
			// 	&digits_string_list[i].to_string(),
			// );
			replaced_text = replaced_text.replace(
				tuple_to_replace[i].1 .0, //&origin_ch_number_for_output[i]
				tuple_to_replace[i].1 .1, //&digits_string_list[i]
			);
		}
	}

	let final_result = C2DResultStruct {
		input_text: ch_text_string.to_string(),
		replaced_text: replaced_text.to_string(),
		ch_number_string_list: origin_ch_number_for_output,
		digits_string_list: digits_string_list,
	};
	return final_result;
}

pub fn take_number_from_string(
	ch_text_string: &str,
	percent_convert: bool,
	traditional_convert: bool,
) -> C2DResultStruct {
	//默认参数设置

	// let percent_convert = true;
	// let traditional_convert = true;
	let final_result: C2DResultStruct =
		take_chinese_number_from_string(ch_text_string, percent_convert, traditional_convert);
	return final_result;
}

#[test]
fn test_core_ch_to_digits() {
	// do test work
	assert_eq!(core_ch_to_digits("三百四十二万".to_string()), "3420000")
}
#[test]
fn test_check_number_seg() {
	let a1 = vec![
		"百".to_string(),
		"分之5".to_string(),
		"负千".to_string(),
		"分之15".to_string(),
	];
	let a2 = "百分之5负千分之15".to_string();
	let a3 = vec!["百分之5".to_string(), "负千分之15".to_string()];
	assert_eq!(check_number_seg(&a1, &a2), a3)
}

#[test]
fn test_check_chinese_number_reasonable() {
	let a2 = "千千万万".to_string();
	assert_eq!(check_chinese_number_reasonable(&a2), false)
}
