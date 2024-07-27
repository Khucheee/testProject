package closer

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var CloseFunctions []func()

func CtrlC() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Receive signal to stop working")
	for iterator := len(CloseFunctions) - 1; iterator >= 0; iterator-- {
		CloseFunctions[iterator]()
	}
	//for _, closeFunction := range CloseFunctions {
	//	closeFunction()
	//}

}
