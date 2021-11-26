# A Documentation
A library to manage Blocks

## Structs

### BlockRecord
```
type BlockRecord struct {
    Hash            cipher.SHA256
    SequenceNo      unint64
    Depth           unint64             // The depth will be the position at which this block is stored in tree.
    Parent          *BlockRecord
    Children        []*BlockRecord
    
    Transactions    []cipher.SHA256              // List of transaction ids
    UxIdSpent       map[cipher.SHA256]*[]btye    // outputs spent / destroyed
    UxIdCreated     map[cipher.SHA256]*[]btye    // ouputs being created
}
```

### BlockRecordTree
In order to use the library we will create a BlockRecordTree instance. Each block tree instance holds a tree of blocks and exposes various methods to interact with them. 
```
type BlockRecordTree struct {
    Root               *BlockRecord                         // Serves as the root of the block tree
    TotalBlocks	       uint64				    	        // Total blocks in the tree
    MaximumDepth       uint64				    	        // Maximum Depth of the tree
    TransactionsMap    map[cipher.SHA256]*Transactions      // Global TxId to *Transactions Map (Transactions struct is defined in [coin/Transactions.go]) 
}
```

### Transactions
A Transaction Record (Copied from src/coin/transactions.go) 
```
type Transaction struct {
	Length    uint32        // length prefix
	Type      uint8         // transaction type
	InnerHash cipher.SHA256 // inner hash SHA256 of In[],Out[]

	Sigs []cipher.Sig        `enc:",maxlen=65535"` // list of signatures, 64+1 bytes each
	In   []cipher.SHA256     `enc:",maxlen=65535"` // ouputs being spent
	Out  []TransactionOutput `enc:",maxlen=65535"` // ouputs being created
}
```


## BlockRecord Routines

### CheckIfUnspentOutputExistsInSpent
- Checks if an unspent output is spent in a block
#### Signature
```
func (b *BlockRecord) CheckIfUnspentOutputExistsInSpent(uxId cipher.SHA256) string {
}
```
#### Output
```
b => {...
    UxIdSpent : {""}
...}
```




### CheckIfUnspentOutputExistsInCreated
- Checks if an unspent output is created in a block
```
func (b *BlockRecord) CheckIfUnspentOutputExistsInCreated(uxId cipher.SHA256) string {
}
```


## BlockRecordTree Routines

### AddBlock
- Adds a block to the block tree
```
func (bt *BlockRecordTree) AddBlock(b *BlockRecord) string {
}
```

### RemoveBlock
- Remove a block from the tree and updates it's parents and children accordingly
```
func (bt *BlockRecordTree) RemoveBlock(b *BlockRecord) string {
}
```

### GetBlockDepth
- Returns the number of blocks between the given block and the root of the block tree
```
func (bt *BlockRecordTree) GetBlockDepth(b *BlockRecord) uint64 {
}
```

### GetAllBlocks
- Returns an list of all blocks from the root to the end of the tree. This function performs a *depth first traversal* of the whole tree returns the list of all blocks 
  in the order they are found. 
```
func (bt *BlockRecordTree) GetAllBlocks() []*BlockRecord {
}
```

### CheckIfUnspentOutputSpendable
- Traverses a tree from root to the given BlockRecord can check if the unspent out was destroyed on it's way from Root to the given block.
- This function can *return* any of the following codes:
  -> "NeverExisted": The unspent output never existed 
  -> "Spent": The unspent output exists but was spent 
  -> "Available": The unspent output is available to spend on the given block
  
```
func (bt *BlockRecordTree) CheckIfUnspentOutputSpendable(UnspentOutput cipher.SHA256, targetBlock *BlockRecord) string {
}
```

### CheckIfMultipleUnspentOutputSpendable
- Traverses a tree from root to the given BlockRecord can check if all the outputs are spendable at the given block.
- This function returns 
  -> true: If all the outputs are spendable at the given block 
  -> false: If anyone of the output is not spendable at the given block 
```
func (bt *BlockRecordTree) CheckIfMultipleUnspentOutputSpendable(UnspentOutputs []cipher.SHA256, targetBlock *BlockRecord) bool {
}
```
