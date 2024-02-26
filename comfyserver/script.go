package comfyserver

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runner/comfy"
	"runner/core"
	"runner/models"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
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
func Startup(node models.Node, port int, waitImg bool) {
	comfyDir := "/home/user/projects/aipg/comfy/goversion"
	validSum := "176f0fd5e1c190f06a0abf801d8ac3c3"
	valid := validateSum(comfyDir, validSum)
	if !valid {
		panic("error")
	}
	fmt.Println("sum is valid")
	run(node, comfyDir, port, waitImg)
}
func run(node models.Node, comfyDir string, port int, waitImg bool) {
	nodeID := node.ID
	portStr := strconv.Itoa(port)
	cmd := exec.Command("conda", "run", "-n", "comfy", "--no-capture-output", "python", "-u", "main.py", "--listen", "0.0.0.0", "--port", portStr)
	additionalEnv := fmt.Sprintf("CUDA_VISIBLE_DEVICES=%s", node.CudaDevice)
	newEnv := append(os.Environ(), additionalEnv)
	cmd.Env = newEnv
	cmd.Dir = comfyDir

	// var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	// stderrIn, _ := cmd.StderrPipe()
	initChan := make(chan string)
	targetString := fmt.Sprintf("To see the GUI go to: http://0.0.0.0:%d", port)
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

	fmt.Println("waiting for comfyui to start")
	<-initChan
	fmt.Println("ComfyUI is started!!!! for realz... loading model by example request")
	imgId := uuid.NewString()
	seed := time.Now().UnixMilli()
	err = comfy.Generate("Chicken", seed, imgId, portStr)
	if err != nil {
		fmt.Println(err)
		panic("could not generate first comfy request")
	}
	fmt.Println("sent request to comfy to populate cache..")
	server := Server{Node: node, NodeID: nodeID, Port: port, Stdout: *stdout, PID: cmd.Process.Pid}
	serverMap[nodeID] = server

	filename := fmt.Sprintf("%s_00001_.png", imgId)

	if waitImg {
		waitImgReady(filename, portStr)
	}
	// go PeriodicImageGen()

	// wg.Wait() not using wg.wait

	/* err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatalf("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr) */
}

func waitImgReady(filename, port string) {
	for {
		_, err := comfy.GetImage(filename, port)
		if err != nil {
			// log.Debugf("Image not ready yet on port %s", port)
			time.Sleep(time.Millisecond * 500)
		} else {
			log.Infof("Image is ready on port %s", port)
			break
		}
	}
}
func PeriodicImageGen(server Server) {

}
