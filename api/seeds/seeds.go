package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/entity"
)

/*
   0: Index
   1: 局
   2: 部門
   3: 学年
   4: 氏名
   5: メール
   6: 電話
*/
type User struct {
	Name 			string
	Mail 			string
	GradeID			int
	DepartmentID	int
	BureauID		int
	RoleID			int
}

type Shift struct {
	UserID  	int
	TaskID  	int
	YearID		int
	DateID  	int
	TimeID  	int
	WeatherID 	int
}

type Task struct {
	Task      	string
	Place     	string
	URL       	string
	Superviser 	string
	Color     	string
	Notes	  	string
	YearID	  	int
}

func main() {
	if err := userInput(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK. Finish user input.")

	if err := taskInput(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK. Finish task input.")

	if err := shiftInput(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK. Finish shift input.")

	fmt.Println("OK. Seed database.")
}

func taskInput() error {
	filename := "./seeds/41th_task.csv"
	yearID := 42

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open csv: %w", err)
	}
	defer f.Close()

	db.ConnectMySQL()
	tx, err := db.ConnectMySQLFromGorm()
	if err != nil {
		return fmt.Errorf("cannot connect db error: %w", err)
	}

	r := csv.NewReader(f)
	record, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	// fmt.Println(record)

	for i := 1; i < len(record); i++ {
		// taskname := strings.ReplaceAll(record[i][0], " ", "")
		// taskname = strings.ReplaceAll(taskname, "　", "")

		task := Task{Task: record[i][0], Place: record[i][1], URL:record[i][2] ,Superviser: record[i][3], Color: record[i][4], Notes: record[i][5], YearID: yearID}
		// fmt.Println(task)
		result := tx.DB().Create(&task)
		if result.Error != nil {
			fmt.Println(task)
			return fmt.Errorf("create db: %w", result.Error)
		}

	}

	return nil
}

func shiftInput() error {
	/*
		filename := []string{
			"40_pre_sunny.csv",
			"40_pre_rainy.csv",
			"40_current_sunny.csv",
			"40_current_rainy.csv",
			"40_cleanup.csv",
		}
	*/
	yearID := 42
	filename := []string{
		"./seeds/41th_current_1_sunny.csv",
	}
	for _, v := range filename {

		f, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("cannot open csv: %w", err)
		}

		var weatherID int
		if strings.Contains(v, "sunny") {
			weatherID = 1
		} else if strings.Contains(v, "rainy") {
			weatherID = 2
		} else {
			weatherID = 3
		}

		var dateID int
		if strings.Contains(v, "pre") {
			dateID = 1
		} else if strings.Contains(v, "current_1") {
			dateID = 2
		} else if strings.Contains(v, "current_2") {
			dateID = 3
		} else {
			dateID = 4
		}

		db.ConnectMySQL()
		tx, err := db.ConnectMySQLFromGorm()
		if err != nil {
			return fmt.Errorf("cannot connect db error: %w", err)
		}

		r := csv.NewReader(f)

		record, err := r.ReadAll()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		// 39thのシフトと変更点があるので修正必須
		// 学年と局の情報が追加されます。
		for i := 2; i < len(record[0]); i++ {
			var user entity.User

			name := strings.ReplaceAll(record[3][i], " ", "")
			name = strings.ReplaceAll(name, "　", "")
			fmt.Println(name)

			if err := tx.DB().Table("users").Where("name = ?", name).First(&user).Error; err != nil {
				fmt.Println(err)
				i++
				break
			}

			for j := 4; j < len(record); j++ {

				var task entity.Task
				if err := tx.DB().Table("tasks").Where("task = ?", record[j][i]).First(&task).Error; 
				err != nil {
				}

				var time entity.Time
				if err := tx.DB().Table("times").Where("time = ?", record[j][0]).First(&time).Error; 
				err != nil {
				}

				shift := Shift{TaskID:task.ID, UserID:user.ID,  YearID:yearID, DateID:dateID, TimeID:time.ID, WeatherID:weatherID}
				// fmt.Println(shift)
				result := tx.DB().Create(&shift)
				if result.Error != nil {
					fmt.Println(shift)
					//				return fmt.Errorf("create db: %w", result.Error)
				}
			}
		}
		f.Close()
	}
	return nil
}

func userInput() error {
	f, err := os.Open("./seeds/user.csv")
	if err != nil {
		return fmt.Errorf("cannot open csv: %w", err)
	}

	var gradeID int
	var departmentID int
	var bureauID int
	var roleID int
	var user User

	db.ConnectMySQL()
	tx, err := db.ConnectMySQLFromGorm()
	if err != nil {
		return fmt.Errorf("cannot connect db error: %w", err)
	}

	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)

		}

		name := strings.ReplaceAll(record[4], " ", "")
		name = strings.ReplaceAll(name, "　", "")

		mail := strings.ReplaceAll(record[8], " ", "")
		mail = strings.ReplaceAll(mail, "　", "")

		grade := strings.ReplaceAll(record[3], " ", "")
		grade = strings.ReplaceAll(grade, "　", "")
		switch grade {
			case `B1`:
				gradeID = 1
			case `B2`:
				gradeID = 2
			case `B3`:
				gradeID = 3
			case `B4`:
				gradeID = 4
			case `M1`:
				gradeID = 5
			case `M2`:
				gradeID = 6
			case `D1`:
				gradeID = 7
			case `D2`:
				gradeID = 8
			case `D3`:
				gradeID = 9
			case `OB`:
				gradeID = 10
			default:
				gradeID = 0
		}

		department := strings.ReplaceAll(record[7], " ", "")
		department = strings.ReplaceAll(department, "　", "")
		switch department {
			case `未所属`:
				departmentID = 1
			case `機械`:
				departmentID = 2
			case `電気`:
				departmentID = 3
			case `情経`:
				departmentID = 4
			case `生物`, `物生`, `物材`:
				departmentID = 5
			case `環社`:
				departmentID = 6
			case `原子力`:
				departmentID = 7
			default:
				departmentID = 0
		}

		bureau := strings.ReplaceAll(record[1], " ", "")
		bureau = strings.ReplaceAll(bureau, "　", "")
		switch bureau {
			case `委員長`:
				bureauID = 1
			case `副委員長`:
				bureauID = 2
			case `総務局`:
				bureauID = 3
			case `企画局`:
				bureauID = 4
			case `渉外局`:
				bureauID = 5
			case `財務局`:
				bureauID = 6
			case `制作局`:
				bureauID = 7
			case `情報局`:
				bureauID = 8
			default:
				bureauID = 0
		}

		role := strings.ReplaceAll(record[2], " ", "")
		role = strings.ReplaceAll(role, "　", "")
		switch role {
			case `局長`:
				roleID = 2
			default:
				roleID = 1
		}

		if(gradeID != 0){
			user = User{Name: name, Mail: mail, GradeID: gradeID, DepartmentID: departmentID, BureauID: bureauID, RoleID: roleID}
			result := tx.DB().Create(&user)
			if result.Error != nil {
				fmt.Println(user)
				return fmt.Errorf("create db: %w", result.Error)
			}
		}
	}

	return nil
}