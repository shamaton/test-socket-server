package room

import "github.com/gorilla/websocket"

// client is one of users in room.
type Client struct {
	socket     *websocket.Conn // socket
	sendByte   chan []byte     // send channel
	sendString chan string     // send channel
	isFinalize bool
}

func CreateClient(socket *websocket.Conn, sendSize int) *Client {
	client := &Client{
		socket:     socket,
		sendByte:   make(chan []byte, sendSize),
		sendString: make(chan string),
		isFinalize: false,
	}
	return client
}

func (c *Client) Run() {
	// join myself.
	gRoom.join <- c

	// start mode write @ go routine
	go c.writeByte()
	go c.writeString()

	// start mode read
	c.read()
}

func (c *Client) read() {
	for {
		if msgType, msg, err := c.socket.ReadMessage(); err == nil {
			// broadcast
			if msgType == websocket.BinaryMessage {
				gRoom.broadcastByte <- msg
			} else if msgType == websocket.TextMessage {
				gRoom.broadcastString <- string(msg)
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
	gRoom.leave <- c

	// close channel
	close(c.sendString)
	close(c.sendByte)

	// socket close
	c.socket.Close()
}
