package api

import (
	"app/room"
	"app/trace"
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

	// ソケットを生成
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// ルーム生成
	// TODO : 何かしらのランダム値をkeyにする
	roomNo++
	r := room.CreateRoom(roomNo)
	r.SetTracer(trace.New(os.Stdout))

	// クライアント生成
	cli := room.CreateClient(r, sock, messageBufferSize)

	// ルーム稼働開始
	go r.Run()

	// ソケット通信開始
	cli.Run()

	return nil
}

func GetSocket(c echo.Context) error {

	// 部屋番号を取得
	roomIdStr := c.FormValue("room_id")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		log.Println("[ERROR] room id is invalid!!", roomIdStr)
		return err
	}

	// 部屋取得
	r, err := room.Get(roomId)
	if err != nil {
		log.Println("[ERROR] ", err, roomId)
		return err
	}

	// ソケットを生成
	sock, err := getSocket(c)
	if err != nil {
		log.Println("[FAILED] create socket:", err)
		return err
	}

	// クライアント生成
	cli := room.CreateClient(r, sock, messageBufferSize)

	// ソケット通信開始
	cli.Run()

	return nil
}

func getSocket(c echo.Context) (*websocket.Conn, error) {
	w := c.Response().Writer()
	h := c.Request()
	socket, err := upGrader.Upgrade(w, h, nil)
	return socket, err
}
