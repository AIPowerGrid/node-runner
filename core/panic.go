package core
import (
	"fmt"
)
func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic:", r)
	}
}