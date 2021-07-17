package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type User struct {
	Id   uint
	Name string
}

type Belong struct {
	UserID  uint
	GroupID uint
}

type Group struct {
	Id   uint
	Name string
}

type Share struct {
	GroupID      uint
	AssignmentID uint
}

type Assignment struct {
	Id   uint
	Name string
	Due  time.Time
}

type Do struct {
	UserID       uint
	AssignmentID uint
	Status       uint
	Ranking      uint
	UpdateAt     time.Time
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
