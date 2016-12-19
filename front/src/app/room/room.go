package room

import (
	"app/trace"
	"errors"
)

type Room struct {
	id              int              // room id
	broadCastByte   chan []byte      // to send raw data
	broadCastString chan string      // to send string message
	join            chan *Client     // use if client join
	leave           chan *Client     // use if client leave
	clients         map[*Client]bool // all clients
	tracer          trace.Tracer     // logger
}

// all rooms management
var roomMap = map[int]*Room{}

func CreateRoom(id int) *Room {
	r := &Room{
		id:              id,
		broadCastByte:   make(chan []byte),
		broadCastString: make(chan string),
		join:            make(chan *Client),
		leave:           make(chan *Client),
		clients:         make(map[*Client]bool),
		tracer:          trace.Off(),
	}

	roomMap[id] = r
	return r
}

func Get(id int) (*Room, error) {
	r := roomMap[id]
	if r == nil {
		return nil, errors.New("room not found!!")
	}
	return r, nil
}

func (r *Room) Run() {
	r.tracer.Trace("[", r.id, "] : opened")

	for {
		select {
		case cli := <-r.join:
			r.clients[cli] = true
			r.tracer.Trace("[", r.id, "] : ", "join new client")

		case cli := <-r.leave:
			if _, isExist := r.clients[cli]; isExist {
				delete(r.clients, cli)
				r.tracer.Trace("[", r.id, "] : ", "leave a client")
			}

		case bytes := <-r.broadCastByte:
			r.tracer.Trace("[", r.id, "] : ", "receive data: ", len(bytes))
			// broadcast
			for cli := range r.clients {
				select {
				case cli.sendByte <- bytes:
					r.tracer.Trace("[", r.id, "] : ", " -- has sent data")
				}
			}

		case msg := <-r.broadCastString:
			r.tracer.Trace("[", r.id, "] : ", "receive message: ", string(msg))
			// broadcast
			for cli := range r.clients {
				select {
				case cli.sendString <- msg:
					r.tracer.Trace("[", r.id, "] : ", " -- has sent message")
				}
			}
		}
	}
}

func (r *Room) SetTracer(tracer trace.Tracer) {
	r.tracer = tracer
}
