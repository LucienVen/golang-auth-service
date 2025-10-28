package utils

import (
	"encoding/json"
	"fmt"
	"math"
)

// 注意：FirstInit 函数已移至 cmd/main.go，因为 utils 包不应该依赖应用层逻辑
// 这样可以避免循环依赖问题

func StructPrintf(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// 浮点数保留 N 位小数(返回 float) 不四舍五入
func FormatFloat2Float(num float64, decimal int) float64 {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		d = math.Pow10(decimal)
	}

	return math.Trunc(num*d) / d
}
