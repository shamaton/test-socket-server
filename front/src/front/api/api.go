package api

import (
	"front/room"
	"front/trace"
	"log"
	"os"
	"strconv"

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

func GetSocketAndCreateRoom(c echo.Context) error {

	// create socket
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// create room
	// TODO : make ramdom room id
	roomNo++
	r := room.CreateRoom(roomNo)
	r.SetTracer(trace.New(os.Stdout))

	// create client
	cli := room.CreateClient(r, sock, messageBufferSize)

	// run
	go r.Run()
	cli.Run()

	return nil
}

func GetSocket(c echo.Context) error {

	// get room id
	roomIdStr := c.FormValue("room_id")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("[ERROR] room id is invalid!!", roomIdStr)
		return err
	}

	// get room
	r, err := room.Get(roomId)
	if err != nil {
		log.Println("[ERROR] ", err, roomId)
		return err
	}

	// create socket
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// create client
	cli := room.CreateClient(r, sock, messageBufferSize)

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
