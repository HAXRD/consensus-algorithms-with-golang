package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"encoding/json"
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
}
var mutex = &sync.Mutex{}

/**
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

// broadcast broadcasts passed-in msg to node's sockets
func (node *Node) broadcast(msg string) {
	mutex.Lock()
	defer mutex.Unlock()
	for url, conn := range node.Sockets {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("Error broadcasting msg [%s] to [%s], %v", msg, url, err)
			conn.Close()
			delete(node.Sockets, url)
		}
	}
}

// wsMsgHandler is a websocket handler that handles any requests from peer nodes.
func (node *Node) wsMsgHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Remote address [%s] connected!", r.RemoteAddr)
	// upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error, %s\n", err)
		return
	}
	defer conn.Close()

	// add incoming connection to peers
	remoteUrl := fmt.Sprintf("ws://%s/ws", conn.RemoteAddr().String())
	mutex.Lock()
	node.Sockets[remoteUrl] = conn
	mutex.Unlock()

	// handle requests infinitely
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v\n", err)
			} else {
				log.Printf("Websocket closed gracefully: %v\n", err)
			}
			log.Printf("Remote address [%s] disconnected!", r.RemoteAddr)
			break
		}
		log.Printf("recv: %s\n", msg)

		// parse msg to different types and perform different ops
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Printf("Unmarshal msg failed, %s, skip this one!\n", err)
			continue
		}

		if msgTypeRaw, ok := data["msgType"]; ok {
			if msgType, ok2 := msgTypeRaw.(string); ok2 {

				switch msgType {
				case MsgTx:
					var tx Transaction
					if err := json.Unmarshal(msg, &tx); err != nil {
						log.Printf("Unmarshal msg->tx failed, %s, skip this one!\n", err)
						continue
					}
					// check if tx is valid
					if !node.TxPool.TxExists(tx) &&
						node.TxPool.VerifyTx(tx) &&
						node.Validators.ValidatorExists(tx.From) {
						// add tx to tx pool
						thresholdReached := node.TxPool.AddTx2Pool(tx)
						// broadcast
						node.broadcast(string(msg))

						if thresholdReached {
							log.Println("THRESHOLD REACHED!")
							if node.Blockchain.GetProposer() == chain_util2.Key2Str(node.Wallet.publicKey) {
								log.Println("PROPOSING A NEW BLOCK!")
								block := node.Blockchain.CreateBlock(node.Wallet, node.TxPool.pool)
								// TODO: broadcast
								newMsg, err := json.Marshal(block)
								if err != nil {
									log.Printf("Marshal block failed, %s, msg won't be sent, skip this one!\n", err)
									continue
								}
								node.broadcast(string(newMsg))
							}
						}
					}
				case MsgPrePrepare:
					var block Block
					if err := json.Unmarshal(msg, &block); err != nil {
						log.Printf("Unmarshal msg->block failed, %s, skip this one!\n", err)
						continue
					}
					// check if block is valid
					if exists, _ := node.BlockPool.BlockExists(block.Hash); !exists && node.Blockchain.VerifyBlock(block) {
						// add block to block pool
						node.BlockPool.AddBlock2Pool(block)
						// broadcast
						node.broadcast(string(msg))

						// create prepareMsg and broadcast it
						prepareMsg := node.Wallet.CreateMsg(MsgPrepare, block.Hash)
						newMsg, err := json.Marshal(prepareMsg)
						if err != nil {
							log.Printf("Marshal prepareMsg failed, %s, msg won't be sent, skip this one!\n", err)
							continue
						}
						node.broadcast(string(newMsg))
					}
				case MsgPrepare:
					var prepareMsg Message
					if err := json.Unmarshal(msg, &prepareMsg); err != nil {
						log.Printf("Unmarshal msg->prepareMsg failed, %s, skip this one!\n", err)
						continue
					}
					// check if prepareMsg is valid
					if !node.PreparePool.MsgExists(prepareMsg) &&
						node.PreparePool.VerifyMsg(prepareMsg) &&
						node.Validators.ValidatorExists(chain_util2.Key2Str(prepareMsg.PublicKey)) {
						// add prepareMsg to prepare pool
						node.PreparePool.AddMsg2Pool(prepareMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.PreparePool.mapPool[prepareMsg.BlockHash]) >= MIN_APPROVALS {
							// create commitMsg and broadcast it
							commitMsg := node.Wallet.CreateMsg(MsgCommit, prepareMsg.BlockHash)
							newMsg, err := json.Marshal(commitMsg)
							if err != nil {
								log.Printf("Marshal commitMsg failed, %s, msg won't be sent, skip this one!\n", err)
								continue
							}
							node.broadcast(string(newMsg))
						}
					}
				case MsgCommit:
					var commitMsg Message
					if err := json.Unmarshal(msg, &commitMsg); err != nil {
						log.Printf("Unmarshal msg->commitMsg failed, %s, skip this one!\n", err)
						continue
					}
					// check if commitMsg is valid
					if !node.CommitPool.MsgExists(commitMsg) &&
						node.CommitPool.VerifyMsg(commitMsg) &&
						node.Validators.ValidatorExists(chain_util2.Key2Str(commitMsg.PublicKey)) {
						// add commitMsg to commit pool
						node.CommitPool.AddMsg2Pool(commitMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.CommitPool.mapPool[commitMsg.BlockHash]) >= MIN_APPROVALS {
							// add block to the chain
							node.Blockchain.AddUpdatedBlock2Chain(commitMsg.BlockHash, node.BlockPool, node.PreparePool, node.CommitPool)
						}
						// create rcMsg and broadcast it
						rcMsg := node.Wallet.CreateMsg(MsgRC, commitMsg.BlockHash)
						newMsg, err := json.Marshal(rcMsg)
						if err != nil {
							log.Printf("Marshal rcMsg failed, %s, msg won't be sent, skip this one!\n", err)
							continue
						}
						node.broadcast(string(newMsg))
					}
				case MsgRC:
					var rcMsg Message
					if err := json.Unmarshal(msg, &rcMsg); err != nil {
						log.Printf("Unmarshal msg->rcMsg failed, %s, skip this one!\n", err)
						continue
					}
					// check if rcMsg is valid
					if !node.RCPool.MsgExists(rcMsg) &&
						node.RCPool.VerifyMsg(rcMsg) &&
						node.Validators.ValidatorExists(chain_util2.Key2Str(rcMsg.PublicKey)) {
						// add rcMsg to rc pool
						node.RCPool.AddMsg2Pool(rcMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.RCPool.mapPool[rcMsg.BlockHash]) >= MIN_APPROVALS {
							// TODO: implement clean mechanism
							log.Println("[REACHED RC!!!!!]")
						}
					}
				default:
					log.Println("[default] unknown msgType!")
				}
			} else {
				log.Println("[inner] unknown msgType!")
			}
		} else {
			log.Println("[outer] unknown msgType!")
		}
	}
}

// connectPeers connects peers
func (node *Node) connectPeers(peerAddrs []string) {
	for _, addr := range peerAddrs {
		url := fmt.Sprintf("ws://%s/ws", addr)
		log.Printf("Trying connecting to peer [%s]...\n", url)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Error connecting to peer [%s]: %s\n", url, err)
			continue
		}
		log.Printf("Connected to peer [%s]\n", url)

		mutex.Lock()
		node.Sockets[url] = conn
		mutex.Unlock()
	}
}

// Listen launches the websocket server and http endpoints,
// then connects to peers.
func (node *Node) Listen(peerAddrs []string) {
	// http endpoint
	http.HandleFunc("/makeTx", node.makeTxHandler)
	go func() {
		url := chain_util2.FormatUrl(node.Host, node.Port)
		log.Printf("Http server listening on [%s]...\n", url)
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatalf("Http server listen failed, %s\n", err)
		}
	}()

	// websocket server
	http.HandleFunc("/ws", node.wsMsgHandler)
	go func() {
		url := chain_util2.FormatUrl(node.Host, node.WsPort)
		log.Printf("Websocket server listening on [%s]...\n", url)
		err := http.ListenAndServe(url, nil)
		if err != nil {
			log.Fatalf("Websocket server listen failed, %s\n", err)
		}
	}()

	/* Add connections to sockets */
	// use websocket client to connect to this node's websocket server
	for {
		wsUrl := fmt.Sprintf("ws://%s/ws", chain_util2.FormatUrl(node.Host, node.WsPort))
		log.Printf("Trying connecting to itself [%s]...\n", wsUrl)
		wsClientConn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			time.Sleep(1 * time.Second)
			log.Printf("Still trying connecting to itself [%s]...\n", wsUrl)
			continue
		}
		log.Printf("Connected to itself [%s]\n", wsUrl)
		mutex.Lock()
		node.Sockets[wsUrl] = wsClientConn
		mutex.Unlock()
		break
	}
	// connect to peers
	if len(peerAddrs) > 0 {
		node.connectPeers(peerAddrs)
	}
}

// ========= Below are HTTP handlers ==========

// makeTxHandler make a tx
func (node *Node) makeTxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	tx := node.Wallet.CreateTx(time.Now().String() + "This is a test tx")
	msg, err := json.Marshal(tx)
	if err != nil {
		log.Printf("Marshal tx failed, [%s], skip this one!\n", err)
		return
	}
	_, err = fmt.Fprintf(w, string(msg))
	if err != nil {
		log.Printf("Write to HTTP client failed, [%s], skip this one!\n", err)
		return
	}
	url := fmt.Sprintf("ws://%s/ws", chain_util2.FormatUrl(node.Host, node.WsPort))
	wsClientConn := node.Sockets[url]
	err = wsClientConn.WriteMessage(websocket.TextMessage, []byte(msg))
	log.Printf("Writing msg to [%s]\n", url)
	if err != nil {
		log.Printf("Error broadcasting msg [%s] to [%s], %v", msg, url, err)
		wsClientConn.Close()
		mutex.Lock()
		delete(node.Sockets, url)
		mutex.Unlock()
	}
}
