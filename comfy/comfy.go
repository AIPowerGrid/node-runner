package comfy

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runner/core"
	"time"

	json "github.com/goccy/go-json"
	"github.com/google/uuid"
)

var clientId = uuid.NewString()
var (
	base_filname = _getfile()
	log          = core.GetLogger()
	fileData     = ""
	// apiUrl       = "http://host.docker.internal:8188"
	apiUrl = "http://localhost"
)

type fullPayload struct {
	ClientID string   `json:"client_id"`
	Prompt   Workflow `json:"prompt"`
}

func Generate(prompt string, seed int64, jobid string, port string) error {
	log.Infof("comfy.generate called with prompt %s and seed %d", prompt, seed)
	ts := time.Now().UnixMilli()
	fmt.Println(ts)
	// fmt.Println(fileData)
	data, err := parseJSONPrompt(fileData)
	if err != nil {
		log.Error(err)
		return err

	}

	// #set the text prompt for our positive CLIPTextEncode
	data.Num6.Inputs.Text = prompt
	// data["6"]["inputs"]["text"] = prompt

	// #set the seed for our KSampler node
	data.Num3.Inputs.Seed = seed
	// data["3"]["inputs"]["seed"] = seed

	data.Num9.Inputs.FilenamePrefix = jobid
	// log.Debug(jobid)
	log.Debugf("http://localhost:8188/view?filename=%s_00001_.png", jobid)
	time.Sleep(time.Second * 5)
	fullData := fullPayload{ClientID: clientId, Prompt: data}

	b, err := json.Marshal(fullData)
	if err != nil {
		log.Error(err)
		return err
	}
	// log.Debug(string(b))
	result, err := apiExample(b, port)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug(result)
	return nil

}
func apiExample(b []byte, port string) (map[string]interface{}, error) {
	var result map[string]interface{}
	url := fmt.Sprintf("%s:%s/prompt", apiUrl, port)
	// url = "https://httpbin.org/post"
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: time.Millisecond * 10000}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	// s, _ := ioutil.ReadAll(resp.Body)
	// log.Debug(string(s))
	return result, err

}

func GetImage(filename, port string) (string, error) {
	url := fmt.Sprintf("%s:%s/view?filename=%s", apiUrl, port, filename)
	// log.Info(url)
	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{Timeout: time.Millisecond * 10000}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return "", errors.New("not found")
	}
	// log.Info(resp.StatusCode)
	s, _ := ioutil.ReadAll(resp.Body)
	// log.Debug(string(s))
	// writeFile(string(s))
	return string(s), nil
}
func writeFile(z string) {
	if err := os.WriteFile("file.png", []byte(z), 0666); err != nil {
		log.Warn(err)
	}
}

func readFile() (string, error) {
	file, err := os.Open(base_filname)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil

}
func Startup() error {
	s, err := readFile()
	if err != nil {
		return err
	}
	var p map[string]interface{}
	err = json.Unmarshal([]byte(s), &p)
	if err != nil {
		return err
	}
	fileData = s
	// log.Debugf("Got fileData: %s", s)
	return err

}

func _getfile() string {
	path := "/home/user/projects/aipg/node-runner"
	// path, _ := os.Getwd()
	file := "/comfy/workflow_api.json"
	return fmt.Sprintf("%s%s", path, file)
}

func parseJSONPrompt(f string) (Workflow, error) {
	// fmt.Println(f == "")
	var v Workflow
	err := json.Unmarshal([]byte(f), &v)
	return v, err

}
