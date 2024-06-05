package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
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

// Listen launches the http endpoints, websocket server, websocket client
// and then connects to peers.
func (node *Node) Listen() {
	// http endpoints
	http.HandleFunc("/queryNodeInfo", node.queryNodeInfoHandler)
	go func() {
		url := chain_util2.FormatUrl(node.Host, node.Port)
		log.Printf("Http server listening on [%s]...\n", url)
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatalf("Http server listen failed, %s\n", err)
		}
	}()

}
