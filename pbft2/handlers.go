package pbft2

import (
	"consensus-algorithms-with-golang/pbft2/chain_util2"
	"fmt"
	"log"
	"net/http"
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

//1. makeTxHandler
//3. queryBlockPoolHandler
//4. queryPreparePoolHandler
//5. queryCommitPoolHandler
//6. queryRCPoolHandler
//7. queryBlockchainHandler
