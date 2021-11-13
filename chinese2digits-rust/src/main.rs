// #[warn(unused_assignments)]
use std::collections::HashMap;
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
		let val_from_hash_map = chinese_char_number_dict.get(&char_to_get[..]);
		
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
						let pre_val_temp_option = chinese_char_number_dict.get(&chinese_chars.chars().nth(i-1).unwrap().to_string()[..]);
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



fn main() {
    println!("{}",core_ch_to_digits("三百四十二".to_string()));
}