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

func (r *Room) BroadCast() chan []byte {
	return r.broadCastByte
}

func (r *Room) Run() {
	r.tracer.Trace("ルーム[", r.id, "] が開始されました")

	for {
		select {
		case cli := <-r.join:
			// 参加
			r.clients[cli] = true
			r.tracer.Trace("[", r.id, "]", "新しいクライアントが参加しました")

		case cli := <-r.leave:
			// 退室
			delete(r.clients, cli)
			r.tracer.Trace("[", r.id, "]", "クライアントが退室しました")

		case bytes := <-r.broadCastByte:
			r.tracer.Trace("[", r.id, "]", "データを受信しました: ", len(bytes))
			// すべてのクライアントにメッセージを転送
			for cli := range r.clients {
				select {
				case cli.sendByte <- bytes:
					// メッセージを送信
					r.tracer.Trace("[", r.id, "]", " -- クライアントに送信されました")
				}
			}

		case msg := <-r.broadCastString:
			r.tracer.Trace("[", r.id, "]", "メッセージを受信しました: ", string(msg))
			// すべてのクライアントにメッセージを転送
			for cli := range r.clients {
				select {
				case cli.sendString <- msg:
					// メッセージを送信
					r.tracer.Trace("[", r.id, "]", " -- クライアントに送信されました")
				}
			}
		}
	}
}

func (r *Room) SetTracer(tracer trace.Tracer) {
	r.tracer = tracer
}
