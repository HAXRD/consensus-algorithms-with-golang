package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// queryNodeInfoHandler queries the sockets of the node
func (node *Node) queryNodeInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	str := "Node info:\n"
	// validators
	str += "[Validators]\n"
	for i, pubKey := range node.Validators.list {
		str += fmt.Sprintf("%d: %v\n", i, chain_util2.FormatHash(pubKey))
	}
	// blockchain
	str += "[Blockchain]\n"
	for i, block := range node.Blockchain.chain {
		str += fmt.Sprintf("[%v]", block.Hash)
		if i < len(node.Blockchain.chain)-1 {
			str += "-->"
		} else {
			str += "\n"
		}
	}
	// available sockets
	str += "[Sockets]\n"
	for key := range node.Sockets {
		conn, ok := node.Sockets[key]
		str += fmt.Sprintf("%s: %v\n", key, conn != nil && ok)
	}
	// TxPool
	str += "[TxPool]\n"
	for i, tx := range node.TxPool.pool {
		str += fmt.Sprintf("%d: %v by %v\n", i, chain_util2.FormatHash(tx.Hash), chain_util2.FormatHash(tx.From))
	}
	// BlockPool
	str += "[BlockPool]\n"
	for i, block := range node.BlockPool.pool {
		str += fmt.Sprintf("%d: %v by %v\n", i, chain_util2.FormatHash(block.Hash), chain_util2.FormatHash(chain_util2.Key2Str(block.Proposer)))
	}
	// PreparePool
	str += "[PreparePool]\n"
	for bh, msgs := range node.PreparePool.mapPool {
		str += fmt.Sprintf("--> %v\n", chain_util2.FormatHash(bh))
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %v\n", i, chain_util2.FormatHash(chain_util2.Key2Str(msg.PublicKey)))
		}
	}
	// CommitPool
	str += "[CommitPool]\n"
	for bh, msgs := range node.CommitPool.mapPool {
		str += fmt.Sprintf("--> %v\n", chain_util2.FormatHash(bh))
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %v\n", i, chain_util2.FormatHash(chain_util2.Key2Str(msg.PublicKey)))
		}
	}
	// RCPool
	str += "[RCPool]\n"
	for bh, msgs := range node.RCPool.mapPool {
		str += fmt.Sprintf("--> %v\n", chain_util2.FormatHash(bh))
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %v\n", i, chain_util2.FormatHash(chain_util2.Key2Str(msg.PublicKey)))
		}
	}
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write to HTTP client failed, [%s], skip this one!\n", err)
	}
}

// broadcast broadcasts the given message to all sockets
func (node *Node) broadcast(msg string) {
	time.Sleep(5 * time.Second) // TODO: change this to random period of time
	mutex.Lock()
	defer mutex.Unlock()
	for url, conn := range node.Sockets {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("Error broadcasting message [%s] to [%s], %v", msg, url, err)
			conn.Close()
			delete(node.Sockets, url)
		}
	}
}

// wsServerHandler is the websocket server handler that holds the main PBFT logic
func (node *Node) wsServerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade to websocket failed, [%s]\n", err)
		return
	}
	log.Printf("Remote address [%s] connected!\n", r.RemoteAddr)
	defer conn.Close()

	// add incoming connections to sockets
	remoteUrl := fmt.Sprintf("ws://%s/ws", conn.RemoteAddr().String())
	mutex.Lock()
	node.Sockets[remoteUrl] = conn
	mutex.Unlock()

	// handle any requests by broadcasting it
	// TODO[broadcast storm]: this should result in a broadcast storm, remove it later
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read message failed, [%s]\n", err)
			break
		}
		log.Printf("recv: [%s]\n", msg)

		node.broadcast(string(msg))
	}
}

// makeTestCallHandler is a test call
// TODO: remove this later [broadcast storm]
func (node *Node) makeTestCallHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	msg := "Hello!"

	// wait for WsClient(Relay) up-online
	for {
		time.Sleep(1 * time.Second)
		if node.Relay != nil {
			break
		}
	}
	// Write to web page
	w.Write([]byte(msg))
	// Write to WsClient(Relay)
	err := node.Relay.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Write message failed, [%s]\n", err)
	}
}

// makeTxHandler makes a tx on current node
func (node *Node) makeTxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	tx := node.Wallet.CreateTx(time.Now().String() + " " + "this is a test message")
	msg, err := json.Marshal(tx)
	if err != nil {
		log.Printf("Marshal tx failed, [%s]\n", err)
		return
	}

	// write message
	for {
		if node.Relay != nil {
			break
		}
	}
	// Write to web page
	w.Write([]byte(msg))
	// Write to WsClient(Relay)
	err = node.Relay.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Write message failed, [%s]\n", err)
	}
}

//1. makeTxHandler
//3. queryBlockPoolHandler
//4. queryPreparePoolHandler
//5. queryCommitPoolHandler
//6. queryRCPoolHandler
//7. queryBlockchainHandler
