package models

import (
	"flag"
	"time"

	"github.com/shinji1095/anonymous_bk/db"
)

type User struct {
	Id   uint
	name string
}

type Belong struct {
	UserID  uint
	GroupID uint
}

type Group struct {
	Id   uint
	name string
}

type Share struct {
	GroupID      uint
	AssignmentID uint
}

type Assignment struct {
	Id   uint
	name string
	due  time.Time
}

type Do struct {
	UserID       uint
	AssignmentID uint
	status       uint
	ranking      uint
	UpdateAt     time.Time
}

func main() {
	flag.Parse()
	command := flag.Arg(0)

	// go run model.go up
	// マイグレーションの実行
	switch command {
	case "up":

	case "down":
		migrate()
	}

	// go run model.go down
	// マイグレーションの削除
}

func migrate() {
	db := db.SqlConnect()
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
	db := db.SqlConnect()
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
