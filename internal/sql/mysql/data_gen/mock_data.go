package data_gen

import (
	"fmt"
	"strings"

	"github.com/LucienVen/golang-auth-service/bootstrap"
	"github.com/LucienVen/golang-auth-service/internal/entity"
	"github.com/LucienVen/golang-auth-service/internal/utils"
	"github.com/LucienVen/golang-auth-service/pkg/log"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CreatorMock string = "mock"
	UpdaterMock string = "mock"
)

func Mock() error {

	// 链接数据库
	//db := bootstrap.App.GetDB().Mysql

	_ = gofakeit.Seed(0)

	// 创建老师
	//	学生
	//nowTime := time.Now().Unix()

	InsertSubjects()
	InsertTeachers()
	InsertClasses()
	InsertExam()

	InsertStudent() // 实现了关联班级

	// 教师分配学科
	err := TeacherSubjectRelation()
	if err != nil {
		log.Error("TeacherSubjectRelation err", zap.Error(err))
	}
	// 班级分配班主任
	MainClassTeacherRelation()

	// 成绩表（关联考试 ，关联学生）
	err = InsertExamScoreWithRelation()
	if err != nil {
		log.Error("InsertExamScoreWithRelation err", zap.Error(err))
	}

	return nil
}

// insertStudent 插入学生数据
func InsertStudent() {
	db := bootstrap.App.GetDB()

	studentNames := GenerateBatchChineseNames(300)
	baseMode := entity.GenBaseModel(CreatorMock, UpdaterMock)

	class, _ := entity.GetAllClass()
	students := make([]entity.Student, 0)

	classStudentNum := 10 // 每个班 10 人

	for times, item := range class {
		//fmt.Println(fmt.Sprintf("%d年%d班", item.GradeLevel, item.ClassNum))

		// 截取人名数组
		usedName := studentNames[times*classStudentNum : (times+1)*classStudentNum]
		//log.Info("usedName", zap.Any("usedName", usedName))

		for _, name := range usedName {

			gender := 1
			if !name.IsMale {
				gender = 2
			}

			id, _ := uuid.NewV7()

			students = append(students, entity.Student{
				Id:        id.String(),
				Name:      name.Name,
				Gender:    int64(gender),
				ClassId:   item.Id,
				Phone:     gofakeit.Phone(),
				Email:     gofakeit.Email(),
				Passwd:    gofakeit.Password(true, true, true, true, false, 16),
				BaseModel: baseMode,
			})
		}

	}

	//log.Info("len data", zap.Int("len", len(students)))
	//return

	//分批插入
	//每次插入 10 条
	//batchSize := 10
	//for i := 0; i < len(students); i += batchSize {
	//	handlerData := students[i:min(i+batchSize, len(students))]
	//	res, err := db.NamedExec(`INSERT INTO student (id, name, gender, class_id, phone, email, passwd, is_delete, creator, updater, create_time, update_time)
	//			VALUES (:id, :name, :gender, :class_id, :phone, :email, :passwd, :is_delete, :creator, :updater, :create_time, :update_time)`, handlerData)
	//
	//	if err != nil {
	//		log.Error("insert student error", zap.Error(err))
	//		return
	//
	//	} else {
	//		log.Info("insert student success", zap.Any("students", students), zap.Any("res", res))
	//	}
	//
	//}

	res, err := db.NamedExec(`INSERT INTO student (id, name, gender, class_id, phone, email, passwd, is_delete, creator, updater, create_time, update_time)
				VALUES (:id, :name, :gender, :class_id, :phone, :email, :passwd, :is_delete, :creator, :updater, :create_time, :update_time)`, students)
	if err != nil {
		log.Error("insert student error", zap.Error(err))
		return

	} else {
		log.Info("insert student success", zap.Any("students", students), zap.Any("res", res))
	}
}

// insertClasses 插入班级数据
func InsertClasses() {
	db := bootstrap.App.GetDB()
	classNum := 3
	gradeLevel := 6

	baseMode := entity.GenBaseModel(CreatorMock, UpdaterMock)

	classData := make([]entity.Class, 0)

	for i := 1; i <= gradeLevel; i++ {
		for j := 1; j <= classNum; j++ {

			classId, _ := uuid.NewV7()

			classData = append(classData, entity.Class{
				Id:            classId.String(),
				MainTeacherId: "",
				ClassNum:      int64(j),
				GradeLevel:    int64(i),
				BaseModel:     baseMode,
			})
		}
	}

	res, err := db.NamedExec(`INSERT INTO class (id, class_num, grade_level, main_teacher_id, is_delete, creator, updater, create_time, update_time) 
							values (:id, :class_num, :grade_level, :main_teacher_id, :is_delete, :creator, :updater, :create_time, :update_time)`, classData)

	if err != nil {
		log.Error("insert class error", zap.Error(err))
	} else {
		log.Info("insert class success", zap.Any("class", classData), zap.Any("res", res))
	}
}

