package main

import "consensus-algorithms-with-golang/exp"

func main() {
	host := "localhost:8080"
	server := exp.NewServer(host)
	server.Listen()
}
