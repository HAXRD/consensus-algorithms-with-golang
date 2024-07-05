package main

import (
	"consensus-algorithms-with-golang/pow"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	SECRET := flag.String("SECRET", "", "secret key")
	HOST := flag.String("HOST", "localhost", "host")
	WSPORT := flag.Uint64("WSPORT", 8090, "websocket port")
	PEERS := flag.String("PEERS", "", "comma seperated list of peers")
	N_MINERS := flag.Uint64("N_MINERS", 1, "number of miners")
	flag.Parse()

	validators := pow.NewValidators(pow.NUM_OF_NODES)
	blockchain := pow.NewBlockchain(*validators)
	wallet := pow.NewWallet(*SECRET)
	txPool := pow.NewTxPool()

	var peers []string
	if *PEERS != "" {
		peers = strings.Split(*PEERS, ",")
	} else {
		peers = nil
	}

	node := pow.NewNode(
		*HOST,
		*WSPORT,
		pow.DIFFICULTY,
		*N_MINERS,
		*validators,
		*blockchain,
		*wallet,
		*txPool,
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
