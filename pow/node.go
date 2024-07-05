package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
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

type Node struct {
	Host        string
	WsPort      uint64
	Port        uint64
	Sockets     map[string]*websocket.Conn
	Relay       *websocket.Conn
	Difficulty  uint64
	NumOfMiners uint64
	Miners      []string
	Validators  Validators
	Blockchain  Blockchain
	Wallet      Wallet
	TxPool      TxPool

	fetchChan       chan []Transaction
	doneChan        chan struct{}
	addMinerChan    chan struct{}
	removeMinerChan chan struct{}
	wg              sync.WaitGroup
}

func NewNode(host string, wsPort uint64, difficulty uint64, numOfMiners uint64,
	vs Validators, bc Blockchain, w Wallet, txp TxPool) *Node {
	return &Node{
		Host:        host,
		WsPort:      wsPort,
		Port:        wsPort + 10000,
		Sockets:     make(map[string]*websocket.Conn),
		Relay:       nil,
		Difficulty:  difficulty,
		NumOfMiners: numOfMiners,
		Validators:  vs,
		Blockchain:  bc,
		Wallet:      w,
		TxPool:      txp,

		fetchChan:       make(chan []Transaction, 1),
		doneChan:        make(chan struct{}),
		addMinerChan:    make(chan struct{}),
		removeMinerChan: make(chan struct{}),
	}
}

func (node *Node) Listen(peers []string) {
	// http endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/queryNodeInfo", node.queryNodeInfoHandler)
	mux.HandleFunc("/setDifficulty", node.setDifficultyHandler)
	//mux.HandleFunc("/addMiner", node.addMinerHandler)
	//mux.HandleFunc("/removeMiner", node.removeMinerHandler)
	mux.HandleFunc("/makeTx", node.makeTxHandler)
	//mux.HandleFunc("/reset", node.resetHandler)
	go node.launchHttpServer(mux)

	// websocket server
	http.HandleFunc("/ws", node.wsServerHandler)
	go node.launchWsServer()

	// websocket client
	go node.launchWsClient()

	// peers
	node.connectPeers(peers)

	// launch periodicallyFetch goroutine
	go node.periodicallyFetch()

	// launch miners
	node.launchMiners()
}

