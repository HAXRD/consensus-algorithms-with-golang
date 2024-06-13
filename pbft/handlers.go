package pbft

import (
	"consensus-algorithms-with-golang/pbft/chain_util"
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
	str := fmt.Sprintf("Node[%s] Info:\n", chain_util.BytesToHex(node.Wallet.publicKey)[:6])
	// validators
	str += "\n[Validators]\n"
	for i, pubKey := range node.Validators.list {
		str += fmt.Sprintf("%d: %s\n", i, chain_util.BytesToHex(pubKey)[:6])
	}
	// blockchain
	str += "\n[Blockchain]\n"
	for i, block := range node.Blockchain.chain {
		str += fmt.Sprintf("[%s]", chain_util.BytesToHex(block.Hash)[:6])
		if i < len(node.Blockchain.chain)-1 {
			str += "-->"
		} else {
			str += "\n"
		}
	}
	// available sockets
	str += "\n[Sockets]\n"
	for key := range node.Sockets {
		conn, ok := node.Sockets[key]
		str += fmt.Sprintf("%s: %v\n", key, conn != nil && ok)
	}
	// TxPool
	str += "\n[TxPool]\n"
	for i, tx := range node.TxPool.pool {
		str += fmt.Sprintf("%d: %s by %s\n", i, chain_util.BytesToHex(tx.Hash)[:6], chain_util.BytesToHex(tx.From)[:6])
	}
	// BlockPool
	str += "\n[BlockPool]\n"
	for i, block := range node.BlockPool.pool {
		str += fmt.Sprintf("%d: %s by %s\n", i, chain_util.BytesToHex(block.Hash)[:6], chain_util.BytesToHex(block.Proposer)[:6])
	}
	// PreparePool
	str += "\n[PreparePool]\n"
	for bh, msgs := range node.PreparePool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:5])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
	}
	// CommitPool
	str += "\n[CommitPool]\n"
	for bh, msgs := range node.CommitPool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:5])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
	}
	// RCPool
	str += "\n[RCPool]\n"
	for bh, msgs := range node.RCPool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:6])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
	}
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Printf("Write to HTTP client failed, [%v], skip this one!\n", err)
	}
}

