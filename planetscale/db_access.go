package planetscale

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
CREATE TABLE `reminders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `body` varchar(1024) NOT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `is_done` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
*/

type Reminder struct {
	Id     int `gorm:"primaryKey"`
	Body   string
	IsDone bool
}

func dbSetup() (*gorm.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	_, _ = fmt.Fprintf(os.Stderr, "DATABASE_URL is '%s'\n", dbUrl)
	re := regexp.MustCompile("mysql2://([A-Za-z0-9_]+)@([^/]+)/(\\S+)")
	matches := re.FindStringSubmatch(dbUrl)
	dsn := ""
	if matches != nil && len(matches) > 0 {
		fmt.Printf("Found %d matches:\n", len(matches))
		for _, match := range matches {
			fmt.Printf("  %s\n", match)
		}
		dsn = fmt.Sprintf("%s:@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", matches[1], matches[2])
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot parse dbUrl '%s'\n", dbUrl))
	}

	fmt.Printf("dsn is '%s'\n", dsn)
	// DATABASE_URL is 'mysql2://root@127.0.0.1:3306/firstexample'
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open database: %v\n", err))
	}

	return db, nil
}

func AddReminder(body string) {
	db, _ := dbSetup()
	fmt.Println("AddReminder called with body", body, "db set up.", db)

	todo := Reminder{
		Body: body,
	}
	result := db.Create(&todo)
	if result.Error != nil {
		log.Fatalf("Unable to create a reminder: %v\n", result.Error)
	}
}

func ShowReminders(showUnfinishedOnly bool, showFinishedOnly bool) {
	if showUnfinishedOnly && showFinishedOnly {
		fmt.Println("Error: Flags prevent showing any reminders")
	} else {
		db, _ := dbSetup()
		reminders := make([]Reminder, 0)

		if showUnfinishedOnly {
			fmt.Println("Showing unfinished reminders:")
			result := db.Where("is_done = ?", false).Find(&reminders)
			if result.Error != nil {
				log.Fatalf("Unable to read reminders: %v\n", result.Error)
			}
			fmt.Println("Found", result.RowsAffected)
			for index, reminder := range reminders {
				fmt.Printf("Reminder #%d: %v\n", index, reminder)
			}
		} else if showFinishedOnly {
			fmt.Println("Showing finished reminders:")
			result := db.Where("is_done = ?", true).Find(&reminders)
			if result.Error != nil {
				log.Fatalf("Unable to read reminders: %v\n", result.Error)
			}
			fmt.Println("Found", result.RowsAffected)
			for index, reminder := range reminders {
				fmt.Printf("Reminder #%d: %v\n", index, reminder)
			}
		} else {
			fmt.Println("Showing all reminders:")
			result := db.Find(&reminders)
			if result.Error != nil {
				log.Fatalf("Unable to read reminders: %v\n", result.Error)
			}
			fmt.Println("Found", result.RowsAffected)
			for index, reminder := range reminders {
				fmt.Printf("Reminder #%d: %v\n", index, reminder)
			}
		}
	}
}

func GetReminder(id int) (Reminder, error) {
	db, err := dbSetup()
	if err != nil {
		return Reminder{}, err
	}

	var todo Reminder
	result := db.First(&todo, id)
	if result.Error != nil {
		log.Printf("Unable to retrieve ID %d: %v\n", id, result.Error)
		return Reminder{}, result.Error
	}

	return todo, nil
}

func ToggleReminder(id int) {
	db, _ := dbSetup()

	var todo Reminder
	result := db.First(&todo, id)
	if result.Error != nil {
		log.Fatalf("Unable to retrieve ID %d: %v\n", id, result.Error)
	}

	todo.IsDone = !todo.IsDone
	result = db.Save(&todo)
	if result.Error != nil {
		log.Fatalf("Unable to update ID %d: %v\n", id, result.Error)
	}
}
