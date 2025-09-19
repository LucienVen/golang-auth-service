package main

import "fmt"

// add 函数用于计算两个整数的和
func add(a, b int) int {
	return a + b
}

func main() {
	x := 10
	y := 20
	result := add(x, y)
	fmt.Printf("计算结果: %d + %d = %d\n", x, y, result)
}
