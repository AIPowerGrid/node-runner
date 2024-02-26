package exllama

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"runner/core"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	json "github.com/goccy/go-json"
)

var (
	conns     = make(map[int]net.Conn)
	log       = core.GetLogger()
	callBacks = make(map[string]chan Resp)
	// callBacks map[string]chan string{}
)

func Connect(port int) {
	curTry := 0
	var cn net.Conn
	for {
		host := fmt.Sprintf("127.0.0.1:%d", port)
		u := url.URL{Scheme: "ws", Host: host}

		ctx := context.Background()
		c, _, _, err := ws.Dial(ctx, u.String())
		if err != nil {
			log.Warn(err)
			if curTry > 20 {
				panic(err)
			}
			curTry += 1
			time.Sleep(time.Millisecond * 100)
		} else {
			conns[port] = c
			cn = c
			break
		}
	}
	fmt.Println("websocket connection opened...", port)
	go ReadAll(cn)
	// defer c.Close()

}
func Send(port int, msg []byte) error {
	conn, found := conns[port]
	if !found {
		return errors.New("no conn")
	}
	err := wsutil.WriteClientMessage(conn, ws.OpText, msg)
	return err
}
func AddRequestCallback(id string, ch chan Resp) {
	callBacks[id] = ch
}
func RemoveCallback(id string) {
	ch, found := callBacks[id]
	if found {
		delete(callBacks, id)
		close(ch)
	}
}

func ReadAll(conn net.Conn) {
	for {
		msg, op, err := wsutil.ReadServerData(conn)
		if err != nil {
			log.Error(err)
			continue
		}
		if op == ws.OpText {
			// log.Debugf("Received Server Message:%s", string(msg))
			// }
			var resp Resp
			err = json.Unmarshal(msg, &resp)
			if err != nil {
				log.Error(err)
			} else {
				// log.Debugf("Got Resp struct from websocket...")
				ch, found := callBacks[resp.RequestID]
				if found {
					// log.Debugf("sending struct to channel...")
					ch <- resp
				}
			}
			// err := json.Unmarshal(msg)

		}
	}
}
