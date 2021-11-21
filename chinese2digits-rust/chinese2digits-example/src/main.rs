// extern crate chese2digits;
// include!["../../chinese2digits/src/lib.rs"];
use chinese2digits;

fn main(){
	println!("{}", chinese2digits::take_number_from_string(&"三零万二零千拉阿拉啦啦30万20千嚯嚯或百四嚯嚯嚯四百三十二分之2345啦啦啦啦".to_string()).replaced_text);
}