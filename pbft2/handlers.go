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
	str := fmt.Sprintf("Node[%s] Info:\n", chain_util2.BytesToHex(node.Wallet.publicKey)[:5])
	// validators
	str += "\n[Validators]\n"
	for i, pubKey := range node.Validators.list {
		str += fmt.Sprintf("%d: %s\n", i, chain_util2.BytesToHex(pubKey)[:5])
	}
	// blockchain
	str += "\n[Blockchain]\n"
	for i, block := range node.Blockchain.chain {
		str += fmt.Sprintf("[%s]", chain_util2.BytesToHex(block.Hash)[:5])
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
		str += fmt.Sprintf("%d: %s by %s\n", i, chain_util2.BytesToHex(tx.Hash)[:5], chain_util2.BytesToHex(tx.From)[:5])
	}
	// BlockPool
	str += "\n[BlockPool]\n"
	for i, block := range node.BlockPool.pool {
		str += fmt.Sprintf("%d: %s by %s\n", i, chain_util2.BytesToHex(block.Hash)[:5], chain_util2.BytesToHex(block.Proposer)[:5])
	}
	// PreparePool
	str += "\n[PreparePool]\n"
	for bh, msgs := range node.PreparePool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:5])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util2.BytesToHex(msg.PublicKey)[:5])
		}
	}
	// CommitPool
	str += "\n[CommitPool]\n"
	for bh, msgs := range node.CommitPool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:5])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util2.BytesToHex(msg.PublicKey)[:5])
		}
	}
	// RCPool
	str += "\n[RCPool]\n"
	for bh, msgs := range node.RCPool.mapPool {
		str += fmt.Sprintf("--> %s\n", bh[:5])
		for i, msg := range msgs {
			str += fmt.Sprintf("%d: %s\n", i, chain_util2.BytesToHex(msg.PublicKey)[:5])
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

	// handle any requests by broadcasting it
	// TODO[broadcast storm]: this should result in a broadcast storm, remove it later
	//for {
	//	_, msg, err := conn.ReadMessage()
	//	if err != nil {
	//		log.Printf("Read message failed, [%s]\n", err)
	//		break
	//	}
	//	log.Printf("recv: [%s]\n", msg)
	//
	//	node.broadcast(string(msg))
	//}

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
							if chain_util2.BytesToHex(node.Blockchain.GetProposer()) == chain_util2.BytesToHex(node.Wallet.publicKey) {
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
					} else {
						log.Printf("[MsgTx] NOT PASSING CONDITIONS!, %v, %v, %v\n",
							!node.TxPool.TxExists(tx),
							node.TxPool.VerifyTx(tx),
							node.Validators.ValidatorExists(tx.From))
						fmt.Println(tx.Event)
						fmt.Println(tx.Event)
						//log.Printf("[VerifyTx], %v, %v, %v, %v, %v\n",
						//	tx.MsgType,
						//	MsgTx,
						//	tx.MsgType == MsgTx,
						//	tx.Hash != chain_util2.Hash(tx.Event),
						//	!chain_util2.Verify(chain_util2.Str2Key(tx.From), tx.Hash, tx.Signature))
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
						node.Validators.ValidatorExists(prepareMsg.PublicKey) {
						// add prepareMsg to prepare pool
						node.PreparePool.AddMsg2Pool(prepareMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.PreparePool.mapPool[chain_util2.BytesToHex(prepareMsg.BlockHash)]) >= MIN_APPROVALS {
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
						node.CommitPool.AddMsg2Pool(commitMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.CommitPool.mapPool[chain_util2.BytesToHex(commitMsg.BlockHash)]) >= MIN_APPROVALS {
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
						node.Validators.ValidatorExists(rcMsg.PublicKey) {
						// add rcMsg to rc pool
						node.RCPool.AddMsg2Pool(rcMsg)
						// broadcast
						node.broadcast(string(msg))

						// PBFT MINIMUM VOTING REQUIREMENT
						if len(node.RCPool.mapPool[chain_util2.BytesToHex(rcMsg.BlockHash)]) >= MIN_APPROVALS {
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
