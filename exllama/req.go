package exllama

import (
	"fmt"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/google/uuid"
)

var (
	phindTemplate = `### System Prompt
You are an intelligent programming assistant.

### User Message
{instruction}

### Assistant
`
	wizTemplate = `Below is an instruction that describes a task. Write a response that appropriately completes the request.

### Instruction:
{instruction}

### Response:
`
)

func init() {
	fmt.Println(phindTemplate)
}
func N() {

}

type Req struct {
	Action    string `json:"action"`
	RequestID string `json:"request_id"`
	Stream    bool   `json:"stream"`
	Text      string `json:"text"`
	MaxTokens int    `json:"max_new_tokens"`
}
type Resp struct {
	Action       string `json:"action"`
	RequestID    string `json:"request_id"`
	ResponseType string `json:"response_type"`
	Chunk        string `json:"chunk"`
	Response     string `json:"response,omitempty"`
}

/*
action: 'infer',
  request_id: 'a1e31422',
  stream: true,
  text: '### System Prompt\n' +
    'You are an intelligent programming assistant.\n' +
    '\n' +
    '### User Message\n' +
    'what is golang\n' +
    '\n' +
    '\n' +
    '### Assistant\n',
  max_new_tokens: 2000
}
*/

func MakePrompt(mode string, prompt string) string {
	var z string
	if mode == "w" {
		z = strings.ReplaceAll(wizTemplate, "{instruction}", prompt)
	} else {
		z = strings.ReplaceAll(phindTemplate, "{instruction}", prompt)

	}
	return z

}
func CreateReq(prompt, mode string, maxTokens int) Req {
	id := uuid.NewString()
	full_prompt := MakePrompt(mode, prompt)
	r := Req{
		Action:    "infer",
		RequestID: id,
		Stream:    true,
		Text:      full_prompt,
		MaxTokens: maxTokens,
	}
	return r

}

func SendReq(r Req, cb chan Resp, port int) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	AddRequestCallback(r.RequestID, cb)
	err = Send(port, b)
	return err
}
