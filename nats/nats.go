package nats

import (
	"fmt"
	"os"
	"runner/core"
	"runner/models"
	"time"

	json "github.com/goccy/go-json"

	"github.com/nats-io/nats.go"
)

var (
	log = core.GetLogger()
)

var NatsConnection *nats.Conn

func GetNC() *nats.Conn {
	return NatsConnection

}
func Start(config models.Config) {
	url := os.Getenv("NATS_URL")

	conn, err := nats.Connect(url, nats.UserInfo("admin", "meamadmin321"), nats.Name(config.MachineID))
	if err != nil {
		log.Fatal(err)
	}
	NatsConnection = conn

	log.Info("[nats] connected to " + url)

}
func RegisterMachine(config models.Config) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	reply, err := NatsConnection.Request("registerMachine", b, time.Second*3)
	if err != nil {
		return err
	}
	s := string(reply.Data)
	log.Info(s)
	return err
}

func GetModel(node models.Node) (ModelResponse, error) {
	// requests a model to run  for the card / cards in the node models.Node ..
	var resp ModelResponse
	b, err := json.Marshal(node)
	if err != nil {
		log.Error(err)
		return resp, err
	}
	c := fmt.Sprintf("requestModel.%s", node.Type)
	reply, err := NatsConnection.Request(c, b, time.Second*5)
	if err != nil {
		log.Error(err)
		return resp, err
	}
	err = json.Unmarshal(reply.Data, &resp)
	if err != nil {
		log.Error(err)
		return resp, err
	}
	log.Debug(resp)
	return resp, nil

}
func comfyRequest(m *nats.Msg) {
	var data models.Job
	err := json.Unmarshal(m.Data, &data)
	if err != nil {
		_returnErr(m, err, "error decoding json", true)
	}

}
