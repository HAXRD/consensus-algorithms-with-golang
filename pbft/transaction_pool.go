package pbft

/**
TransactionPool is where all the transactions from nodes are stored.
It has the following functions:
1. add transaction to the pool
2. check if a transaction is valid
3. check if a transaction exists or not
4. empties the pool
*/

type TransactionPool struct {
	transactions []Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{transactions: make([]Transaction, 0, TX_THRESHOLD)}
}

// TransactionExists check if a given tx is in the transaction pool
func (tp TransactionPool) TransactionExists(tx Transaction) bool {
	for _, _tx := range tp.transactions {
		if _tx.id == tx.id {
			return true
		}
	}
	return false
}

// AddTransaction adds a new tx to the pool
// returns true if reaches pool threshold
// return false otherwise
func (tp *TransactionPool) AddTransaction(tx Transaction) bool {
	tp.transactions = append(tp.transactions, tx)
	if len(tp.transactions) >= TX_THRESHOLD {
		return true
	}
	return false
}

// VerifyTransaction check if the transaction is valid or not
func (tp *TransactionPool) VerifyTransaction(tx Transaction) bool {
	return VerifyTransaction(tx)
}