func (node *Node) launchHttpServer(mux *http.ServeMux) {
	httpUrl := pow_util.FormatUrl(node.Host, node.Port)
	err := http.ListenAndServe(httpUrl, corsMiddleware(mux))
	if err != nil {
		log.Fatalf("Http server listening failed, %s\n", err)
	}
	log.Printf("Http server listening on [%s]...\n", httpUrl)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (node *Node) launchWsServer() {
	wsServerUrl := pow_util.FormatUrl(node.Host, node.WsPort)
	err := http.ListenAndServe(wsServerUrl, nil)
	if err != nil {
		log.Fatalf("Websocket server listening failed, %s\n", err)
	}
	log.Printf("Websocket server listening on [%s]...\n", wsServerUrl)
}

func (node *Node) launchWsClient() {
	// keep trying dialing to ws server util ws server is online
	var wsClientConn *websocket.Conn
	var err error
	for {
		wsServerUrl := fmt.Sprintf("ws://%s/ws", pow_util.FormatUrl(node.Host, node.WsPort))
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
	var mt int
	var relayMsg []byte
	for {
		mt, relayMsg, err = node.Relay.ReadMessage()
		if err != nil {
			log.Printf("Relay error reading msg, [%s], skipping...\n", err)
		}
		mutex.Lock()
		err = node.Relay.WriteMessage(mt, relayMsg)
		mutex.Unlock()
		if err != nil {
			log.Printf("Relay error writing msg, [%s], skipping...\n", err)
		}
	}
}

func (node *Node) connectPeers(peers []string) {
	for _, peer := range peers {
		peerUrl := fmt.Sprintf("ws://%s/ws", peer)
		log.Printf("Dialing peer [%s]...\n", peerUrl)
		peerConn, _, err := websocket.DefaultDialer.Dial(peerUrl, nil)
		if err != nil {
			log.Printf("Peer [%s] websocket dialing failed, %s\n", peer, err)
			continue
		}
		log.Printf("Connected to peer [%s]...\n", peer)

		mutex.Lock()
		node.Sockets[peerUrl] = peerConn
		mutex.Unlock()

		go node.launchPeer(peerUrl)
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
	// relaying any received messages to current node's WsClient(Relay)
	for {
		mt, msg, err = node.Sockets[peerUrl].ReadMessage()
		if err != nil {
			log.Printf("Error reading from peer [%s], %v\n", peerUrl, err)
		}
		mutex.Lock()
		err = node.Relay.WriteMessage(mt, msg)
		mutex.Unlock()
		if err != nil {
			log.Printf("Error relaying to WsClient/Relay, %v\n", err)
		}
	}
}

func (node *Node) periodicallyFetch() {
	for {
		time.Sleep(UPDATE_INTERVAL * time.Second)

		// only periodicallyFetch when `waiting` is not empty
		numOfTxs := len(node.TxPool.waiting)
		mutex.Lock()
		if numOfTxs > 0 {
			// convert map into slice
			waitingSlice := make([]Transaction, 0, numOfTxs)
			for _, tx := range node.TxPool.waiting {
				waitingSlice = append(waitingSlice, tx)
			}
			node.fetchChan <- waitingSlice
			log.Printf("Fetched [%d] txs from waiting", numOfTxs)
		} else {
			log.Printf("Waiting pool is empty, nothing to fetch")
		}
		mutex.Unlock()
	}
}

func (node *Node) miner(id int) {
	defer node.wg.Done()
	for {
		select {
		case fetchedTxs := <-node.fetchChan:
			// timeout go routine
			timeoutChan := make(chan struct{})
			go func() {
				time.Sleep((UPDATE_INTERVAL - 1) * time.Second)
				close(timeoutChan)
			}()

			block := node.tryPropose(timeoutChan, fetchedTxs)
			if block != nil {
				// Write to WsClient(Relay)
				msgStr, err := WrapData2MsgStr(*block)
				if err != nil {
					log.Printf("Marshal block failed, [%s]\n", err)
					return
				}
				mutex.Lock()
				err = node.Relay.WriteMessage(websocket.TextMessage, msgStr)
				log.Printf("Miner [%d] proposed a block [%s]!\n", id, pow_util.Byte2Hex(block.Hash)[:node.Difficulty+4])
				mutex.Unlock()
			} else {
				log.Printf("Miner [%d] timeout reached!\n", id)
			}
		case <-node.doneChan:
			fmt.Printf("Miner [%d] stopped\n", id)
		}
	}
}

func (node *Node) tryPropose(timeoutChan chan struct{}, fetchedTxs []Transaction) *Block {
	for {
		select {
		case <-timeoutChan:
			return nil
		default:
			// generate a random nonce
			nonce := rand.Uint64()
			lastHash := node.Blockchain.chain[len(node.Blockchain.chain)-1].Hash
			hash := HashBlock(lastHash, fetchedTxs, nonce)
			// check if meets difficulty requirement
			if pow_util.HashMeetsDifficulty(hash, node.Difficulty) {
				// propose a block
				timestamp := time.Now().String()
				proposer := node.Wallet.pubKey
				signature := node.Wallet.Sign(hash)
				block := node.Wallet.CreateBlock(timestamp, lastHash, hash, fetchedTxs, proposer, signature, nonce)
				return block
			}
		}
	}
}

func (node *Node) launchMiners() {
	for i := 0; i < int(node.NumOfMiners); i++ {
		node.wg.Add(1)
		go node.miner(i + 1)
		log.Printf("Miner [%d] online...\n", i+1)
	}
}