// broadcast broadcasts the given message to all sockets
func (node *Node) broadcast(msg string) {
	time.Sleep(3 * time.Second) // TODO: change this to random period of time
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
		log.Printf("Upgrade to websocket failed, %v\n", err)
		return
	}
	log.Printf("Remote address [%s] connected!\n", r.RemoteAddr)
	defer conn.Close()

	// add incoming connections to sockets
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
		//log.Printf("recv: %s\n", msg)

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
						mutex.Lock()
						poolCopy, success := node.TxPool.AddTx2Pool(tx)
						mutex.Unlock()
						if !success {
							continue
						}
						// broadcast
						node.broadcast(string(msg))

						if poolCopy != nil {
							log.Println("THRESHOLD REACHED!")
							if chain_util.BytesToHex(node.Blockchain.GetProposer()) == chain_util.BytesToHex(node.Wallet.publicKey) {
								log.Println("PROPOSING A NEW BLOCK!")
								block := node.Blockchain.CreateBlock(node.Wallet, poolCopy)

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
						mutex.Lock()
						success := node.BlockPool.AddBlock2Pool(block)
						mutex.Unlock()
						if !success {
							continue
						}
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
						node.Validators.ValidatorExists(prepareMsg.PublicKey) {
						// add prepareMsg to prepare pool
						mutex.Lock()
						success := node.PreparePool.AddMsg2Pool(prepareMsg)
						mutex.Unlock()
						if !success {
							continue
						}
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.PreparePool.mapPool[chain_util.BytesToHex(prepareMsg.BlockHash)]) >= MIN_APPROVALS {
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
						node.Validators.ValidatorExists(commitMsg.PublicKey) {
						// add commitMsg to commit pool
						mutex.Lock()
						success := node.CommitPool.AddMsg2Pool(commitMsg)
						mutex.Unlock()
						if !success {
							continue
						}
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.CommitPool.mapPool[chain_util.BytesToHex(commitMsg.BlockHash)]) >= MIN_APPROVALS {
							// add block to the chain
							mutex.Lock()
							node.Blockchain.AddUpdatedBlock2Chain(commitMsg.BlockHash, node.BlockPool, node.PreparePool, node.CommitPool)
							mutex.Unlock()
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
						node.Validators.ValidatorExists(rcMsg.PublicKey) {
						// add rcMsg to rc pool
						mutex.Lock()
						success := node.RCPool.AddMsg2Pool(rcMsg)
						mutex.Unlock()
						if !success {
							continue
						}
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.RCPool.mapPool[chain_util.BytesToHex(rcMsg.BlockHash)]) >= MIN_APPROVALS {
							// TODO: implement clean mechanism
							log.Println("[REACHED RC!!!!!]")
							mutex.Lock()
							success := node.TxPool.TransferInProgressToCommitted(node.BlockPool.GetBlock(rcMsg.BlockHash).Data)
							mutex.Unlock()
							if success {
								log.Println("[TRANSFERRED IN-PROGRESS TO COMMITTED SUCCESSFULLY!!!]")
							} else {
								log.Println("[TRANSFERRED IN-PROGRESS TO COMMITTED FAILED!!!]")
							}
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
	mutex.Lock()
	err := node.Relay.WriteMessage(websocket.TextMessage, []byte(msg))
	mutex.Unlock()
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
	mutex.Lock()
	err = node.Relay.WriteMessage(websocket.TextMessage, []byte(msg))
	mutex.Unlock()
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

type BlockInfo struct {
	Hash     string `json:"hash"`
	Proposer string `json:"proposer"`
}

type TxPoolItem struct {
	Hash   string `json:"hash"`
	PubKey string `json:"pubKey"`
}

type TxPoolInfo struct {
	Waiting    []TxPoolItem `json:"waiting"`
	InProgress []TxPoolItem `json:"inProgress"`
	Committed  []TxPoolItem `json:"committed"`
}

type BlockPoolItem struct {
	BlockHash string `json:"blockHash"`
	PubKey    string `json:"pubKey"`
}

type MsgPoolItem struct {
	BlockHash string   `json:"blockHash"`
	FromWhos  []string `json:"fromWhos"`
}

// Data defines the data object sent to frontend server for display
type Data struct {
	NodeHash    string          `json:"nodeHash"`
	BlockChain  []BlockInfo     `json:"blockchain"`
	Sockets     []string        `json:"sockets"`
	TxPool      TxPoolInfo      `json:"txPool"`
	BlockPool   []BlockPoolItem `json:"blockPool"`
	PreparePool []MsgPoolItem   `json:"preparePool"`
	CommitPool  []MsgPoolItem   `json:"commitPool"`
	RCPool      []MsgPoolItem   `json:"rcPool"`
}

func (node *Node) queryNodeInfo2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nodeHash := chain_util.BytesToHex(node.Wallet.publicKey)[:6]
	blockChain := make([]BlockInfo, 0, len(node.Blockchain.chain))
	for _, block := range node.Blockchain.chain {
		blockChain = append(blockChain, BlockInfo{
			Hash:     chain_util.BytesToHex(block.Hash)[:6],
			Proposer: chain_util.BytesToHex(block.Proposer)[:6],
		})
	}
	sockets := make([]string, 0, len(node.Sockets))
	for key := range node.Sockets {
		sockets = append(sockets, key)
	}

	txPool := TxPoolInfo{
		Waiting:    make([]TxPoolItem, len(node.TxPool.pool)),
		InProgress: make([]TxPoolItem, 0, len(node.TxPool.inProgress)),
		Committed:  make([]TxPoolItem, 0, len(node.TxPool.committed)),
	}
	for i, transaction := range node.TxPool.pool {
		txPool.Waiting[i] = TxPoolItem{
			Hash:   chain_util.BytesToHex(transaction.Hash)[:6],
			PubKey: chain_util.BytesToHex(transaction.From)[:6],
		}
	}
	for _, transaction := range node.TxPool.inProgress {
		txPool.InProgress = append(txPool.InProgress, TxPoolItem{
			Hash:   chain_util.BytesToHex(transaction.Hash)[:6],
			PubKey: chain_util.BytesToHex(transaction.From)[:6],
		})
	}
	for _, transaction := range node.TxPool.committed {
		txPool.Committed = append(txPool.Committed, TxPoolItem{
			Hash:   chain_util.BytesToHex(transaction.Hash)[:6],
			PubKey: chain_util.BytesToHex(transaction.From)[:6],
		})
	}

	blockPool := make([]BlockPoolItem, 0, len(node.BlockPool.pool))
	for _, block := range node.BlockPool.pool {
		blockPool = append(blockPool, BlockPoolItem{
			BlockHash: chain_util.BytesToHex(block.Hash)[:6],
			PubKey:    chain_util.BytesToHex(block.Proposer)[:6],
		})
	}
	preparePool := make([]MsgPoolItem, 0, len(node.PreparePool.mapPool))
	for blockHash, msgs := range node.PreparePool.mapPool {
		fromWhos := make([]string, 0, len(msgs))
		for _, msg := range msgs {
			fromWhos = append(fromWhos, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
		preparePool = append(preparePool, MsgPoolItem{
			BlockHash: blockHash[:6],
			FromWhos:  fromWhos,
		})
	}
	commitPool := make([]MsgPoolItem, 0, len(node.CommitPool.mapPool))
	for blockHash, msgs := range node.CommitPool.mapPool {
		fromWhos := make([]string, 0, len(msgs))
		for _, msg := range msgs {
			fromWhos = append(fromWhos, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
		commitPool = append(commitPool, MsgPoolItem{
			BlockHash: blockHash[:6],
			FromWhos:  fromWhos,
		})
	}
	rcPool := make([]MsgPoolItem, 0, len(node.RCPool.mapPool))
	for blockHash, msgs := range node.RCPool.mapPool {
		fromWhos := make([]string, 0, len(msgs))
		for _, msg := range msgs {
			fromWhos = append(fromWhos, chain_util.BytesToHex(msg.PublicKey)[:6])
		}
		rcPool = append(rcPool, MsgPoolItem{
			BlockHash: blockHash[:6],
			FromWhos:  fromWhos,
		})
	}

	data := Data{
		NodeHash:    nodeHash,
		BlockChain:  blockChain,
		Sockets:     sockets,
		TxPool:      txPool,
		BlockPool:   blockPool,
		PreparePool: preparePool,
		CommitPool:  commitPool,
		RCPool:      rcPool,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
