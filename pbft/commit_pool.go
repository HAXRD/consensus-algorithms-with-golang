package pbft

import "consensus-algorithms-with-golang/pbft/chainutil"

/**
Commit is a type of message that stores passed-in blockHash, publicKey and signature.
It is used to guarantee the FINALITY
1. NewCommit: create a commit message with given blockHash, publicKey and signature.
*/

type Commit struct {
	blockHash string `json:"blockHash"`
	publicKey string `json:"publicKey"`
	signature string `json:"signature"`
}

func NewCommit(wallet Wallet, prepare Prepare) *Commit {
	return &Commit{blockHash: prepare.blockHash, publicKey: chainutil.Key2Str(wallet.publicKey), signature: wallet.Sign(prepare.blockHash)}
}

/**
CommitPool stores a list of commit message for "each" prepare (actually block info) message

1. NewCommitPool: create a new commit pool
2. InitCommitForGivenPrepare: init an empty list for a given blockHash
3. AddCommit2Pool: pushed a commit message for a blockHash into the list
4. CommitExists: check if a given commit message already exists
5. IsCommitValid: check if the commit message is valid or not
*/

type CommitPool struct {
	mapOfList map[string][]Commit
}

func NewCommitPool() *CommitPool {
	return &CommitPool{mapOfList: make(map[string][]Commit)}
}

func (cp *CommitPool) InitCommitForGivenPrepare(blockHash string) {
	cp.mapOfList[blockHash] = []Commit{}
}

func (cp *CommitPool) AddCommit2Pool(commit Commit) {
	cp.mapOfList[commit.blockHash] = append(cp.mapOfList[commit.blockHash], commit)
}

func (cp *CommitPool) CommitExists(commit Commit) bool {
	for _, c := range cp.mapOfList[commit.blockHash] {
		if c.blockHash == commit.blockHash {
			return true
		}
	}
	return false
}

func (cp *CommitPool) IsCommitValid(commit Commit) bool {
	return chainutil.Verify(commit.publicKey, commit.blockHash, commit.signature)
}
