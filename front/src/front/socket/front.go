package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// client is one of users in room.
type Client struct {
	*websocket.Conn             // socket
	fromBack        chan []byte // receive from back.go
	isFinalize      bool
}

func CreateClient(userId, groupId int, socket *websocket.Conn, sendSize int) *Client {
	client := &Client{
		socket,
		make(chan []byte, sendSize),
		false,
	}
	client.fromBack = make(chan []byte, sendSize)
	client.SetCloseHandler(client.closeSocket)

	// register map
	users[userId] = client
	user2group[userId] = groupId
	group2user[groupId] = append(group2user[groupId], groupId)
	return client
}

func (c *Client) Run() {

	// start mode write @ go routine
	go c.writeByte()

	// start mode read
	c.read()
}

func (c *Client) read() {
	for {
		if msgType, msg, err := c.ReadMessage(); err == nil {
			if msgType == websocket.BinaryMessage {
				back.fromFront <- msg
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
	// message from back
	for bytes := range c.fromBack {
		if err := c.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
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

	// close channel
	close(c.fromBack)

	// socket close
	c.Close()
}

func (c *Client) closeSocket(code int, message string) error {
	fmt.Println("call close!!")
	// close channel
	close(c.fromBack)
	return nil
}
