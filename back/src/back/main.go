package main

import (
	"back/api"
	"back/room"

	"github.com/labstack/echo"
)

const BIND = ":8081"

func main() {
	e := echo.New()

	// create room
	room.CreateAndRun()

	e.GET("/get_and_create", api.GetSocketAndCreateRoom)

	e.Logger.Fatal(e.Start(BIND))
}
