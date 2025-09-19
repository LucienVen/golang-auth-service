package entity

import "time"

type BaseModel struct {
	IsDelete   int64  `json:"is_delete" db:"is_delete"`     // 是否删除（逻辑删除标记）
	Creator    string `json:"creator" db:"creator"`         // 创建者
	Updater    string `json:"updater" db:"updater"`         // 更新者
	CreateTime int64  `json:"create_time" db:"create_time"` // 创建时间
	UpdateTime int64  `json:"update_time" db:"update_time"` // 更新时间
}

func (bm *BaseModel) GenWhenUpdate(handler string) *BaseModel {
	bm.Updater = handler
	bm.UpdateTime = time.Now().Unix()
	return bm
}

func (bm *BaseModel) GenWhenUpdateMap(handler string) map[string]interface{} {
	bm.GenWhenUpdate(handler)
	return map[string]interface{}{
		"updater":     bm.Updater,
		"update_time": bm.UpdateTime,
	}
}

func GenBaseModel(c, u string) BaseModel {
	timenow := time.Now().Unix()
	return BaseModel{
		IsDelete:   0,
		Creator:    c,
		Updater:    u,
		CreateTime: timenow,
		UpdateTime: timenow,
	}
}

// Class 班级表
type Class struct {
	Id            string `json:"id" db:"id"`
	ClassNum      int64  `json:"class_num" db:"class_num"`             // 班级序号
	GradeLevel    int64  `json:"grade_level" db:"grade_level"`         // 年级
	MainTeacherId string `json:"main_teacher_id" db:"main_teacher_id"` // 班主任
	BaseModel
}

// Exam 考试表
type Exam struct {
	Id        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`             // 考试名称
	SubjectId string `json:"subject_id" db:"subject_id"` // 学科 id
	Year      int64  `json:"year" db:"year"`             // 考试年份
	BaseModel
}

// Grade 学生成绩表
type Grade struct {
	Id        string  `json:"id" db:"id"`
	StudentId string  `json:"student_id" db:"student_id"` // 学生 id
	SubjectId string  `json:"subject_id" db:"subject_id"` // 学科
	Year      int64   `json:"year" db:"year"`             // 考试年份
	Score     float64 `json:"score" db:"score"`           // 成绩
	Term      int64   `json:"term" db:"term"`             // 上学期 1，下学期2
	ExamId    string  `json:"exam_id" db:"exam_id"`       // 试卷 Id
	BaseModel
}

// Role 角色表
type Role struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"` // 角色名
	BaseModel
}

// Student 学生信息表
type Student struct {
	Id      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`     // 姓名
	Gender  int64  `json:"gender" db:"gender"` // 性别
	ClassId string `json:"class_id" db:"class_id"`
	Phone   string `json:"phone" db:"phone"`   // 联系方式
	Email   string `json:"email" db:"email"`   // 邮箱
	Passwd  string `json:"passwd" db:"passwd"` // 登录密码
	BaseModel
}

// Subject 学科信息表
type Subject struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`               // 学科
	Description string `json:"description" db:"description"` // 介绍
	DirectorId  string `json:"director_id" db:"director_id"` // 主任(关联教师 id)
	BaseModel
}

// Teacher 教职人员表
type Teacher struct {
	Id        string `json:"id" db:"id"`                 // 教职人员ID
	Name      string `json:"name" db:"name"`             // 姓名
	Age       int64  `json:"age" db:"age"`               // 年龄
	Gender    int64  `json:"gender" db:"gender"`         // 性别(1-男，2-女)
	SubjectId string `json:"subject_id" db:"subject_id"` // 所教学科
	Phone     string `json:"phone" db:"phone"`           // 联系方式
	Email     string `json:"email" db:"email"`           // 邮箱地址
	Passwd    string `json:"passwd" db:"passwd"`         // 密码（已加密）
	Level     int64  `json:"level" bson:"level"`         // 负责年级
	BaseModel
}
