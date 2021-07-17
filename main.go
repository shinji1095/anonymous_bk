package main

import (
	"net/http"

	"github.com/labstack/echo"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := sqlConnect()
	defer db.Close()
	e := echo.New()
	router(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello")
}
