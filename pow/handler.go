package pow

import (
	"consensus-algorithms-with-golang/pow/pow_util"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

type MinerInfo struct {
	Hash string `json:"hash"`
}

type BlockInfo struct {
	Hash      string   `json:"hash"`
	Proposer  string   `json:"proposer"`
	Nonce     uint64   `json:"nonce"`
	Timestamp string   `json:"timestamp"`
	Txs       []string `json:"txs"`
}

type Data struct {
	Difficulty  uint64      `json:"difficulty"`
	NodeAddress string      `json:"nodeAddress"`
	Hash        string      `json:"hash"`
	Miners      []MinerInfo `json:"miners"`
	Blockchain  []BlockInfo `json:"blockchain"`
}

func (node *Node) queryNodeInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nodeAddress := fmt.Sprintf("http://%s", pow_util.FormatUrl(node.Host, node.Port))

	nodeHash := pow_util.Byte2Hex(node.Wallet.pubKey)

	mutex.Lock()
	miners := make([]MinerInfo, node.NumOfMiners)
	for i, miner := range node.Miners {
		miners[i] = MinerInfo{miner}
	}
	mutex.Unlock()

	mutex.Lock()
	blockchain := make([]BlockInfo, len(node.Blockchain.chain))
	for i, block := range node.Blockchain.chain {
		txs := make([]string, len(block.Data))
		for j, tx := range block.Data {
			txs[j] = pow_util.Byte2Hex(tx.Hash)
		}

		blockchain[i] = BlockInfo{
			Hash:      pow_util.Byte2Hex(block.Hash),
			Proposer:  pow_util.Byte2Hex(block.Proposer),
			Nonce:     block.Nonce,
			Timestamp: block.Timestamp,
			Txs:       txs,
		}
	}
	mutex.Unlock()

	data := Data{
		Difficulty:  node.Difficulty,
		NodeAddress: nodeAddress,
		Hash:        nodeHash,
		Miners:      miners,
		Blockchain:  blockchain,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

type Request struct {
	Difficulty uint64 `json:"difficulty"`
}

// TODO: add setDifficulty functionality
func (node *Node) setDifficultyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received request body: %s\n", string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request body read successfully!"))
}

func (node *Node) makeTxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	tx := node.Wallet.CreateTx("Tx made at " + time.Now().String())
	msgStr, err := WrapData2MsgStr(*tx)
	if err != nil {
		log.Printf("Marshal tx failed, [%s]\n", err)
		return
	}

	// waiting relay to be available
	for {
		if node.Relay != nil {
			break
		}
	}
	// Write to web page
	w.Write(msgStr)
	// Write to WsClient(Relay)
	mutex.Lock()
	err = node.Relay.WriteMessage(websocket.TextMessage, msgStr)
	mutex.Unlock()
	if err != nil {
		log.Printf("Write message failed, [%s]\n", err)
	}
}

func (node *Node) broadcast(msg string) {
	if _, ok := node.hasBroadcastSet[msg]; !ok {
		node.hasBroadcastSet[msg] = true
		time.Sleep(3 * time.Second)
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
}

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
			log.Printf("Unmarshal message failed, %s, skip this one!\n", err)
			continue
		}

		if msgTypeRaw, ok := data["type"]; ok {
			if msgType, ok2 := msgTypeRaw.(string); ok2 {
				// extract data
				dataRaw, err := json.Marshal(data["data"])
				if err != nil {
					log.Printf("Extract data from msg failed, %s, skip this one!\n", err)
					continue
				}
				switch msgType {
				case MsgTx:
					var tx Transaction
					if err := json.Unmarshal(dataRaw, &tx); err != nil {
						log.Printf("Unmarshal msg->tx failed, %s, skip this one!\n", err)
						continue
					}
					// check if tx is valid
					if !node.TxPool.TxExists(tx) &&
						node.TxPool.VerifyTx(tx) &&
						node.Validators.ValidatorExists(tx.From) {
						log.Printf("Tx [%s] received\n", pow_util.Byte2Hex(tx.Hash)[:6])
						// add tx to txpool
						mutex.Lock()
						node.TxPool.AddTx2Pool(tx) // if not success, the broadcast will try again
						mutex.Unlock()
						// broadcast
						node.broadcast(string(msg))
					}
				case MsgBlock:
					var block Block
					if err := json.Unmarshal(dataRaw, &block); err != nil {
						log.Printf("Unmarshal msg->block failed, %s, skip this one!\n", err)
						continue
					}
					// add block to chain if it passes the verification
					// and its proposer is valid
					if !node.Blockchain.BlockExists(block) &&
						node.Blockchain.VerifyBlock(block) &&
						node.Validators.ValidatorExists(block.Proposer) {
						if node.Blockchain.VerifyBlockWithLastBlockInChain(block) {
							// add block to chain
							mutex.Lock()
							node.Blockchain.AddBlock(block)
							node.TxPool.UpdateCommitted(block, nil)
							log.Printf(
								"Added block [%s] with [%d] txs to blockchain\n"+
									"Block proposed at [%s] by Node-[%s]\n",
								pow_util.Byte2Hex(block.Hash)[:node.Difficulty+4],
								len(block.Data),
								block.Timestamp,
								pow_util.Byte2Hex(block.Proposer)[:4])
							mutex.Unlock()
						} else if node.Blockchain.BlockConflicts(block) &&
							node.Blockchain.Overwritable(block) {
							// overwrite block in chain
							mutex.Lock()
							blocksThatWereOverwritten := node.Blockchain.OverwriteBlock(block)
							node.TxPool.UpdateCommitted(block, blocksThatWereOverwritten)
							log.Printf(
								"Overwrite blocks with block [%s] with [%d] new txs to blockchain\n"+
									"Block proposed at [%s] by Node-[%s]\n",
								pow_util.Byte2Hex(block.Hash)[:node.Difficulty+4],
								len(block.Data),
								block.Timestamp,
								pow_util.Byte2Hex(block.Proposer)[:4])
							mutex.Unlock()
						}
						// broadcast
						node.broadcast(string(msg))
					}
				default:
					log.Printf("recv UNSUPPORTED: %s\n", msg)
				}
			}
		}
	}
}

func (node *Node) resetHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	node.Blockchain.Clear()
	node.TxPool.Clear()
	node.hasBroadcastSet = make(map[string]bool)
	if node.cancel != nil {
		node.cancel()
	}
	mutex.Unlock()
	log.Printf("NODE-[%s] RESET!!!", pow_util.Byte2Hex(node.Wallet.pubKey)[:6])
}
