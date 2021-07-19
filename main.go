package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	e.GET("/do", get_does)
	e.POST("/user", post_user)
	e.POST("validate/user", validate_user)

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
		message := PostUserMessageFaild{"Registration failed", false, err.Error()}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	db.Create(&user)
	message := PostUserMessageSuccess{"Registration successfull", true}

	return c.JSON(http.StatusOK, message)
}

type ValidateUserErrorMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Error   string `json:"error"`
}

type ValidateUserSuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	User    User   `json:"user"`
}

func validate_user(c echo.Context) (err error) {
	db := sqlConnect()
	defer db.Close()

	user := new(User)
	if err = c.Bind(user); err != nil {
		message := ValidateUserErrorMessage{"vaiudation failed", false, err.Error()}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	validate := new(User)
	fmt.Print(user.Email)
	fmt.Print(user.Password)

	if err := db.Where("email = ? AND password = ?", user.Email, user.Password).First(&validate); err.Error != nil {
		message := ValidateUserErrorMessage{
			"validation failed",
			false,
			"cannot find any records",
		}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	message := ValidateUserSuccessMessage{
		"Validation successfull",
		true,
		*validate,
	}

	return c.JSON(http.StatusOK, message)
}

type GetDoSuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    []Do   `json:"data"`
}

type GetDoErrorMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func get_does(c echo.Context) (err error) {
	db := sqlConnect()
	defer db.Close()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	fmt.Println(time.Date(2014, 12, 31, 8, 4, 18, 0, loc))

	userID, _ := strconv.Atoi(c.QueryParam("userID"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	year, _ := strconv.Atoi(c.QueryParam("year"))

	// t1 月初め、　t2 月終わり
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	t2 := t1.AddDate(0, 1, -1)
	fmt.Print(t1, "\n")
	fmt.Print(t2, "\n")

	do := []Do{}
	if err := db.Select("ranking, updateAt").Where("userID = ? AND status = ? AND updateAt >= ? AND updateAt < ?", userID, 2, t1, t2).Find(&do); err.Error != nil {
		return c.JSON(http.StatusBadRequest, GetDoErrorMessage{"Can't find records", false})
	}

	message := GetDoSuccessMessage{
		"Record successfully get",
		true,
		do,
	}

	return c.JSON(http.StatusOK, message)
}
