package main

import (
	"consensus-algorithms-with-golang/pbft"
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
	flag.Parse()

	validators := pbft.NewValidators(pbft.NUM_OF_NODES)
	blockchain := pbft.NewBlockchain(*validators)
	wallet := pbft.NewWallet(*SECRET)
	txPool := pbft.NewTxPool()
	blockPool := pbft.NewBlockPool()
	preparePool := pbft.NewMsgPool()
	commitPool := pbft.NewMsgPool()
	rcPool := pbft.NewMsgPool()

	var peers []string
	if *PEERS != "" {
		peers = strings.Split(*PEERS, ",")
	} else {
		peers = nil
	}

	node := pbft.NewNode(
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
