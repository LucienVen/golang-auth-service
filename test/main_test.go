package main

import (
	"fmt"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	// 简单的测试用例
	x := 10
	y := 20
	result := add(x, y)

	// 添加一些可以设置断点的操作
	time.Sleep(100 * time.Millisecond) // 添加延时，方便调试

	if result != 30 {
		t.Errorf("Expected 30, got %d", result)
	}

	// 打印一些信息
	fmt.Printf("测试结果: %d + %d = %d\n", x, y, result)
}

func Example() {
	// 示例函数
	x := 5
	y := 3
	result := add(x, y)
	fmt.Printf("示例结果: %d + %d = %d\n", x, y, result)
	// Output: 示例结果: 5 + 3 = 8
}