// insertTeachers 插入老师数据
func InsertTeachers() {
	db := bootstrap.App.GetDB()
	//teachers := []string{"张伟", "王伟", "王芳", "李伟", "王秀英", "李秀英", "李娜", "张秀英", "刘伟", "张敏"}
	teachers := GenerateBatchChineseNames(20)

	baseMode := entity.GenBaseModel(CreatorMock, UpdaterMock)

	teacherData := make([]entity.Teacher, len(teachers))
	for index, name := range teachers {
		gender := 1
		if !name.IsMale {
			gender = 2
		}

		id, _ := uuid.NewV7()

		teacherData[index] = entity.Teacher{
			Id:        id.String(),
			Name:      name.Name,
			Age:       int64(gofakeit.Number(20, 60)),
			Gender:    int64(gender),
			Phone:     gofakeit.Phone(),
			Email:     gofakeit.Email(),
			Passwd:    gofakeit.Password(true, true, true, true, false, 16),
			BaseModel: baseMode,
		}
	}

	res, err := db.NamedExec(`INSERT INTO teacher (id, name, age, gender, phone, email, passwd, is_delete, creator, updater, create_time, update_time)
				VALUES (:id, :name, :age, :gender, :phone, :email, :passwd, :is_delete, :creator, :updater, :create_time, :update_time)`, teacherData)

	if err != nil {
		log.Error("insert teacher error", zap.Error(err))
	} else {
		log.Info("insert teacher success", zap.Any("teachers", teacherData), zap.Any("res", res))
	}
}

// insertSubjects 插入学科数据
func InsertSubjects() {
	db := bootstrap.App.GetDB()
	subjects := []string{"语文", "数学", "英语", "科学", "体育"}

	subjectData := make([]entity.Subject, len(subjects))

	for index, name := range subjects {
		id, _ := uuid.NewV7()
		baseModel := entity.GenBaseModel(CreatorMock, UpdaterMock)

		subjectData[index] = entity.Subject{
			Id:          id.String(),
			Name:        name,
			Description: gofakeit.Sentence(5),
			DirectorId:  "",
			BaseModel:   baseModel,
		}
	}

	res, err := db.NamedExec(`INSERT INTO subject (id, name, description, is_delete, creator, updater, create_time, update_time)
				VALUES (:id, :name, :description, :is_delete, :creator, :updater, :create_time, :update_time)`, subjectData)

	if err != nil {
		log.Error("insert subjects error", zap.Error(err))
		return
	}

	log.Info("insert subjects success", zap.Any("subjects", subjectData), zap.Any("res", res))

}

// 考试
func InsertExam() {
	db := bootstrap.App.GetDB()
	// 2024年，2023年，
	// 语数英，各 4 场，上下学期各 2

	// 获取学科 ID
	subjects, err := entity.GetSubjectByName("语文", "数学", "英语")
	if err != nil {
		log.Error("get subject err", zap.Error(err))
		return
	}

	//fmt.Println(subjects)

	baseMode := entity.GenBaseModel(CreatorMock, UpdaterMock)

	// 生成考试名称
	//examName := make(map[string]int64, 0) // map[考试名称]学科 id

	examData := make([]entity.Exam, 0)

	for _, subject := range subjects {
		for i, s := range []string{"2023年上学期", "2023年下学期", "2024年上学期", "2024年下学期"} {

			en := s + subject.Name + "考试"
			//examName[s+subject.Name+"考试"] = subject.Id

			year := 2025
			mockCreateTime := "202504"
			switch i {
			case 1:
				year = 2023
				mockCreateTime = "202303"
			case 2:
				year = 2023
				mockCreateTime = "202310"
			case 3:
				year = 2024
				mockCreateTime = "202403"
			case 4:
				year = 2024
				mockCreateTime = "202410"
			default:
			}

			mockTimestamp, _ := GenTimestamp(mockCreateTime)
			baseMode.CreateTime = mockTimestamp
			baseMode.UpdateTime = mockTimestamp

			genId, err := uuid.NewV7()
			if err != nil {
				log.Error("uuid generate error", zap.Error(err))
				panic(err)
			}

			examData = append(examData, entity.Exam{
				Id:        genId.String(),
				Name:      en,
				SubjectId: subject.Id,
				Year:      int64(year),
				BaseModel: baseMode,
			})

		}
	}

	// insert
	res, err := db.NamedExec(`INSERT INTO exam (id, name, subject_id, year, is_delete, creator, updater, create_time, update_time)
				VALUES (:id, :name, :subject_id, :year, :is_delete, :creator, :updater, :create_time, :update_time)`, examData)
	if err != nil {
		log.Error("insert exam error", zap.Error(err))
	}

	log.Info("insert exam success", zap.Any("res", res))

}

