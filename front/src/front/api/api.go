package api

import (
	"front/socket"
	"log"

	"strconv"

	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upGrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

var roomNo = 0

func GetSocket(c echo.Context) error {

	// get user id
	userIdStr := c.FormValue("uid")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("[ERROR] user id is invalid!!", userIdStr)
		return err
	}

	groupIdStr := c.FormValue("gid")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		log.Println("[ERROR] group id is invalid!!", groupIdStr)
		return err
	}

	userName := c.FormValue("name")
	if userName == "" {
		log.Println("[ERROR] user name is empty!!")
		return errors.New("name is empty")
	}

	// debug
	fmt.Println("userid : groupid ", userIdStr, groupIdStr)

	// register map
	if socket.IsExistUser(userId) {
		return errors.New("user has already existed!!")
	}

	// create socket
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// create client
	cli := socket.CreateClient(userId, groupId, userName, sock, messageBufferSize)

	// run
	cli.Run()

	return nil
}

func getSocket(c echo.Context) (*websocket.Conn, error) {
	w := c.Response().Writer
	h := c.Request()
	socket, err := upGrader.Upgrade(w, h, nil)
	return socket, err
}

func Ping(c echo.Context) error {
	c.String(200, "pong")
	return nil
}
