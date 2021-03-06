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
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/labstack/echo"
)

type User struct {
	Id        uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	GroupID   uint   `json:"groupID"`
	Group     Group
}

type Belong struct {
	UserID  uint `json:"userID"`
	GroupID uint `json:"groupID"`
}

type Group struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
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
	Id           uint      `json:"id" gorm:"AUTO_INCREMENT"`
	UserID       uint      `json:"userID" gorm:"foreignKey:UserID"`
	AssignmentID uint      `json:"assignmentID" gorm:"foreignKey:AssignmentID"`
	Status       uint      `json:"status"`
	Ranking      uint      `json:"ranking"`
	UpdateAt     time.Time `json:"updateAt"`
}

func main() {
	db := sqlConnect()
	defer db.Close()
	e := echo.New()
	router(e)
	var port string
	switch env := os.Getenv("env"); env {
	case "production":
		port = os.Getenv("PORT")
	default:
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func router(e *echo.Echo) {
	e.GET("/", root)
	e.GET("/user/:userID", get_user)
	e.GET("/migrate/:command", migrate)
	e.GET("/do", get_dos)
	e.GET("/do/userID/:userID", get_do_spec)
	e.GET("/do/week/:userID", get_do_week)
	e.GET("/group", get_group_all)
	e.GET("/group/:groupID", get_group)
	e.GET("/assignment/:groupID", get_ass)
	e.POST("/user", post_user)
	e.POST("/validate/user", validate_user)
	e.POST("/assignment", post_assignment)
	e.POST("/group", post_group)
	e.POST("/do", post_do)
	e.POST("/share", post_share)
	e.PUT("/belong/group/:groupID", belong_group)
	e.PUT("/do", put_do)
	e.PUT("/user/:userID", put_user)
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
		log.Print("access as production\n")
		DBMS = "postgres"
		URL = os.Getenv("DATABASE_URL")
		URL += "?sslmode=require"
		log.Print("acceccing ", URL)
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
				fmt.Println("DB????????????")
				panic(err)
			}
			db, err = gorm.Open(DBMS, URL)
		}
	}
	fmt.Println("DB????????????")

	return db
}

