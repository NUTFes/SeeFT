package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/NUTFes/SeeFT/mysql/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/entity"
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
	Name string
	Mail string
}

type Shift struct {
	UserID  int
	Date    string
	Time    string
	WorkID  int
	Weather string
}

type Work struct {
	Name      string
	URL       string
	President string
	TEL       string
	Place     string
	Color     string
}

func main() {
	if err := userInput(); err != nil {
		log.Fatal(err)
	}

	if err := workInput(); err != nil {
		log.Fatal(err)
	}

	if err := shiftInput(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK. Seed database.")
}

func workInput() error {
	filename := "40_work_list.csv"

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open csv: %w", err)
	}
	defer f.Close()

	db.Init()
	tx := db.GetDB()

	r := csv.NewReader(f)
	record, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	fmt.Println(record)

	for i := 0; i < len(record); i++ {
		work := Work{Name: record[i][0], URL: record[i][1], President: record[i][2], TEL: record[i][3], Place: record[i][4], Color: record[i][5]}
		fmt.Println(work)
		result := tx.Create(&work)
		if result.Error != nil {
			fmt.Println(work)
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

	filename := []string{
		"40_pre.csv",
		"40_current.csv",
	}
	for _, v := range filename {

		f, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("cannot open csv: %w", err)
		}

		var weather string
		if strings.Contains(v, "sunny") {
			weather = "sunny"
		} else if strings.Contains(v, "rainy") {
			weather = "rainy"
		} else {
			weather = "none"
		}

		var date string
		if strings.Contains(v, "pre") {
			date = "preparationDay"
		} else if strings.Contains(v, "current") {
			date = "currentDay"
		} else {
			date = "cleanupDay"
		}

		db.Init()
		tx := db.GetDB()

		r := csv.NewReader(f)

		record, err := r.ReadAll()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		// 39thのシフトと変更点があるので修正必須
		// 学年と局の情報が追加されます。
		for i := 2; i < len(record[0]); i++ {
			for j := 4; j < len(record); j++ {
				var user entity.User

				name := strings.ReplaceAll(record[3][i], " ", "")
				name = strings.ReplaceAll(name, "　", "")
				fmt.Println(name)

				if err := tx.Table("users").Where("name = ?", name).First(&user).Error; err != nil {
					fmt.Println(err)
					i++
					break
				}

				var work entity.Work
				if err := tx.Table("works").Where("name = ?", record[j][i]).First(&work).Error; err != nil {
				}

				shift := Shift{user.ID, date, record[j][0], work.ID, weather}
				fmt.Println(shift)
				result := tx.Create(&shift)
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
	f, err := os.Open("user.csv")
	if err != nil {
		return fmt.Errorf("cannot open csv: %w", err)
	}

	var user User

	db.Init()
	tx := db.GetDB()

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

		user = User{Name: name, Mail: record[5]}
		result := tx.Create(&user)
		if result.Error != nil {
			fmt.Println(user)
			return fmt.Errorf("create db: %w", result.Error)
		}
	}

	return nil
}