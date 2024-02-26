package comfyserver

import (
	"fmt"
	"os/exec"
	"strings"
)

func validateSum(dir, should string) bool {
	sum, err := getSum(dir)
	if err != nil {
		return false
	}
	if sum == should {
		return true
	}
	return false
}
func getSum(dir string) (string, error) {
	/*
		find . -type f -exec md5sum {} + | LC_ALL=C sort | md5sum
		option from website
		find somedir -type f -exec md5sum {} \; | sort -k 2 | md5sum
		option form forum
		tar -cf - | md5sum other option, but includes metadata unfortunately..
	*/
	// cmd := exec.Command("tar", "-cf","-" ComfyUI | md5sum")
	c := "find . -name \"*.py\" -type f -exec md5sum {} + | LC_ALL=C sort | md5sum"
	cmd := exec.Command("bash", "-c", c)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	s := string(out)
	fmt.Println(s)
	n := strings.Split(s, " ")[0]
	final := strings.TrimSpace(n)
	return final, nil
}
