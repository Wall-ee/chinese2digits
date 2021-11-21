// extern crate chese2digits;
// include!["../../chinese2digits/src/lib.rs"];
use chinese2digits;

fn main(){
	println!("{}", chinese2digits::take_number_from_string(&"负百分之点二八你好啊百分之三五是不是点五零百分之负六十五点二八",true,true).replaced_text);
}