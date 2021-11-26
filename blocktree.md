# A Documentation
A library to manage Block Records

## Structs

### BlockRecord
```
type BlockRecord struct {
    Hash            cipher.SHA256
    SequenceNo      unint64
    Depth           unint64             // The depth will be the position at which this blockRecord is stored in tree.
    Parent          *BlockRecord
    Children        []*BlockRecord
    
    Transactions    []cipher.SHA256              // List of transaction ids
    UxIdSpent       map[cipher.SHA256]*[]btye    // outputs spent / destroyed
    UxIdCreated     map[cipher.SHA256]*[]btye    // ouputs being created
}
```

### BlockRecordTree
In order to use the library we will create a BlockRecordTree instance. Each blockRecord tree instance holds a tree of blocks and exposes various methods to interact with them. 
```
type BlockRecordTree struct {
    Root               *BlockRecord                         // Serves as the root of the blockRecord tree
    TotalBlocks	       uint64				    // Total blocks in the tree
    MaximumDepth       uint64				    // Maximum Depth of the tree
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
- Checks if an unspent output is spent in a blockRecord
#### Signature
```
func (b *BlockRecord) CheckIfUnspentOutputExistsInSpent(uxId cipher.SHA256) string {
}
```
#### Sample Run / Tests
```
let 

b => {...
    UxIdSpent : {
        "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4": [0xb1, 0xb2, 0xb3]}
...}

then

b.CheckIfUnspentOutputExistsInSpent("03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4")  => true
b.CheckIfUnspentOutputExistsInSpent("xxxxxxxx")  => false

```

### CheckIfUnspentOutputExistsInCreated
- Checks if an unspent output is created in a blockRecord
#### Signature
```
func (b *BlockRecord) CheckIfUnspentOutputExistsInCreated(uxId cipher.SHA256) string {
}
```
#### Sample Run / Tests
```
let 

b => {...
    UxIdCreated : {
        "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4": [0xb1, 0xb2, 0xb3]}
...}

then

b.CheckIfUnspentOutputExistsInCreated("03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4")  => true
b.CheckIfUnspentOutputExistsInCreated("xxxxxxxx")  => false

```


## BlockRecordTree Routines

### AddBlock
- Adds a blockRecord to the blockRecord tree
#### Signature
```
func (bt *BlockRecordTree) AddBlock(b *BlockRecord) string {
}
```
#### Sample Run / Tests
```
let 

bt => {...
    Root : nil
...}

b1 => { ... Hash : "<Sha256#1>" ...}
b2 => { ... Hash : "<Sha256#2>" ...}

then

bt.AddBlock(b1)
bt.AddBlock(b2)

