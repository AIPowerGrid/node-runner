package nats

import (
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func ReturnErr(m *nats.Msg, err error, msg string, logErr bool) {
	_returnErr(m, err, msg, logErr)

}
func _returnErr(m *nats.Msg, err error, msg string, logErr bool) {
	if logErr {
		log.Error(err)
	}
	var r JSResponse
	if msg != "" {
		r = JSResponse{Success: false, Message: msg}
	} else {
		r = JSResponse{Success: false, Message: err.Error()}
	}
	b, _ := json.Marshal(r)
	m.Respond(b)
}

func natsPanic(cb func([]byte) error) {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic:", r)
		s := "Recovered Panic: "
		s += fmt.Sprint(r)
		f := JSResponse{Success: false, Message: s}
		b, err := json.Marshal(f)
		if err == nil {
			cb(b)
		}
	}
}
