package data_gen

import (
	"testing"

	"github.com/LucienVen/golang-auth-service/bootstrap"
	"github.com/LucienVen/golang-auth-service/internal/utils"
	"github.com/LucienVen/golang-auth-service/pkg/log"
)

func First() {
	bootstrap.Run()
	bootstrap.App.StartHeartbeat()
	log.InitLogger()
}

func TestMock(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()

	err := Mock()
	if err != nil {
		t.Error(err)
	}

	t.Log("Mock data generated successfully")
}

func TestInsertTeachers(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	InsertTeachers()
}

func TestInsertClasses(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	InsertClasses()
}

func TestInsertStudent(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	InsertStudent()
}

func TestInsertExam(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	InsertExam()
}

func TestInsertSubjects(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	InsertSubjects()
}

func TestTeacherSubjectRelation(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	err := TeacherSubjectRelation()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("TeacherSubjectRelation success")
}

func TestMainClassTeacherRelation(t *testing.T) {
	First()
	defer bootstrap.App.CloseApplication()
	MainClassTeacherRelation()

	t.Log("MainClassTeacherRelation success")

}

func TestInsertExamScoreWithRelation(t *testing.T) {
	clean := utils.FirstInit()
	defer clean()

	err := InsertExamScoreWithRelation()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("InsertExamScoreWithRelation success")

}

func TestGetGradeTerm(t *testing.T) {
	examName := "2023年下学期数学考试"
	res := GetGradeTerm(examName)
	t.Log(res)
}