func migrate(c echo.Context) error {
	command := c.Param("command")

	// go run model.go up
	// ?????????????????????????????????
	switch command {
	case "up":
		upTable()
	case "down":
		downTable()
	}

	// go run model.go down
	// ?????????????????????????????????
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

func get_user(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	userID := c.Param("userID")
	user := new(User)
	if err := db.Find(&user, "id=?", userID); err.Error != nil {
		message := GetDoErrorMessage{"Invalid query param", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}

	message := ValidateUserSuccessMessage{"Find user", true, *user}
	return c.JSON(http.StatusOK, message)
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
	message := ValidateUserSuccessMessage{"Registration successfull", true, *user}

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
		return c.JSON(http.StatusOK, message)
	}

	validate := new(User)

	if err := db.Where("email = ? AND password = ?", user.Email, user.Password).First(&validate); err.Error != nil {
		message := ValidateUserErrorMessage{
			"validation failed",
			false,
			"cannot find records",
		}
		return c.JSON(http.StatusOK, message)
	}
	fmt.Print(validate.Group.Id)
	if validate.Group.Id != 0 {
		db.Model(&validate).Related(&validate.Group.Id)
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

	// t1 ???????????????t2 ????????????
	t1 := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	t2 := t1.AddDate(0, 1, -1)
	fmt.Print(t1, "\n")
	fmt.Print(t2, "\n")

	do := []Do{}
	if err := db.Where("user_id= ? AND status = ? AND update_at >= ? AND update_at < ?", userID, 2, t1, t2).Find(&do); err.Error != nil {
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
		return c.JSON(http.StatusOK, message)
	}

	db.Create(&assignment)

	groupID := assignment.GroupID
	group := new(Group)
	if err := db.Find(&group, "id=?", groupID); err.Error != nil {
		message := GetDoErrorMessage{"Group get failed", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	db.Model(&group).Related(&group.Users)
	fmt.Print("\n", group)
	for _, user := range group.Users {
		db.Create(&Do{
			UserID:       user.Id,
			AssignmentID: assignment.Id,
			UpdateAt:     time.Now(),
		})
	}
	message := PostAssignmentSuccessMessage{"Assignment successfully create", true}

	return c.JSON(http.StatusOK, message)
}

type PostGroupSuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    Group  `json:"data"`
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
	message := PostGroupSuccessMessage{"Group successfully create", true, *group}

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

	userID := c.Param("userID")

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now()
	t1 := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, loc)
	t2 := t1.AddDate(0, 0, -7)

	dos := []Do{}
	if err := db.Where("update_at >= ? AND user_id = ?", t2, userID).Find(&dos); err.Error != nil {
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

	userID, _ := strconv.Atoi(c.QueryParam("userID"))
	assignmentID, _ := strconv.Atoi(c.QueryParam("assignmentID"))
	status, _ := strconv.Atoi(c.QueryParam("status"))
	fmt.Print("userID: ", userID)
	fmt.Print("\nassID: ", assignmentID)
	fmt.Print(status)

	do := new(Do)
	if err := db.Where("user_id = ? AND assignment_id = ?", userID, assignmentID).First(&do); err.Error != nil {
		message := GetDoErrorMessage{"Cant find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	fmt.Print("\ndo: ", do)
	// do.Status = uint(status)
	// db.Save(&do)
	db.Model(&do).Update("Status", uint(status))

	message := SuccessMessage{"Putting do success", true}
	return c.JSON(http.StatusOK, message)
}

func belong_group(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	groupID, _ := strconv.Atoi(c.Param("groupID"))
	userID := c.QueryParam("userID")
	fmt.Print("groupID:", groupID)
	fmt.Print("\nuserID: ", userID)

	user := new(User)

	if err := db.Find(&user, "id=?", userID); err.Error != nil {
		message := GetDoErrorMessage{"Cannot find user record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}

	// user.GroupID = uint(groupID)
	// db.Save(&user)
	db.Model(&user).Update("GroupID", uint(groupID))

	message := SuccessMessage{"Record updated", true}
	return c.JSON(http.StatusOK, message)
}

type GetGroupSuccessMessage struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    Group  `json:"data"`
}

type GetGroupsSuccessMessage struct {
	Message string  `json:"message"`
	Status  bool    `json:"status"`
	Data    []Group `json:"data"`
}

func get_group_all(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	groups := []Group{}
	if err := db.Find(&groups); err.Error != nil {
		message := GetDoErrorMessage{"Cannot find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	for i := range groups {
		db.Model(&groups[i]).Related(&groups[i].Users, "Users")
	}
	return c.JSON(http.StatusOK, GetGroupsSuccessMessage{"Success", true, groups})
}

func get_group(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	groupID := c.Param("groupID")
	group := new(Group)
	if err := db.Find(&group, "id=?", groupID); err.Error != nil {
		message := GetDoErrorMessage{"Cannot find group record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	db.Model(&group).Related(&group.Users)

	message := GetGroupSuccessMessage{"Get group Successful", true, *group}
	return c.JSON(http.StatusOK, message)
}

func put_user(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	userID, _ := strconv.Atoi(c.Param("userID"))
	groupID, _ := strconv.Atoi(c.QueryParam("groupID"))
	user := new(User)
	if err := db.Find(&user, "id=?", userID); err.Error != nil {
		message := GetDoErrorMessage{"Cant find record", false, err.Error}
		return c.JSON(http.StatusOK, message)
	}
	user.GroupID = uint(groupID)
	db.Save(&user)

	message := SuccessMessage{"Putting do success", true}
	return c.JSON(http.StatusOK, message)
}