// 教师分配学科
func TeacherSubjectRelation() error {
	// 查询老师，查询学科
	subjects, err := entity.GetAllSubject()
	if err != nil {
		return fmt.Errorf("get all subject err: %w", err)
	}
	//
	teachers, err := entity.GetAllTeacher()
	if err != nil {
		return fmt.Errorf("get all teacher err: %w", err)
	}

	subjectLen := len(subjects)

	db := bootstrap.App.GetDB()

	for _, item := range teachers {
		instance := entity.NewTeacherEntity(db, item)

		rand := gofakeit.Number(0, subjectLen-1)
		subjectId := subjects[rand].Id

		err = instance.UpdateSubjectId(subjectId, "lxt-02")
		if err != nil {
			return fmt.Errorf("update subject err: %w, teacherId:%s, subjectId:%s", err, item.Id, subjectId)
		}
	}

	return nil
}

// 班级分配班主任
func MainClassTeacherRelation() {
	// 查询老师，查询班级
	teachers, err := entity.GetAllTeacher()
	if err != nil {
		log.Error("get all teacher err", zap.Error(err))
		return
	}

	classes, err := entity.GetAllClass()
	if err != nil {
		log.Error("get all class err", zap.Error(err))
		return
	}

	teacherNum := len(teachers) - 1

	if teacherNum < len(classes) {
		log.Error("老师数量不足")
		return
	}

	workedTeacher := make(map[string]struct{})

	for _, item := range classes {
		instance := entity.NewClassEntity(bootstrap.App.Mysql, item)

		teacherId := func(num int) string {
			tempId := ""
			// 检查是否已经分配
			for tempId == "" {
				rand := gofakeit.Number(0, num)
				t := teachers[rand].Id

				if _, ok := workedTeacher[t]; !ok {
					tempId = t
					workedTeacher[t] = struct{}{}
				}
			}
			return tempId
		}(teacherNum)

		err = instance.UpdateMainTeacherId(teacherId, "lxt")
		if err != nil {
			log.Error("update main teacher err", zap.Error(err))
			return
		}
	}

	log.Info("班级分配班主任成功")
}

// 成绩表（关联考试 ，关联学生）
func InsertExamScoreWithRelation() error {

	// 获取考试
	exams, err := entity.GetAllExams()
	if err != nil {
		return fmt.Errorf("get all exam err: %w", err)
	}

	// 获取学生
	students, err := entity.GetAllStudent()
	if err != nil {
		return fmt.Errorf("get all student err: %w", err)
	}

	gradeData := make([]entity.Grade, 0)

	basemode := entity.GenBaseModel(CreatorMock, UpdaterMock)

	// 每个学生，每场考试，都有成绩
	for _, exam := range exams {
		examId := exam.Id

		// 判断是上学期，还是下学期
		whichPart := GetGradeTerm(exam.Name)

		for _, student := range students {

			newId, _ := uuid.NewV7()

			gradeData = append(gradeData, entity.Grade{
				Id:        newId.String(),
				StudentId: student.Id,
				SubjectId: exam.SubjectId,
				Year:      exam.Year,
				Score:     utils.FormatFloat2Float(float64(gofakeit.IntRange(30, 99)), 1),
				Term:      whichPart,
				ExamId:    examId,
				BaseModel: basemode,
			})
		}

	}

	db := bootstrap.App.GetDB()
	// insert
	res, err := db.NamedExec(`INSERT INTO grade (id, student_id, subject_id, year, score, term, exam_id, is_delete, creator, updater, create_time, update_time)
				VALUES (:id, :student_id, :subject_id, :year, :score, :term, :exam_id, :is_delete, :creator, :updater, :create_time, :update_time)`, gradeData)
	if err != nil {
		return fmt.Errorf("insert grade err: %w", err)
	}

	log.Info("insert grade success", zap.Any("res", res))
	return nil
}

func GetGradeTerm(name string) int64 {
	whichPart := int64(1)
	if strings.Contains(name, "下") {
		whichPart = 2
	}

	return whichPart
}
