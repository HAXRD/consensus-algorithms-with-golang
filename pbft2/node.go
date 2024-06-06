package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var mutex = &sync.Mutex{}

/*
*
A Node represents a single node in a blockchain system.
It has the following properties:
- Host: host address of the node
- WsPort: websocket port
- Port: http server port
- Sockets: the addresses of itself and all connected peers
- Validators: system acknowledged public keys (wallet/node)
- Blockchain: a copy of the blockchain
- Wallet: node's wallet
- TxPool: node's tx pool
- BlockPool: node's block pool
- PreparePool: node's prepare pool
- CommitPool: node's commit pool
- RCPool: node's round-change pool

It features the following methods:
1. NewNode
2. broadcast
3. wsMsgHandler
4. connectPeers
5. Listen
=======below are http handlers=============
1. makeTxHandler
2. queryTxPoolHandler
3. queryBlockPoolHandler
4. queryPreparePoolHandler
5. queryCommitPoolHandler
6. queryRCPoolHandler
7. queryBlockchainHandler
*/

type Node struct {
	Host        string
	WsPort      uint64
	Port        uint64
	Sockets     map[string]*websocket.Conn
	Relay       *websocket.Conn
	Validators  Validators
	Blockchain  Blockchain
	Wallet      Wallet
	TxPool      TransactionPool
	BlockPool   BlockPool
	PreparePool MsgPool
	CommitPool  MsgPool
	RCPool      MsgPool
}

// NewNode creates a new node with given info
func NewNode(host string, wsPort uint64, vs Validators, bc Blockchain, w Wallet,
	tp TransactionPool, bp BlockPool, pp MsgPool, cp MsgPool, rcp MsgPool) *Node {
	return &Node{
		Host:        host,
		WsPort:      wsPort,
		Port:        wsPort + 10000,
		Sockets:     make(map[string]*websocket.Conn),
		Relay:       nil,
		Validators:  vs,
		Blockchain:  bc,
		Wallet:      w,
		TxPool:      tp,
		BlockPool:   bp,
		RCPool:      rcp,
		PreparePool: pp,
		CommitPool:  cp,
	}
}

// launchHttpServer launches the Http endpoint server
func (node *Node) launchHttpServer() {
	httpUrl := chain_util2.FormatUrl(node.Host, node.Port)
	log.Printf("Http server listening on [%s]...\n", httpUrl)
	err := http.ListenAndServe(httpUrl, nil)
	if err != nil {
		log.Fatalf("Http server listen failed, %s\n", err)
	}
}

// launchWsServer launches the websocket server
func (node *Node) launchWsServer() {
	wsServerUrl := chain_util2.FormatUrl(node.Host, node.WsPort)
	log.Printf("Websocket server listening on [%s]...\n", wsServerUrl)
	err := http.ListenAndServe(wsServerUrl, nil)
	if err != nil {
		log.Fatalf("Websocket server listen failed, %s\n", err)
	}
}

// launchWsClient launches the websocket client to connection to
// current node's ws server
// The message source of WsClient(Relay) could be
//
//	CASE 1. Direct `WriteMessage()` calls from any http endpoints requests
//	CASE 2. Direct `WriteMessage()` calls from peer ws socket connections
//	   after `peer.ReadMessage()` then calling `Relay.WriteMessage()`
//	CASE 3. Relaying by `Relay.ReadMessage()` then `Relay.WriteMessage()`
func (node *Node) launchWsClient() {
	// keep trying dialing to ws server util ws server is online
	var wsClientConn *websocket.Conn
	var err error
	for {
		wsServerUrl := fmt.Sprintf("ws://%s/ws", chain_util2.FormatUrl(node.Host, node.WsPort))
		log.Printf("Dialing itself [%s]...\n", wsServerUrl)
		wsClientConn, _, err = websocket.DefaultDialer.Dial(wsServerUrl, nil)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Connected to itself [%s]...\n", wsServerUrl)

		node.Relay = wsClientConn
		break
	}

	// relaying any incoming messages to ws server
	// (this is CASE 3)
	var mt int
	var relayMsg []byte
	for {
		mt, relayMsg, err = node.Relay.ReadMessage()
		if err != nil {
			log.Printf("Relay error reading msg, [%s], skipping...\n", err)
		}
		err = node.Relay.WriteMessage(mt, relayMsg)
		if err != nil {
			log.Printf("Relay error writing msg, [%s], skipping...\n", err)
		}
	}
}

func (node *Node) launchPeer(peerUrl string) {
	var mt int
	var msg []byte
	var err error
	// wait for Relay up-online
	for {
		time.Sleep(1 * time.Second)
		if node.Relay != nil {
			break
		}
	}
	// relaying any received message to current node's WsClient(Relay)
	for {
		mt, msg, err = node.Sockets[peerUrl].ReadMessage()
		if err != nil {
			log.Printf("Error reading from peer [%s], %v\n", peerUrl, err)
		}
		err = node.Relay.WriteMessage(mt, msg)
		if err != nil {
			log.Printf("Error relaying to WsClient/Relay, %v\n", err)
		}
	}
}

// connectPeers connects current node's peers
func (node *Node) connectPeers(peers []string) {
	for _, peer := range peers {
		peerUrl := fmt.Sprintf("ws://%s/ws", peer)
		log.Printf("Dialing peer [%s]...\n", peerUrl)
		peerConn, _, err := websocket.DefaultDialer.Dial(peerUrl, nil)
		if err != nil {
			log.Printf("Error connecting to peer [%s], skipping...\n", err)
			continue
		}
		log.Printf("Connected to peer [%s]...\n", peerUrl)

		mutex.Lock()
		node.Sockets[peerUrl] = peerConn
		mutex.Unlock()

		// launchPeer
		go node.launchPeer(peerUrl)
	}
}

// Listen launches the http endpoints, websocket server, websocket client
// and then connects to peers.
func (node *Node) Listen(peers []string) {
	// http endpoints
	http.HandleFunc("/queryNodeInfo", node.queryNodeInfoHandler)
	http.HandleFunc("/makeTx", node.makeTxHandler)
	go node.launchHttpServer()

	// websocket server
	http.HandleFunc("/ws", node.wsServerHandler)
	go node.launchWsServer()

	// websocket client
	go node.launchWsClient()

	// peers
	node.connectPeers(peers)
}
