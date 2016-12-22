package main

import (
	"back/api"
	"back/socket"

	"github.com/labstack/echo"
)

const BIND = ":8081"

func main() {
	e := echo.New()

	// create room
	socket.StartRoom()

	e.GET("/", api.GetSocket)
	e.GET("/ping", api.Ping)

	e.Logger.Fatal(e.Start(BIND))
}
