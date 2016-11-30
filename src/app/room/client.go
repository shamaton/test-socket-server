package room

import (
	"app/convert"
	"app/game"

	"fmt"

	"github.com/gorilla/websocket"
)

// clientはチャットを行っている1人のユーザーを表します。
type Client struct {
	socket     *websocket.Conn // socket
	sendByte   chan []byte     // send channel
	sendString chan string     // send channel
	room       *Room           // room controller
}

func CreateClient(room *Room, socket *websocket.Conn, sendSize int) *Client {
	client := &Client{
		socket:   socket,
		sendByte: make(chan []byte, sendSize),
		room:     room,
	}
	return client
}

func (c *Client) Room() *Room {
	return c.room
}

func (c *Client) Run() {
	// join myself.
	c.room.join <- c

	// start mode write @ go routine
	go c.writeByte()
	go c.writeString()

	// start mode read
	c.read()
}

func (c *Client) read() {
	for {
		if msgType, msg, err := c.socket.ReadMessage(); err == nil {

			if msgType == websocket.BinaryMessage {
				// データをconverterへ
				converter := convert.Create(msg)

				fmt.Println(converter.CommandId())

				// 退出処理は別
				if converter.CommandId() == 1 {
					c.room.leave <- c
				} else {
					game.Dispatch(converter)
					if converter.IsPacked() {
						c.room.broadCastByte <- converter.PackedData()
					}
					if len(converter.Message) > 0 {

					}
				}
			} else if msgType == websocket.TextMessage {
				/*
					cmd := string(msg)
					if cmd == "enter" {
						c.room.join <- c
					} else if cmd == "leave" {
						c.room.leave <- c
					} else {
						c.room.forward <- msg
					}
				*/
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *Client) writeByte() {
	for bytes := range c.sendByte {
		if err := c.socket.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (c *Client) writeString() {
	for msg := range c.sendString {
		if err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			break
		}
	}
	c.socket.Close()
}
