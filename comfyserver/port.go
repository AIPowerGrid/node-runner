package comfyserver

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func CheckPort(port string) error {
	pid, used := isPortUsed(port)
	if !used {
		return nil
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	err = p.Kill()
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		// fmt.Println("loop", i)
		_, used := isPortUsed(port)
		if !used {
			break
		}
		i++
		time.Sleep(time.Millisecond * 400)
	}

	return nil
}

func isPortUsed(port string) (pid int, u bool) {

	defer func() {
		if r := recover(); r != nil {
			pid = 0
			u = false
		}
	}()
	c := fmt.Sprintf("ss -ltnup \"sport = :%s\" | grep %s", port, port)
	cmd := exec.Command("bash", "-c", c)
	out, err := cmd.Output()
	if err != nil { // exit status 1 means nothing on port with grep
		return 0, false
	}
	s := string(out)
	// fmt.Println(s)
	// fmt.Println("above...")
	gg := strings.Split(s, "pid=")
	if len(gg) > 0 {
		second := gg[1]
		pidStr := strings.Split(second, ",")[0]

		pid, _ := strconv.Atoi(pidStr)
		return pid, true
	} else {
		return 0, false
	}
}
func _oldport(port string) (int, bool) {
	c := fmt.Sprintf("ss -ltnup \"sport = :%s\" | grep %s", port, port)
	cmd := exec.Command("bash", "-c", c)
	out, err := cmd.Output()
	if err != nil { // exit status 1 means nothing on port with grep
		return 0, false
	}
	s := string(out)
	// fmt.Println(s)
	// fmt.Println("above...")
	gg := strings.Split(s, "pid=")
	second := gg[1]
	pidStr := strings.Split(second, ",")[0]

	pid, _ := strconv.Atoi(pidStr)
	return pid, true
}
