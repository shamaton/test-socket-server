package socket

import (
	"fmt"
	"front/convert"
	"front/game"
	"net/http"

	"front/errstack"

	"github.com/gorilla/websocket"
)

// todo : tmp
const messageBufferSize = 256

// client is one of users in room.
type connBack struct {
	socket     *websocket.Conn // socket
	fromFront  chan []byte     // receive from front.go
	isFinalize bool
}

var back *connBack

func ConnectBack(url string) error {
	dialer := websocket.Dialer{
		Subprotocols:    []string{},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//header := http.Header{"Accept-Encoding": []string{"gzip"}}
	header := http.Header{}

	socket, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}

	// create back
	back = &connBack{
		socket:     socket,
		fromFront:  make(chan []byte, messageBufferSize),
		isFinalize: false,
	}

	// running
	go back.Run()
	return nil
}

func (c *connBack) Run() {

	// start mode write @ go routine
	go c.writeByte()

	// start mode read
	c.read()
}

const (
	_ = iota
	world
	group
	private
)

func (c *connBack) read() {
	for {
		if msgType, msg, err := c.socket.ReadMessage(); err == nil {

			if msgType == websocket.BinaryMessage {

				fmt.Println("front -> back : receive from backend...")
				// convert raw data
				converter := convert.Create(msg)

				switch converter.CommandId() {
				case 0:
				// todo : leave ?
				case 1:
					typ, id, es := c.getRangeInfo(converter)
					if es.HasErr() {
						fmt.Println(es.Err())
						continue
					}
					c.send2front(typ, id, msg)
				default:

					game.Dispatch(converter)
					if converter.IsPacked() {
						//c.room.broadCastByte <- converter.PackedData()
					}
				}

			}
		} else {
			/* error or close signal */
			break
		}
	}

	// if this line reach, finalize client
	c.finalize()
}

func (c *connBack) writeByte() {
	for bytes := range c.fromFront {
		fmt.Println("front -> back : send to backend")
		if err := c.socket.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
			/* if error occurred, finalize */
			break
		}
	}
	// if this line reach, finalize client
	c.finalize()
}

func (c *connBack) finalize() {
	// already finalized ?
	if c.isFinalize {
		return
	}
	c.isFinalize = true

	// close channel
	close(c.fromFront)

	// socket close
	c.socket.Close()
}

func (c *connBack) close(code int, message string) error {
	fmt.Println("call close conback!!")
	// close channel
	close(c.fromFront)
	return nil
}

func (c *connBack) getRangeInfo(converter *convert.Converter) (int, int, errstack.Stacker) {
	type receiveData struct {
		RangeType int
		RangeId   int
		FromId    int
		Name      string
		Message   string
	}
	var r receiveData
	es := converter.Unpack(&r)
	fmt.Println(r)
	return r.RangeType, r.RangeId, es
}

func (c *connBack) send2front(rangeType, rangeId int, data []byte) {
	switch rangeType {
	case world:
		for _, c := range users {
			c.fromBack <- data
		}
	case group:
		us, ok := group2user[rangeId]
		if !ok {
			return
		}
		for id, _ := range us {
			if user, ok := users[id]; ok {
				user.fromBack <- data
			}
		}
	case private:
		if user, ok := users[rangeId]; ok {
			user.fromBack <- data
		}
	default:
	}
}
