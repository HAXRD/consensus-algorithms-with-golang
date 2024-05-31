package exp

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleOSSignals() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
}