bt => {...
    Root : { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#2>", Depth: 2 ...},
    TotalBlocks: 2,
    MaximumDepth: 2
...}

```

### RemoveBlock
- Remove a blockRecord from the tree and updates it's parents and children accordingly
#### Signature
```
func (bt *BlockRecordTree) RemoveBlock(b *BlockRecord) string {
}
```
#### Sample Run / Tests
```
let 

bt => {...
    Root : { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#2>", Depth: 2 ...} -> { ... Hash : "<Sha256#3>", Depth: 3 ...},
    TotalBlocks: 3,
    MaximumDepth: 3
...}

b2 => { ... Hash : "<Sha256#2>", Depth: 2 ...}

then

bt.RemoveBlock(b2)

bt => {...
    Root : { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#3>", Depth: 2 ...}.
    TotalBlocks: 2,
    MaximumDepth: 2
...}

```

### GetBlockDepth
- Returns the number of blocks between the given blockRecord and the root of the blockRecord tree
#### Signature
```
func (bt *BlockRecordTree) GetBlockDepth(b *BlockRecord) uint64 {
}
```
#### Sample Run / Tests
```
let 

bt => {...
    Root : { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#2>", Depth: 2 ...} -> { ... Hash : "<Sha256#3>", Depth: 3 ...},
    TotalBlocks: 3,
    MaximumDepth: 3
...}

b1 => { ... Hash : "<Sha256#1>", Depth: 1 ...}
b2 => { ... Hash : "<Sha256#2>", Depth: 2 ...}
b3 => { ... Hash : "<Sha256#3>", Depth: 3 ...}

then

bt.GetBlockDepth(b1) => 1
bt.GetBlockDepth(b2) => 2
bt.GetBlockDepth(b3) => 3

```


### GetAllBlocks
- Returns an list of all blocks from the root to the end of the tree. This function performs a *depth first traversal* of the whole tree returns the list of all blocks in the order they are found. 
#### Signature
```
func (bt *BlockRecordTree) GetAllBlocks() []*BlockRecord {
}
```
#### Sample Run / Tests
```
let 

bt => {...
    Root : { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#2>", Depth: 2 ...} -> { ... Hash : "<Sha256#3>", Depth: 3 ...},
    TotalBlocks: 3,
    MaximumDepth: 3
...}

then

bt.GetAllBlocks() => { ... Hash : "<Sha256#1>", Depth: 1 ...} -> { ... Hash : "<Sha256#2>", Depth: 2 ...} -> { ... Hash : "<Sha256#3>", Depth: 3 ...}
```

### CheckIfUnspentOutputSpendable
- Traverses a tree from root to the given BlockRecord can check if the unspent out was destroyed on it's way from Root to the given blockRecord.
- This function can *return* any of the following codes:
  -> "NeverExisted": The unspent output never existed 
  -> "Spent": The unspent output exists but was spent 
  -> "Available": The unspent output is available to spend on the given blockRecord
  
#### Signature
```
func (bt *BlockRecordTree) CheckIfUnspentOutputSpendable(UnspentOutput cipher.SHA256, targetBlock *BlockRecord) string {
}
```
#### Sample Run / Tests
```
let 

b1 => { ... UxIdCreated : {"<Sha256#1>": [], "<Sha256#2>": [] "<Sha256#3>": []} ...} 
b2 => { ... UxIdSpent : {"<Sha256#3>": []},  ...} 
b3 => { ... UxIdCreated : {"<Sha256#4>": []} ...}

bt => {...
    Root : b1 -> b2 -> b3,
    TotalBlocks: 3,
    MaximumDepth: 3
...}

then

bt.CheckIfUnspentOutputSpendable("<Sha256#1>", b1) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#2>", b1) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#3>", b1) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#4>", b1) => "NeverExisted" 

bt.CheckIfUnspentOutputSpendable("<Sha256#1>", b2) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#2>", b2) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#3>", b2) => "Spent" 
bt.CheckIfUnspentOutputSpendable("<Sha256#4>", b2) => "NeverExisted" 

bt.CheckIfUnspentOutputSpendable("<Sha256#1>", b3) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#2>", b3) => "Available" 
bt.CheckIfUnspentOutputSpendable("<Sha256#3>", b3) => "Spent" 
bt.CheckIfUnspentOutputSpendable("<Sha256#4>", b3) => "Available" 
```

### CheckIfMultipleUnspentOutputSpendable
- Traverses a tree from root to the given BlockRecord can check if all the outputs are spendable at the given blockRecord.
- This function returns 
  -> true: If all the outputs are spendable at the given blockRecord 
  -> false: If anyone of the output is not spendable at the given blockRecord 
#### Signature
```
func (bt *BlockRecordTree) CheckIfMultipleUnspentOutputSpendable(UnspentOutputs []cipher.SHA256, targetBlock *BlockRecord) bool {
}
```
#### Sample Run / Tests
```
let 

b1 => { ... UxIdCreated : {"<Sha256#1>": [], "<Sha256#2>": [] "<Sha256#3>": []} ...} 
b2 => { ... UxIdSpent : {"<Sha256#3>": []},  ...} 
b3 => { ... UxIdCreated : {"<Sha256#4>": []} ...}

bt => {...
    Root : b1 -> b2 -> b3,
    TotalBlocks: 3,
    MaximumDepth: 3
...}

then

bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>", "<Sha256#3>"], b1) => true 
bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>", "<Sha256#3>", "<Sha256#4>"], b1) => false 

bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>"], b2) => true 
bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>", "<Sha256#3>"], b2) => false 

bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>", "<Sha256#4>"], b3) => true 
bt.CheckIfMultipleUnspentOutputSpendable(["<Sha256#1>", "<Sha256#2>", "<Sha256#3>", "<Sha256#4>"], b3) => false 
```
