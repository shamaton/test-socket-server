/*
  Package error message stacker.
*/
package errstack

/*
  error.go

  エラーメッセージをスタックする
 */
import (
	"runtime"
	"strings"
)

type errArgs []interface{}

type Stacker struct {
	errMsgs errArgs
}

/*
  オブジェクトを生成

  argument
    msg : エラーメッセージ（なくてもOK）

  return
    生成済オブジェクト
*/
func NewStack(msg ...interface{}) Stacker {

	ew := Stacker{}

	// msgがある場合追加しておく
	if len(msg) > 0 {
		ew.errMsgs = append(ew.errMsgs, msg...)
		ew = ew.addCallerMsg(2)
	}

	return ew
}

/*
  エラーメッセージを追加する

  argument
    msg : エラーメッセージ（なくてもOK）

  return
    スタック済オブジェクト
*/
func (this Stacker) Write(msg ...interface{}) Stacker {
	this.errMsgs = append(this.errMsgs, msg...)
	// 呼び出し元
	this = this.addCallerMsg(2)
	return this
}

/*
  エラーが発生したか

  return
    true or false
*/
func (this Stacker) HasErr() bool {
	if len(this.errMsgs) > 0 {
		return true
	}
	return false
}

/*
  スタックしたエラーを取得する

  return
    エラーメッセージ
*/
func (this Stacker) Err() errArgs {
	return this.errMsgs
}

/*
  呼び出し元の情報をエラーに追加する

  argument
    skip : callerに渡すskip値

  return
    スタック済オブジェクト
*/
func (this Stacker) addCallerMsg(skip int) Stacker {
	// 呼び出し元
	pc, file, line, _ := runtime.Caller(skip)
	callerName := runtime.FuncForPC(pc).Name()

	// 定型文
	addArgs := this.fixedPhrase(file, line, callerName)

	// 追加
	this.errMsgs = append(this.errMsgs, addArgs...)

	return this
}

/*
  callerの情報を整形する

  argument
    file       : ファイル名
    line       : 行数
    callerName : 呼び出し元

  return
    メッセージ
*/
func (this Stacker) fixedPhrase(file string, line int, callerName string) errArgs {
	// 一旦、srcでフィルタする
	splits := strings.Split(file, "/src/")

	// 仮に区切れなくてもエラーにせずそのまま利用する
	fileName := file
	if len(splits) == 2 {
		fileName = splits[1]
	}
	addArgs := errArgs{"(" + callerName + ")", "at", fileName, "line", line}
	return addArgs
}

/*
  unshiftする

  argument
    v : エラーメッセージ（なくてもOK）

  return
    スタック済オブジェクト
*/
func (this Stacker) unshift(v ...interface{}) Stacker {
	this.errMsgs = append(v, this.errMsgs...)
	return this
}
