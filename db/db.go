package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func SqlConnect() (database *gorm.DB) {
	var DBMS string
	var USER string
	var PASS string
	var PROTOCOL string
	var DBNAME string
	var URL string
	var env = os.Getenv("env")

	switch env {
	case "production":
		log.Print("access as production")
		URL = os.Getenv("DATABASE_URL")
	default:
		log.Print("access as development")
		DBMS = "mysql"
		USER = "user"
		PASS = "ppp"
		PROTOCOL = "tcp(db:3306)"
		DBNAME = "another"
		URL = USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	}

	count := 0
	db, err := gorm.Open(DBMS, URL)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				fmt.Println("DB接続失敗")
				panic(err)
			}
			db, err = gorm.Open(DBMS, URL)
		}
	}
	fmt.Println("DB接続成功")

	return db
}
