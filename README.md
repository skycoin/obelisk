# Obelisk Sim Readme
This Script performs a simulation of the obelisk consensus algorithm.

- [Global Public Methods](#global-public-methods)
    + [GetRandomPubKey](#getrandompubkey)
      - [Signature](#signature)
    + [GetRandomSHA256](#getrandomsha256)
      - [Signature](#signature-1)
    + [GetSimulation](#getsimulation)
      - [Signature](#signature-2)
    + [GenerateRandomBlockTree](#generaterandomblocktree)
      - [Signature](#signature-3)
    + [NewRandomBlockRecord](#newrandomblockrecord)
      - [Signature](#signature-4)
    + [NewRandomNode](#newrandomnode)
      - [Signature](#signature-5)
    + [NewNodeStateBlockMeta](#newnodestateblockmeta)
      - [Signature](#signature-6)
- [Struct Simulation](#struct-simulation)
  * [Data](#data)
  * [Methods](#methods)
    + [AdvanceTicks](#advanceticks)
      - [Signature](#signature-7)
    + [InitSimulation](#initsimulation)
      - [Signature](#signature-8)
    + [PrintAllNodes](#printallnodes)
      - [Signature](#signature-9)
- [Struct BlockRecordTree](#struct-blockrecordtree)
  * [Data](#data-1)
  * [Methods](#methods-1)
    + [GetAllBlockRecords](#getallblockrecords)
      - [Signature](#signature-10)
- [Struct BlockRecord](#struct-blockrecord)
  * [Data](#data-2)
- [Struct Node](#struct-node)
  * [Data](#data-3)
  * [Methods](#methods-2)
    + [InitializeNode](#initializenode)
      - [Signature](#signature-11)
    + [ValidateNodeState](#validatenodestate)
      - [Signature](#signature-12)
    + [UpdateNodeState](#updatenodestate)
      - [Signature](#signature-13)
- [Struct NodeBlockMeta](#struct-nodeblockmeta)
  * [Data](#data-4)
- [Overall Flow](#overall-flow)
- [Dry Run](#dry-run) 
- [How to build / run?](#how-to-build---run-)
  * [Sample Run](#sample-run)
  * [Sample Output](#sample-output)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

## Global Public Methods
Following are the public methods that will be used by the cmd script
#### GetRandomPubKey
Generates a random cipher.PubKey.
##### Signature
```
func GetRandomPubKey() cipher.PubKey
```
#### GetRandomSHA256
Generates a random sha256 hash. It Basically generates a random number and then perform sha256 hash on it. 
##### Signature
```
func GetRandomSHA256() cipher.SHA256
```
#### GetSimulation
This method returns the globally active singleton simulation object
##### Signature
```
func GetSimulation() *Simulation
```
#### GenerateRandomBlockTree
This method recursively generates a random block tree for a given number of nodes and a given number of children of each node. 
- It will traverse the tree in the breadth first search manner and will keep adding blocks until totalBlocks is reached. 
- To Add a new block. It create a BlockRecord struct, use InitializeRandomBlock to generate BlockRecord with pre-initialized Hash and Parent Node set
##### Signature
```
func NewRandomBlockRecordTree(totalBlocks int, childrenPerNode int) (*BlockRecordTree, error) {}
```

#### NewRandomBlockRecord
Creates a block record object based on random Hash Value. It basically
1- Generates a random hash (cipher.SHA256) using GetRandomBlockHash and assigned it to b.Hash
2- and sets: b.seqNo = 0 | b.Parent = parent (parameter) | b.Children = []*BlockRecord{}  
##### Signature
```
func NewRandomBlockRecord() *BlockRecord {}
```
#### NewRandomNode
Creates a random node object with given id and initializes it with a random public key
##### Signature
```
func NewRandomNode(id int) *Node {}
```
#### NewNodeStateBlockMeta
Creates a random node block meta with a given block record
##### Signature
```
func NewNodeStateBlockMeta(blockRecord *BlockRecord) *NodeStateBlockMeta {}
```

## Struct Simulation
The Simulation struct will hold all the data required for the running instance of the obelisk simulation. This struct will be maintained as singleton
### Data
```
type Simulation struct {
	VerboseMode   bool                    // Simulation Running in Verbose Mode
	Iterations    int                     // Number of Iterations to run the simulation
	Ticks         int                     // Number of Simulation Ticks
	Nodes         []*Node                 // Nodes in Simulation
	RootBlockTree *BlockRecordTree        // Root Block Tree for Simulation
}
```
### Methods
Following are the methods supported by the Simulation Struct
#### AdvanceTicks
This method simply increments ticks held by the running simulation application. It tracks the total number of updates done on any of the nodes.
##### Signature
```
func (sim *Simulation) AdvanceTicks() {}
```

#### InitSimulation
This method initializes the global simulation object based on command line arguments
##### Signature
```
func (sim *Simulation) InitSimulation(totalRootBlockTreeNodes int, totalRootBlockTreeChildrenPerNode int, numberOfNodes int, numberOfSubscribers int, iterations int, verboseMode bool) {}
```

#### PrintAllNodes
Print all nodes along with their states as csv
##### Signature
```
func (sim *Simulation) PrintAllNodes() {}
```
## Struct BlockRecordTree
The BlockRecordTree struct will hold data for the root block tree
### Data
```
type BlockRecordTree struct {
	Root *BlockRecord           // Root node of the Block Record Tree
}
```
### Methods
Following are the methods supported by the BlockRecordTree Struct required for this simulation
#### GetAllBlockRecords
Returns all blocks of the root block tree as a list
##### Signature
```
func (brt *BlockRecordTree) GetAllBlockRecords() []*BlockRecord {}
```

## Struct BlockRecord
The BlockRecord struct will hold data to simulate a Block Record
### Data
```
type BlockRecord struct {
	hash     cipher.SHA256  // Hash of the Block
	seqNo    int            // SeqNo of the block
	parent   *BlockRecord   // Pointer to the parent of the block record
	children []*BlockRecord // List of children of the block record
}
```
## Struct Node
The Node struct holds the Node information for the running simulation
### Data
```
type Node struct {
	id            int                                   // id of the node
	pubKey        cipher.PubKey                         // Node's public key
	seqNo         int                                   // Node's sequence number tracking the number of updates done on the node
	subscriptions []*Node                               // List of Nodes subscribed by the current Node
	state         map[cipher.SHA256]*NodeStateBlockMeta // A mapping from BlockRecord Hash to current Node's separate copy of NodeStateBlockMeta
}
```
### Methods
Following are the methods supported by the Node Struct required for this simulation
#### InitializeNode
Initializes the current node's state:
- Iterate through the global block record tree held by Simulation struct 
- Foreach of the block record adds it to the state and then initialize the weight = (weight of parent) / (number of children of parent)
- Adds number of subscribers to the node in a random fashion driven by seed
##### Signature
```
func (n *Node) InitializeNode(brt *BlockRecordTree, nodes []*Node, numberOfSubscribers int) {}
```
#### ValidateNodeState
For NodeBlockMeta entry in current node's state, verifies that it's weight is equal to that sum of the weights of it's children
##### Signature
```
func (n *Node) ValidateNodeState() {
}
```

#### UpdateNodeState
Updates the current node's state:
1- Increment the current Node's seqNo By 1
2- Get the state of each of the subscribed nodes
- Foreach NodeBlockMeta in current node's state:
    - copy the highest seqNo from the corresponding NodeBlockMeta(s) of the subscribed nodes' (correspondences can be done by hash of the block record).
    - get avg of the weights of the corresponding NodeBlockMeta(s) of the subscribed nodes' (correspondences can be done by hash of the block record).
    - assign the above calculated seqNO and avgWeight to the current Node's NodeBlockMeta.
##### Signature
```
func (n *Node) UpdateNodeState() {
}
```
## Struct NodeBlockMeta
The NodeBlockMeta struct holds each node's individual copy of block record details.
### Data
```
struct NodeBlockMeta {
    blockRecord   *NodeBlockRecord    // The corresponding BlockRecord from the tree
    seqNo         uint64              // Here we maintain the highest seqNo among the nodes considered while syncing the states  
    ticks         uint64              // capture the ticks from the global simulation at the time when the NodeBlockMeta was synced
    weight        float              // weight of NodeBlockMeta. This weight will be sum of the weight of the children of the blockRecord  
}
```
## Overall Flow
- Print each node id along with public key
- Print each node along with it's subscribers 
- For the number of iterations provided:
  - Based on the given seed, generate a random number in range (1, number of nodes)
  - Get the corresponding node based on the index of the node at above generated random number and call node.UpdateState.

## Dry Run
1- Print each node id along with public key
```console
N1 => 1, <PK-Node1>
N2 => 2, <PK-Node2>
N3 => 3, <PK-Node3>
```
2- Print each node along with it's subscribers
```console
N1 => 1, [2,3]
N2 => 2, [1,3]
N3 => 3, [1,2]
```

3- Generate a Random Block Tree
```
b1 = { HASH-B1 }
b2 = { HASH-B2 }
b3 = { HASH-B3 }

simulation.Root = b1 => [b2, b3]
```

4- Initialize Each Node's state as per the block tree. NodeBlockMeta is represented as {BlockRecord, seqNo, ticks, weight}
```
N1 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
N2 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

5- Here we perform three iterations for simulation
##### Iteration 1 (Update State Called on N1)
*Before*
```
simulation.ticks = 0
(seq=0) N1 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
(seq=0) N2 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
(seq=0) N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

*After*
```
simulation.ticks = 1
(seq=1) N1 => {b1, 0, 1, (1+1)/2=1.00} => [ {b2, 0, 1, (0.5+0.5)/2=0.50} {b3, 0, 1, (0.5+0.5)/2=0.50} ]
(seq=0) N2 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
(seq=0) N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

##### Iteration 2 (Update State Called on N2)
*Before*
```
simulation.ticks = 1
(seq=1) N1 => {b1, 0, 1, 1.00} => [ {b2, 0, 1, 0.50} {b3, 0, 1, 0.50} ]
(seq=0) N2 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
(seq=0) N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

*After*
```
simulation.ticks = 2
(seq=1) N1 => {b1, 0, 1, 1.00} => [ {b2, 0, 1, 0.50} {b3, 0, 1, 0.50} ]
(seq=1) N2 => {b1, 1, 2, (1+1)/2=1.00} => [ {b2, 1, 2, (0.5+0.5)/2=0.50} {b3, 1, 2, (0.5+0.5)/2=0.50} ]
(seq=0) N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

##### Iteration 3 (Update State Called on N3)
*Before*
```
simulation.ticks = 2
(seq=1) N1 => {b1, 0, 1, 1.00} => [ {b2, 0, 1, 0.50} {b3, 0, 1, 0.50} ]
(seq=1) N2 => {b1, 1, 2, 1.00} => [ {b2, 1, 2, 0.50} {b3, 1, 2, 0.50} ]
(seq=0) N3 => {b1, 0, 0, 1.00} => [ {b2, 0, 0, 0.50} {b3, 0, 0, 0.50} ]
```

*After*
```
simulation.ticks = 3
(seq=1) N1 => {b1, 0, 1, 1.00} => [ {b2, 0, 1, 0.50} {b3, 0, 1, 0.50} ]
(seq=1) N2 => {b1, 1, 2, 1.00} => [ {b2, 1, 2, 0.50} {b3, 1, 2, 0.50} ]
(seq=1) N3 => {b1, 1, 3, (1+1)/2=1.00} => [ {b2, 1, 3, (0.5+0.5)/2=0.50} {b3, 1, 3, (0.5+0.5)/2=0.50} ]
```

*Final State*
```
simulation.ticks = 3
(seq=1) N1 => {b1, 0, 1, 1.00} => [ {b2, 0, 1, 0.50} {b3, 0, 1, 0.50} ]
(seq=1) N2 => {b1, 1, 2, 1.00} => [ {b2, 1, 2, 0.50} {b3, 1, 2, 0.50} ]
(seq=1) N3 => {b1, 1, 3, 1.00} => [ {b2, 1, 3, 0.50} {b3, 1, 3, 0.50} ]
```



## How to build / run?  
The simulation will be run as a command line script
### Sample Run
```console
<dir-Path>/obelisk$ go build ./src/simulation
<dir-Path>/obelisk$ ./simulation -block-record-count 3 -children-per-block 2 -nodes 3 -subcribers 2 -iterations 1000
```
### Sample Output
```console

#begin Simulation Initial State:

Node (id=1 seqNo=0) Details:
PubKey:[180 133 196 192 101 210 240 238 186 238 131 109 78 169 31 210 92 6 244 61 173 119 238 38 66 189 195 117 53 37 50 228 59]
Subscriptions:[3 2]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.51
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.49

Node (id=2 seqNo=0) Details:
PubKey:[146 164 121 54 187 205 148 113 120 248 228 192 216 122 59 231 3 196 155 225 216 131 160 124 72 135 219 148 78 63 15 241 105]
Subscriptions:[3 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.51
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.49

Node (id=3 seqNo=0) Details:
PubKey:[28 248 23 47 180 196 85 145 40 35 217 154 84 193 190 168 128 243 25 245 109 162 150 19 153 132 16 19 155 68 106 227 47]
Subscriptions:[2 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.51
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0 | 0 | 0.49

#end Simulation Initial State



#begin Simulation Final State:

Node (id=1 seqNo=8) Details:
PubKey:[180 133 196 192 101 210 240 238 186 238 131 109 78 169 31 210 92 6 244 61 173 119 238 38 66 189 195 117 53 37 50 228 59]
Subscriptions:[3 2]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 9 | 23 | 1.00
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 9 | 23 | 1.00
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 9 | 23 | 0.00

Node (id=2 seqNo=7) Details:
PubKey:[146 164 121 54 187 205 148 113 120 248 228 192 216 122 59 231 3 196 155 225 216 131 160 124 72 135 219 148 78 63 15 241 105]
Subscriptions:[3 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 9 | 24 | 1.00
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 9 | 24 | 1.00
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 9 | 24 | 0.00

Node (id=3 seqNo=9) Details:
PubKey:[28 248 23 47 180 196 85 145 40 35 217 154 84 193 190 168 128 243 25 245 109 162 150 19 153 132 16 19 155 68 106 227 47]
Subscriptions:[2 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 0000000000000000000000000000000000000000000000000000000000000000 | 7 | 22 | 1.00
19722d7643e4fa884d197907538c91595753a8bced44dd569db2814b1cfb0e0d | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 7 | 22 | 1.00
3bf0835bfca39af7786e55e3aefd6a9f2abc95a5d9e2ecef54a3cafa50781e67 | fae99589eb0dc2c68982e8c1bc7b5f99ffa704806bc75d76bfbb6302d08c60a6 | 7 | 22 | 0.00

#end Simulation Final State

Iteration 23: Convergence Achieved!!!

```
