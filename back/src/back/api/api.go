package api

import (
	"back/room"
	"log"

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

func GetSocketAndCreateRoom(c echo.Context) error {

	// create socket
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// create client
	cli := room.CreateClient(sock, messageBufferSize)

	// run
	cli.Run()

	return nil
}

func getSocket(c echo.Context) (*websocket.Conn, error) {
	w := c.Response().Writer()
	h := c.Request()
	socket, err := upGrader.Upgrade(w, h, nil)
	return socket, err
}
