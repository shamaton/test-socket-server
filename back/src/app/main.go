package main

import "github.com/labstack/echo"

const BIND = ":8081"

func main() {
	e := echo.New()

	//e.GET("/get_and_create", api.GetSocketAndCreateRoom)
	//e.GET("/get", api.GetSocket)
	e.Logger.Fatal(e.Start(BIND))
}
