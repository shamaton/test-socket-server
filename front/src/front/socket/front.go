package socket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// client is one of users in room.
type Client struct {
	userId     int
	groupId    int
	socket     *websocket.Conn // socket
	fromBack   chan []byte     // receive from back.go
	isFinalize bool
}

func CreateClient(userId, groupId int, socket *websocket.Conn, sendSize int) *Client {
	client := &Client{
		userId:     userId,
		groupId:    groupId,
		socket:     socket,
		fromBack:   make(chan []byte, sendSize),
		isFinalize: false,
	}
	client.fromBack = make(chan []byte, sendSize)
	client.socket.SetCloseHandler(client.close)

	// register map
	users[userId] = client
	user2group[userId] = groupId

	_, exist := group2user[groupId]
	if !exist {
		group2user[groupId] = map[int]bool{}
	}
	group2user[groupId][userId] = true
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
		if msgType, msg, err := c.socket.ReadMessage(); err == nil {
			if msgType == websocket.BinaryMessage {
				fmt.Println("user -> front : send to back")
				back.fromFront <- msg
			}
		} else {
			/* error or close signal */
			break
		}
	}

	// if this line reach, finalize client
	c.finalize()
}

func (c *Client) writeByte() {
	// message from back
	for bytes := range c.fromBack {
		fmt.Println("back -> front : send to user")
		if err := c.socket.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
			/* if error occurred, finalize */
			break
		}
	}
	// if this line reach, finalize client
	c.finalize()
}

func (c *Client) finalize() {
	// already finalized ?
	if c.isFinalize {
		return
	}
	c.isFinalize = true

	// delete from map
	delete(users, c.userId)
	delete(user2group, c.userId)
	delete(group2user[c.groupId], c.userId)

	// close channel
	close(c.fromBack)

	// socket close
	c.socket.Close()
}

func (c *Client) close(code int, message string) error {
	fmt.Println("call close!!", code, message)
	// finalize
	c.finalize()
	return nil
}
