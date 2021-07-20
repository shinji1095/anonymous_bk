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
	Id        uint   `json:"id" gorm:"AUTO_INCREMENT"`
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
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Assignment []Assignment
}

type Share struct {
	GroupID      uint `json:"groupID"`
	AssignmentID uint `json:"assignmentID"`
}

type Assignment struct {
	Id      uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Name    string `json:"name"`
	Due     string `json:"due"`
	GroupID int    `json:"groupID"`
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
	e.GET("/do", get_dos)
	e.GET("/do/userID/:userID", get_do_spec)
	e.GET("/do/week", get_do_week)
	e.GET("/assignment/:groupID", get_ass)
	e.POST("/user", post_user)
	e.POST("/validate/user", validate_user)
	e.POST("/assignment", post_assignment)
	e.POST("/group", post_group)
	e.POST("/do", post_do)
	e.POST("/share", post_share)
	e.PUT("/do", put_do)
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
	if err = c.Bind(&user); err != nil {
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
		return c.JSON(http.StatusBadRequest, message)
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
		return c.JSON(http.StatusBadRequest, message)
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
	Error   error  `json:"error"`
}

func get_dos(c echo.Context) (err error) {
	db := sqlConnect()
	defer db.Close()

	loc, _ := time.LoadLocation("Asia/Tokyo")

	userID, _ := strconv.Atoi(c.QueryParam("userID"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	year, _ := strconv.Atoi(c.QueryParam("year"))

	// t1 月初め、　t2 月終わり
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	t2 := t1.AddDate(0, 1, -1)
	fmt.Print(t1, "\n")
	fmt.Print(t2, "\n")

	do := []Do{}
	if err := db.Select("assignment_id, ranking, update_at").Where("user_ id= ? AND status = ? AND update_at >= ? AND update_at < ?", userID, 2, t1, t2).Find(&do); err.Error != nil {
		return c.JSON(http.StatusOK, GetDoErrorMessage{"Can't find records", false, err.Error})
	}

	message := GetDoSuccessMessage{
		"Record successfully get",
		true,
		do,
	}

	return c.JSON(http.StatusOK, message)
}

type PostAssignmentErrorMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Error   string `json:"error"`
}

type PostAssignmentSuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func post_assignment(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	assignment := new(Assignment)
	if err := c.Bind(&assignment); err != nil {
		message := PostAssignmentErrorMessage{"Creation failed", false, err.Error()}
		return c.JSON(http.StatusBadRequest, message)
	}

	db.Create(&assignment)
	message := PostAssignmentSuccessMessage{"Assignment successfully create", true}

	return c.JSON(http.StatusOK, message)
}

func post_group(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	group := new(Group)
	if err := c.Bind(&group); err != nil {
		message := PostAssignmentErrorMessage{"Creation failed", false, err.Error()}
		return c.JSON(http.StatusBadRequest, message)
	}

	db.Create(&group)
	message := PostAssignmentSuccessMessage{"Group successfully create", true}

	return c.JSON(http.StatusOK, message)
}

type ErrorMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Error   string `json:"error"`
}

type SuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func post_do(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	do := new(Do)
	if err := c.Bind(do); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Creation failed", false, err.Error()})
	}
	db.Create(&do)
	message := SuccessMessage{"Creation successfully done", true}
	return c.JSON(http.StatusOK, message)
}

func post_share(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	share := new(Share)
	if err := c.Bind(share); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Creation failed", false, err.Error()})
	}
	db.Create(&share)
	message := SuccessMessage{"Creation successfully done", true}
	return c.JSON(http.StatusOK, message)
}

func get_do_spec(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	dos := []Do{}

	userID := c.Param("userID")
	if err := db.Find(&dos, "user_id = ?", userID); err.Error != nil {
		message := GetDoErrorMessage{"Cant't find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	message := GetDoSuccessMessage{
		"Record successfully get",
		true,
		dos}
	return c.JSON(http.StatusOK, message)
}

func get_do_week(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now()
	t1 := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, loc)
	t2 := t1.AddDate(0, 0, -7)

	dos := []Do{}
	if err := db.Where("update_at >= ?", t2).Find(&dos); err.Error != nil {
		message := GetDoErrorMessage{"Cant't find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}

	message := GetDoSuccessMessage{
		"Record successfully get",
		true,
		dos}
	return c.JSON(http.StatusOK, message)
}

type GetAssMeSuccessMessage struct {
	Message string       `json:"message"`
	Status  bool         `json:"status"`
	Data    []Assignment `json:"data"`
}

func get_ass(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	groupID := c.Param("groupID")

	//group := new(Group)
	assignments := []Assignment{}
	if err := db.Where("group_id = ?", groupID).Find(&assignments); err.Error != nil {
		message := GetDoErrorMessage{"Cant any record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}

	message := GetAssMeSuccessMessage{
		"ok",
		true,
		assignments,
	}
	return c.JSON(http.StatusOK, message)
}

func put_do(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	userID := c.Param("userID")
	assignmentID := c.Param("assignmentID")
	status, _ := strconv.Atoi(c.Param("status"))

	do := new(Do)
	if err := db.Where("user_id = ? AND assignment_id = ?", userID, assignmentID).Find(&do); err.Error != nil {
		message := GetDoErrorMessage{"Cant find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}

	do.Status = uint(status)
	db.Save(&do)

	message := SuccessMessage{"Putting do success", true}
	return c.JSON(http.StatusOK, message)
}
