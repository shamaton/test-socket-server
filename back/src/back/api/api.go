package api

import (
	"back/socket"
	"log"
	"net/http"

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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func GetSocket(c echo.Context) error {

	// create socket
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// create client
	cli := socket.CreateClient(sock, messageBufferSize)

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
