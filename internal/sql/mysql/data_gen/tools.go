package data_gen

import (
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"time"
)

// 常见单姓 + 复姓
var surnames = []string{
	// 复姓
	"欧阳", "司马", "诸葛", "上官", "东方", "夏侯", "皇甫", "尉迟",
	// 常见单姓
	"李", "王", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴",
	"徐", "孙", "胡", "朱", "高", "林", "何", "郭", "马", "罗",
}

// 常用汉字用于名字符号（精选）
var commonNameChars = []rune([]rune("伟刚勇毅俊峰强军平保东文辉力明永健世广志才涛超波龙磊斌鹏成康星光天达安岩中茂进林有坚和彬富顺信子杰涛昊诚轩凡嘉晨皓"))

var femaleNameChars = []rune([]rune("秀娟英华慧巧美娜静淑惠珠翠雅芝玉萍红娥玲芬芳燕彩春菊兰凤洁梅琳素云莲真环雪荣爱妹霞香月莺媛艳瑞凡佳嘉诗怡菲可欣彤璐梦"))

// 生成一个中国人名（带去重）
func generateChineseName(existing map[string]bool, isMale bool) string {
	for {
		surname := surnames[rand.Intn(len(surnames))]
		var name string

		// 随机生成1或2个字的名字
		length := rand.Intn(2) + 1
		for i := 0; i < length; i++ {
			if isMale {
				name += string(commonNameChars[rand.Intn(len(commonNameChars))])
			} else {
				name += string(femaleNameChars[rand.Intn(len(femaleNameChars))])
			}
		}

		fullName := surname + name

		if !existing[fullName] {
			existing[fullName] = true
			return fullName
		}
	}
}

// 批量生成不重复的人名
type MockName struct {
	Name   string `json:"name"`
	IsMale bool   `json:"isMale"`
}

func GenerateBatchChineseNames(count int) []MockName {
	rand.Seed(time.Now().UnixNano())
	//names := make([]string, 0, count)
	names := make([]MockName, count)

	existing := make(map[string]bool)

	//for len(names) < count {
	//	isMale := gofakeit.Bool() // 随机性别(true 男性， false 女性)
	//	name := generateChineseName(existing, isMale)
	//}

	for i := 0; i < count; i++ {
		isMale := gofakeit.Bool() // 随机性别(true 男性， false 女性)
		name := generateChineseName(existing, isMale)

		names[i] = MockName{
			Name:   name,
			IsMale: isMale,
		}
	}

	return names
}

// 年月生成对应时间戳
func GenTimestamp(input string) (int64, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	layout := "200601"

	t, err := time.Parse(layout, input)
	if err != nil {
		return 0, err
	}

	// 设置时间为 1 号晚上 20:00（默认是零点）
	t = time.Date(t.Year(), t.Month(), 1, 20, 0, 0, 0, loc)
	return t.Unix(), nil
}
