package main

import (
	"fmt"
	"./chinese2digits"
)


func main(){
	fmt.Println(chinese2digits.TakeChineseNumberFromString("负百分之点二八百分之三五",true,true))
}