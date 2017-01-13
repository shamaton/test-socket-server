package socket

import (
	"fmt"
	"front/convert"
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
	world = iota
	group
	private
)

const (
	_ = iota
	sendMessage
	updateState
	getMemberInfo
)

func (c *connBack) read() {
	for {
		if msgType, msg, err := c.socket.ReadMessage(); err == nil {

			if msgType == websocket.BinaryMessage {

				fmt.Println("front -> back : receive from backend...")
				// convert raw data
				converter := convert.Create(msg)

				switch converter.CommandId() {
				case sendMessage:
					typ, tid, fid, es := c.getRangeInfo(converter)
					if es.HasErr() {
						fmt.Println(es.Err())
						continue
					}

					// NOTE : does not work if server scaling
					c.send2front(typ, tid, msg)
					if typ == private {
						c.send2front(typ, fid, msg)
					}

				case updateState:
					// notice to world
					c.send2front(world, -1, msg)

				case getMemberInfo:
					id := c.getMemberInfo(converter)
					// target and myself
					c.send2front(private, id, converter.PackedData())

				default:
					/*
						game.Dispatch(converter)
						if converter.IsPacked() {
							c.room.broadCastByte <- converter.PackedData()
						}
					*/
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

func (c *connBack) closeSelf(code int, message string) error {
	fmt.Printf("client close : [%d] %s", code, message)
	// close channel
	close(c.fromFront)
	return nil
}

func (c *connBack) getRangeInfo(converter *convert.Converter) (int, int, int, errstack.Stacker) {
	type receiveData struct {
		RangeType int
		RangeId   int
		FromId    int
		Name      string
		Message   string
	}
	var r receiveData
	es := converter.Unpack(&r)
	return r.RangeType, r.RangeId, r.FromId, es
}

func (c *connBack) getMemberInfo(converter *convert.Converter) int {

	type receive struct {
		UserId int
	}
	rec := receive{}
	converter.Unpack(&rec)

	type response struct {
		UserId   int
		UserName string
		Status   int
	}
	infos := []response{}
	for k, _ := range users {
		s := response{}
		s.UserId = k
		s.UserName = user2name[k]
		s.Status = 1
		infos = append(infos, s)
	}
	es := converter.Pack(getMemberInfo, infos)
	if es.HasErr() {
		fmt.Println(es.Err())
	}
	return rec.UserId
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
