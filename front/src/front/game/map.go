package game

import "front/convert"

// 受信コマンド一覧
// TODO : 削除するとずれてしまうのでenumは使うべきではない気がする
const (
	_ = iota
	LeaveRoom
	ReceiveMessage
	Test
	R_1000 = iota + 1000
	R_1001
)

var mapper = map[uint32]func(*convert.Converter){
//1: leaveRoom,
//2: receiveMessage,
}

func Dispatch(converter *convert.Converter) {
	//f := mapper[converter.CommandId()]
	//f(converter)
}

/*
func leaveRoom(converter *convert.Converter) {
	// データ確認
	var dummy bool
	converter.Unpack(&dummy)
	fmt.Println("ret dummy ->", dummy)
	// 退出する

}



// TODO : エラーチェックなど
func receiveMessage(converter *convert.Converter) {
	type receiveData struct {
		rangeType int
		rangeId int
		fromId int
		message string
	}
	r := new(receiveData)
	converter.Unpack(r)

	converter.Pack(2, message)
	fmt.Println("message ->", message)
}
*/
// 応答コマンド一覧
