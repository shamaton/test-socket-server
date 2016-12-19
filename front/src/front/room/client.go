package room

import (
	"front/convert"
	"front/game"

	"fmt"

	"github.com/gorilla/websocket"
)

// client is one of users in room.
type Client struct {
	socket     *websocket.Conn // socket
	sendByte   chan []byte     // send channel
	sendString chan string     // send channel
	room       *Room           // room controller
	isFinalize bool
}

func CreateClient(room *Room, socket *websocket.Conn, sendSize int) *Client {
	client := &Client{
		socket:     socket,
		sendByte:   make(chan []byte, sendSize),
		sendString: make(chan string),
		room:       room,
		isFinalize: false,
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
				// convert raw data
				converter := convert.Create(msg)

				fmt.Println(converter.CommandId())

				// if command is 1, leave
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
				msgStr := string(msg)
				c.room.broadCastString <- msgStr
			}
		} else {
			/* error or close signal */
			break
		}
	}

	// if this line reach, finalize client
	c.Finalize()
}

func (c *Client) writeByte() {
	for bytes := range c.sendByte {
		if err := c.socket.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
			/* if error occurred, finalize */
			break
		}
	}
	// if this line reach, finalize client
	c.Finalize()
}

func (c *Client) writeString() {
	for msg := range c.sendString {
		if err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			/* if error occurred, finalize */
			break
		}
	}
	// if this line reach, finalize client
	c.Finalize()
}

func (c *Client) Finalize() {
	// already finalized ?
	if c.isFinalize {
		return
	}
	c.isFinalize = true

	// leave room
	c.room.leave <- c

	// close channel
	close(c.sendString)
	close(c.sendByte)

	// socket close
	c.socket.Close()
}
