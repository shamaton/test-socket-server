package socket

import (
	"back/trace"

	"os"
)

type Room struct {
	broadcastByte   chan []byte      // to send raw data
	broadcastString chan string      // to send string message
	join            chan *Client     // use if client join
	leave           chan *Client     // use if client leave
	clients         map[*Client]bool // all clients
	tracer          trace.Tracer     // logger
}

var gRoom *Room

func StartRoom() {
	if gRoom != nil {
		return
	}

	gRoom = &Room{
		broadcastByte:   make(chan []byte),
		broadcastString: make(chan string),
		join:            make(chan *Client),
		leave:           make(chan *Client),
		clients:         make(map[*Client]bool),
		tracer:          trace.Off(),
	}
	gRoom.SetTracer(trace.New(os.Stdout))

	go gRoom.Run()
}

func (r *Room) Run() {
	r.tracer.Trace("room : opened")

	for {
		select {
		case cli := <-r.join:
			r.clients[cli] = true
			r.tracer.Trace("room : join new client")

		case cli := <-r.leave:
			if _, isExist := r.clients[cli]; isExist {
				delete(r.clients, cli)
				r.tracer.Trace("room : leave a client")
			}

		case bytes := <-r.broadcastByte:
			r.tracer.Trace("room : receive data: ", len(bytes))
			// broadcast
			for cli := range r.clients {
				select {
				case cli.sendByte <- bytes:
					r.tracer.Trace("room : -- has sent data")
				}
			}

		case msg := <-r.broadcastString:
			r.tracer.Trace("room : receive message: ", string(msg))
			// broadcast
			for cli := range r.clients {
				select {
				case cli.sendString <- msg:
					r.tracer.Trace("room :  -- has sent message")
				}
			}
		}
	}
}

func (r *Room) SetTracer(tracer trace.Tracer) {
	r.tracer = tracer
}
