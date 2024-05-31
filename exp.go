package main

import (
	"consensus-algorithms-with-golang/exp"
	"flag"
	"strings"
)

func main() {
	HOST := flag.String("HOST", "localhost", "Host to bind to")
	WSPORT := flag.Int("WSPORT", 8080, "Port to bind to")
	PEERS := flag.String("PEERS", "", "Comma separated list of PEERS")

	flag.Parse()
	var peers []string
	if *PEERS != "" {
		peers = strings.Split(*PEERS, ",")
	} else {
		peers = nil
	}

	server := exp.NewServer(*HOST, *WSPORT)
	server.Listen(peers)

	// handle OS interruption
	exp.HandleOSSignals()
	for _, conn := range server.Peers {
		conn.Close()
	}
}
