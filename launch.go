package main

import (
	"consensus-algorithms-with-golang/pbft2"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	SECRET := flag.String("SECRET", "", "secret key")
	HOST := flag.String("HOST", "localhost", "Hostname")
	WSPORT := flag.Uint64("WSPORT", 8080, "WebSocket port")
	PEERS := flag.String("PEERS", "", "Comma separated list of peers")
	validators := pbft2.NewValidators(pbft2.NUM_OF_NODES)
	blockchain := pbft2.NewBlockchain(*validators)
	wallet := pbft2.NewWallet(*SECRET)
	txPool := pbft2.NewTxPool()
	blockPool := pbft2.NewBlockPool()
	preparePool := pbft2.NewMsgPool()
	commitPool := pbft2.NewMsgPool()
	rcPool := pbft2.NewMsgPool()

	flag.Parse()
	var peers []string
	if *PEERS != "" {
		peers = strings.Split(*PEERS, ",")
	} else {
		peers = nil
	}

	node := pbft2.NewNode(
		*HOST,
		*WSPORT,
		*validators,
		*blockchain,
		*wallet,
		*txPool,
		*blockPool,
		*preparePool,
		*commitPool,
		*rcPool,
	)
	node.Listen(peers)

	// handle system interruption
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
	for _, conn := range node.Sockets {
		conn.Close()
	}
}
