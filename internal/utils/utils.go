package utils

import (
	"encoding/json"
	"fmt"
	"github.com/LucienVen/golang-auth-service/internal/app"
	"math"
	"os"

	"github.com/LucienVen/golang-auth-service/pkg/log"
)

func FirstInit() {
	// 创建应用实例
	application := app.NewApplication()

	// 启动应用
	if err := application.Start(); err != nil {
		log.Errorf("应用启动失败: %v", err)
		os.Exit(1)
	}

	//log2.Printf("%+v", app.Mysql)
	//log2.Printf("%+v", bootstrap.App.Mysql)
	//log2.Printf("%+v", bootstrap.App.GetDB())

	log.InitLogger()

}

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
