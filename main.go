package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo"
)

type User struct {
	Id        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Belong struct {
	UserID  uint `json:"userID"`
	GroupID uint `json:"groupID"`
}

type Group struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type Share struct {
	GroupID      uint `json:"groupID"`
	AssignmentID uint `json:"assignmentID"`
}

type Assignment struct {
	Id   uint      `json:"id"`
	Name string    `json:"name"`
	Due  time.Time `json:"due"`
}

type Do struct {
	UserID       uint      `json:"userID"`
	AssignmentID uint      `json:"assignmentID"`
	Status       uint      `json:"status"`
	Ranking      uint      `json:"ranking"`
	UpdateAt     time.Time `json:"updateAt"`
}

func main() {
	db := sqlConnect()
	defer db.Close()
	e := echo.New()
	router(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func router(e *echo.Echo) {
	e.GET("/", root)
	e.GET("/migrate/:command", migrate)
	e.POST("/user", post_user)
}

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello")
}

func sqlConnect() (database *gorm.DB) {
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

func migrate(c echo.Context) error {
	command := c.Param("command")

	// go run model.go up
	// マイグレーションの実行
	switch command {
	case "up":
		upTable()
	case "down":
		downTable()
	}

	// go run model.go down
	// マイグレーションの削除
	return c.JSON(http.StatusOK, "migration complete")
}

func upTable() {
	db := sqlConnect()
	defer db.Close()

	db.AutoMigrate(
		&User{},
		&Belong{},
		&Group{},
		&Share{},
		&Assignment{},
		&Do{},
	)

}

func downTable() {
	db := sqlConnect()
	defer db.Close()

	db.DropTable(
		&User{},
		&Belong{},
		&Group{},
		&Share{},
		&Assignment{},
		&Do{},
	)
}

type PostUserMessageSuccess struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type PostUserMessageFaild struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Error   string `json:"error"`
}

func post_user(c echo.Context) (err error) {
	db := sqlConnect()
	defer db.Close()

	user := new(User)
	if err = c.Bind(user); err != nil {
		message := PostUserMessageFaild{"Registration successfull", false, err.Error()}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	db.Create(&user)
	message := PostUserMessageSuccess{"Registration successfull", true}

	return c.JSON(http.StatusOK, message)
}
