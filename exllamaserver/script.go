package exllamaserver

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runner/core"
	"runner/exllama"
	"runner/models"
	"strconv"
	"sync"
)

type Server struct {
	NodeID string
	Port   int
	Stdout CapturingPassThroughWriter
	PID    int
	Node   models.Node
}

var (
	log       = core.GetLogger()
	serverMap = map[string]Server{}
)

func ServerByID(nodeID string) (Server, error) {
	found := false
	var target Server
	for _, val := range serverMap {
		if val.NodeID == nodeID {
			target = val
			found = true
			break
		}
	}
	if !found {
		return target, errors.New("not found")

	}
	return target, nil
}
func GetServers() map[string]Server {
	return serverMap
}
func GetServer(nodeID string) (Server, error) {
	x, found := serverMap[nodeID]
	if !found {
		return Server{}, errors.New("no server with that ID")
	}
	return x, nil
}

func GetPort(nodeID string) (int, error) {
	x, found := serverMap[nodeID]
	if !found {
		return 0, errors.New("not found")
	}
	return x.Port, nil
}
func Startup(node models.Node, modelPath string, port int, waitGen bool) {
	exllamaDir := "/home/user/fb/ex/exllamav2"
	validSum := "b34d7d994355d02c41f544c3310367b0"
	valid := validateSum(exllamaDir, validSum)
	if !valid {
		panic("error")
	}
	fmt.Println("sum is valid")
	// log.Info("should be running now but we are panicking because we need to code it first")
	run(node, exllamaDir, modelPath, port, waitGen)
	// run(node, comfyDir, port, waitImg)
}
func run(node models.Node, exllamaDir string, modelPath string, port int, waitGen bool) {
	log.Infof("Starting Exllama socket on port %d for cuda device %s", port, node.CudaDevice)
	nodeID := node.ID
	portStr := strconv.Itoa(port)
	hostStr := fmt.Sprintf("0.0.0.0:%s", portStr)
	cmd := exec.Command("conda", "run", "-n", "ml", "--no-capture-output", "python", "-u", "examples/ws_server.py", "-m", modelPath, "--host", hostStr)
	additionalEnv := fmt.Sprintf("CUDA_VISIBLE_DEVICES=%s", node.CudaDevice)
	newEnv := append(os.Environ(), additionalEnv)
	cmd.Env = newEnv
	cmd.Dir = exllamaDir

	// var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	// stderrIn, _ := cmd.StderrPipe()
	initChan := make(chan string)
	targetString := fmt.Sprintf("Starting WebSocket server on 0.0.0.0 port %d", port)
	stdout := NewCapturingPassThroughWriter(os.Stdout, targetString, initChan)
	// stderr := NewCapturingPassThroughWriter(os.Stderr)
	CheckPort(portStr) // kill process on port if being used

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = io.Copy(stdout, stdoutIn)
		// _, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	// _, errStderr = io.Copy(stderr, stderrIn) // read from stder, we should have this in a goroutine and then read from the val

	fmt.Println("waiting for exllama to start")
	<-initChan
	log.Info("Exllama started, loading model by example request")
	server := Server{Node: node, NodeID: nodeID, Port: port, Stdout: *stdout, PID: cmd.Process.Pid}
	serverMap[nodeID] = server
	log.Infof("Connecting to Library websocket...")
	// time.Sleep(time.Second * 3)
	exllama.Connect(port)
	// panic("ok done")
	req := exllama.CreateReq("Create a hello world python script", "p", 2000)
	respChan := make(chan exllama.Resp)
	exllama.SendReq(req, respChan, port)
	var full_resp string
	for {
		v := <-respChan
		if v.ResponseType == "chunk" {
			full_resp += v.Chunk
			fmt.Printf("%s", v.Chunk)
		} else if v.ResponseType == "full" {
			log.Debugf("question fully answered %d..", port)
			break
		}

	}
	exllama.RemoveCallback(req.RequestID)
	// panic("ok")

	// if waitGen {
	// 	log.Debug("should send request here TODO")

	// }
}
