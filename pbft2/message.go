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
 tx      ==(n-to-1)==> block:   	 tx collection, use `Block`
*block   ==(1-to-1)==> prepare: 	 pre-prepare, use `Message` with `type` = "MESSAGE_PRE_PREPARE"
*prepare ==(1-to-1)==> commit:  	 prepare, use `Message` with `type` = "MESSAGE_PREPARE"
*commit  ==(1-to-1)==> round_change: commit, use `Commit` with `type` = "MESSAGE_COMMIT"
*/

// Define Enums
type MsgType int

func (mt MsgType) String() string {
	return MsgName[mt]
}

const (
	MsgPrePrepare = iota
	MsgPrepare
	MsgCommit
)

var MsgName = map[MsgType]string{
	MsgPrePrepare: "PRE-PREPARE",
	MsgPrepare:    "PREPARE",
	MsgCommit:     "COMMIT",
}

/*
*
Message stores passed-in blockHash, publicKey and signature.
1. NewMsg
*/

type Message struct {
	MsgType   MsgType   `json:"msgType"`
	BlockHash string    `json:"blockHash"`
	PublicKey PublicKey `json:"publicKey"`
	Signature string    `json:"signature"`
}

// NewMsg creates a new message that is used for phase transition in PBFT.
func NewMsg(msgType MsgType, blockHash string, publicKey PublicKey, signature string) *Message {
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
	if mp.mapPool[msg.BlockHash] == nil {
		mp.mapPool[msg.BlockHash] = []Message{}
	}
	mp.mapPool[msg.BlockHash] = append(mp.mapPool[msg.BlockHash], msg)
}

// MsgExists checks if a message for a block hash already exists or not
// by comparing its publicKey.
func (mp *MsgPool) MsgExists(msg Message) bool {
	for _, m := range mp.mapPool[msg.BlockHash] {
		if chain_util2.Key2Str(m.PublicKey) == chain_util2.Key2Str(msg.PublicKey) {
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
func (mp *MsgPool) CleanPool(hash string) bool {
	if mp.mapPool[hash] != nil {
		delete(mp.mapPool, hash)
		return true
	} else {
		return false
	}
}
