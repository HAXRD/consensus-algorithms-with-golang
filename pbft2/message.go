package pbft2

import "consensus-algorithms-with-golang/pbft2/chain_util2"

/**
PBFT uses 3 phases to ensure consensus, pre-prepare, prepare, and commit.
- Pre-prepare (solve synchronization issue): ensuring all the
  nodes know the proposal;
- Prepare (fault tolerance): making sure sufficient number of
  nodes agree upon the proposal;
- Commit (Finality): finalize the commit, making sure the proposal
  is final and will not be altered.

Workflow: '*' tagged are included in this file
 tx      ==(n-to-1)==> block:   	 "TX"
*block   ==(1-to-1)==> prepare: 	 "PRE-PREPARE"
*prepare ==(1-to-1)==> commit:  	 "PREPARE"
*commit  ==(1-to-1)==> round_change: "COMMIT"
 round_change =======> new_round:    "RC"
*/

// Define MsgTypes
const (
	MsgTx         = "Tx"
	MsgPrePrepare = "PRE-PREPARE"
	MsgPrepare    = "PREPARE"
	MsgCommit     = "COMMIT"
	MsgRC         = "RC"
)

/*
*
Message stores passed-in blockHash, publicKey and signature.
1. NewMsg
*/

type Message struct {
	MsgType   string    `json:"msgType"`
	BlockHash []byte    `json:"blockHash"`
	PublicKey PublicKey `json:"publicKey"`
	Signature []byte    `json:"signature"`
}

// NewMsg creates a new message that is used for phase transition in PBFT.
func NewMsg(msgType string, blockHash []byte, publicKey PublicKey, signature []byte) *Message {
	return &Message{
		MsgType:   msgType,
		BlockHash: blockHash,
		PublicKey: publicKey,
		Signature: signature,
	}
}

/**
MessagePool stores a pool of messages with a specified message type.
With the same block hash as map key, each element in the pool
represents a message sent from a different node.
It features the following methods:
1. NewMsgPool
2. AddMsg2Pool: pushes a message for a block hash into the map list
3. MsgExists: check if a given message for a block hash already exists
4. VerifyMsg: check if the message is valid or not
5. CleanPool: remove the list with the specified block hash in the map pool
*/

type MsgPool struct {
	mapPool map[string][]Message
}

// NewMsgPool creates a message pool. It uses a map to store a list of
// messages, whose key is the passed-in block hash.
func NewMsgPool() *MsgPool {
	return &MsgPool{
		mapPool: make(map[string][]Message),
	}
}

// AddMsg2Pool adds a message to the pool. It first checks if
// the map pool already has the specified message list whose
// key is message's block hash, init a list if not. Then it
// add the message to the map pool.
func (mp *MsgPool) AddMsg2Pool(msg Message) {
	hashHex := chain_util2.BytesToHex(msg.BlockHash)
	if mp.mapPool[hashHex] == nil {
		mp.mapPool[hashHex] = []Message{}
	}
	mp.mapPool[hashHex] = append(mp.mapPool[hashHex], msg)
}

// MsgExists checks if a message for a block hash already exists or not
// by comparing its publicKey.
func (mp *MsgPool) MsgExists(msg Message) bool {
	hashHex := chain_util2.BytesToHex(msg.BlockHash)
	for _, m := range mp.mapPool[hashHex] {
		if chain_util2.BytesToHex(m.PublicKey) == chain_util2.BytesToHex(msg.PublicKey) {
			return true
		}
	}
	return false
}

// VerifyMsg verifies the passed-in message
func (mp *MsgPool) VerifyMsg(msg Message) bool {
	return chain_util2.Verify(msg.PublicKey, msg.BlockHash, msg.Signature)
}

// CleanPool remove the list with specified block hash in map pool
func (mp *MsgPool) CleanPool(hash []byte) bool {
	hashHex := chain_util2.BytesToHex(hash)
	if mp.mapPool[hashHex] != nil {
		delete(mp.mapPool, hashHex)
		return true
	} else {
		return false
	}
}
