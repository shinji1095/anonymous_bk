package main

import "github.com/labstack/echo"

func router(e *echo.Echo) {
	e.GET("/", root)
	e.GET("/migrate/:command", migrate)
}
