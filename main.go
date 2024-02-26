package main

import (
	"fmt"
	"io/ioutil"
	"runner/comfy"
	"runner/comfyserver"
	"runner/core"
	"runner/exllama"
	"runner/exllamaserver"
	"runner/models"
	"runner/nats"
	"sync"
	"time"

	natsLib "github.com/nats-io/nats.go"

	json "github.com/goccy/go-json"
)

var (
	log        = core.GetLogger()
	_modelsDir = "/home/user/gen-models/"
)

var globalConfig models.Config

func validateConfig(config models.Config) bool {
	valid := true
	if config.MachineID == "" || config.OwnerID == "" {
		valid = false
	}
	nodes := config.Nodes
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		if n.ID == "" || n.VRAM == 0 || n.CudaDevice == "" || n.GPUCount == 0 || n.Type == "" {
			valid = false
		}
	}
	return valid
}
func getConfig() models.Config {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("could not read configuration file")
	}
	var config models.Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Error(err)
		// log.Fatal("Invalid configuration file..")
		panic("invalid configuration file")
	}
	valid := validateConfig(config)
	if !valid {
		panic("invalid config")
	}
	return config
}
func closeAllServers() {

}
func main() {
	defer log.Sync()
	// defer closeAllServers()
	config := getConfig()
	fmt.Println(config)
	fmt.Println(1, config.MachineID)
	comfy.Startup()
	log.Info("connecting to nats server..")
	nats.Start(config)
	log.Info("nats setup done")
	/* step 1 - get model from backend */
	nodes := PickModels(config)
	config.Nodes = nodes
	globalConfig = config
	/* step 2 - run model loaded on gpus  + run first prompt to load model cached*/
	imageNodes := filterNodes("image", nodes)
	textNodes := filterNodes("text", nodes)

	var wg sync.WaitGroup
	if len(imageNodes) > 0 && false {
		log.Info("loading image nodes on gpus")
		wg.Add(1)
		go func() {
			loadComfyGPUS(imageNodes)
			wg.Done()
		}()
	}
	if len(textNodes) > 0 {
		log.Info("Loading Text nodes on GPUs")
		wg.Add(1)
		go func() {
			loadTextGPUS(textNodes)
			log.Info("Text nodes loading done...")
			wg.Done()
		}()

	}
	log.Info("Waiting for All GPU loading to finish")
	wg.Wait()
	log.Info("GPU loading finished, registering with backend to serve inference reqeuests..")

	/* step 3 - send nodes to get saved in backend and start serving requests */
	err := nats.RegisterMachine(config)
	if err != nil {
		// log.Info("Error occured init machine..")
		log.Error(err)
		panic("could not init machine for listening to requests")
	}
	log.Info("Done registering machine")
	/* step 4 - register listeners to requests for all nodes */
	log.Info("Registering Listeners")
	listenImage(imageNodes)
	listenText(textNodes)

	log.Info(nodes)
	select {}

}
func listenText(nodes []models.Node) {

	nc := nats.GetNC()
	for _, node := range nodes {
		c := fmt.Sprintf("node.textgenrequest.%s", node.ID)
		log.Infof("shouldbe listen on %c", c)
		nc.Subscribe(c, func(m *natsLib.Msg) {
			var job models.Job
			err := json.Unmarshal(m.Data, &job)
			if err != nil {
				log.Error(err)
				nats.ReturnErr(m, err, "invalid payload", false)
				return
			}
			log.Infof("Received Job to solve: %v", job)
			textResp := fmt.Sprintf("textgenstream.%s", job.ID)
			m.Respond([]byte("OK"))
			server, err := exllamaserver.ServerByID(node.ID)
			if err != nil {
				log.Error(err)
				return
			}
			req := exllama.CreateReq(job.Prompt, "p", 2000)
			respChan := make(chan exllama.Resp)
			exllama.SendReq(req, respChan, server.Port)
			for {
				v := <-respChan
				if v.ResponseType == "chunk" {
					t := models.TextStream{
						RequestID:    job.ID,
						ResponseType: "chunk",
						Chunk:        v.Chunk,
					}
					b, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
						continue
					} else {
						nc.Publish(textResp, b)
					}
					// fmt.Printf("%s", v.Chunk)
				} else if v.ResponseType == "full" {
					t := models.TextStream{
						RequestID:    job.ID,
						ResponseType: "full",
						Response:     v.Response,
					}
					b, err := json.Marshal(t)
					if err != nil {
						log.Error(err)
						break
					} else {
						nc.Publish(textResp, b)
					}
					// log.Debugf("question fully answered %d..", server.Port)
					break
				}

			}
			exllama.RemoveCallback(req.RequestID)
			// panic("ok")

		})
	}
}
func listenImage(nodes []models.Node) {
	// for _, node := range nodes {
	// c := fmt.Sprintf("comfyNodeRequest.%s", node.ID)
	// }

}
func filterNodes(t string, nodes []models.Node) []models.Node {
	var nn []models.Node
	for _, node := range nodes {

		if node.Model == "sdxl" && t == "image" {
			nn = append(nn, node)
		} else if node.Type == "text" && t == "text" {
			nn = append(nn, node)
		}
	}
	return nn
}
func getModelPath(model string) string {
	// phindStr := "TheBloke_Phind-CodeLlama-34B-v2-GPTQ"
	return _modelsDir + model
}
func loadTextGPUS(nodes []models.Node) {
	startPort := 7865
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		curPort := startPort + i
		nodeID := n.ID
		// if i < 100 {
		log.Infof("starting exllama server on port %d for node %s", curPort, nodeID)
		wg.Add(1)
		go func() {
			modelPath := getModelPath(n.Model)
			log.Infof("Starting Exllama server for model %s", modelPath)
			exllamaserver.Startup(n, modelPath, curPort, true)
			wg.Done()

		}()
		// }
	}
	log.Info("waiting for exllama servers to get running + 1 test prompt")
	wg.Wait()
	elap := time.Since(start)
	log.Infof("exllama servers setup in %v", elap)
}
func loadComfyGPUS(nodes []models.Node) {
	startPort := 8190
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		curPort := startPort + i
		nodeID := n.ID
		log.Infof("starting comfy server on port %d for node %s", curPort, nodeID)
		// if i < 1 {
		wg.Add(1)
		go func() {
			comfyserver.Startup(n, curPort, true)
			wg.Done()
		}()
		// }

	}
	log.Info("waiting for comfy servers to get running + served example request...")

	wg.Wait()
	elap := time.Since(start)
	log.Infof("servers setup in %v...", elap)
}
func PickModels(config models.Config) []models.Node {
	nodes := config.Nodes
	var newNodes []models.Node
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		log.Infof("Setting up node %s", node.ID)
		val, err := nats.GetModel(node)
		if err != nil {
			panic("error during initiation")
		}
		node.Model = val.Model
		node.Type = val.Type
		log.Infof("Got model %s to run on node %s, type: %s", node.Model, node.ID, node.Type)
		newNodes = append(newNodes, node)
		// log.Info(val)
	}

	return newNodes
}
