package room

import (
	"app/trace"
	"errors"
)

type Room struct {
	// 部屋番号
	id int
	// forwardは他のクライアントに転送するためのメッセージを保持するチャネルです。
	broadCastByte   chan []byte
	broadCastString chan string
	// joinはチャットルームに参加しようとしているクライアントのためのチャネルです。
	join chan *Client
	// leaveはチャットルームから退室しようとしているクライアントのためのチャネルです
	leave chan *Client
	// clientsには在室しているすべてのクライアントが保持されます。
	clients map[*Client]bool
	// tracerはチャットルーム上で行われた操作のログを受け取ります。
	tracer trace.Tracer
}

// ルーム管理
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
