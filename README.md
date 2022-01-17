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
- [Struct CommunicationsDelayMatrix](#struct-communicationsdelaymatrix)
  * [Data](#data-2)
  * [Methods](#methods-2)
    + [InitializeCommunicationsDelayMatrix](#initializecommunicationsdelaymatrix)
      - [Signature](#signature-11)
- [Struct NodeGridMap](#struct-nodegridmap)
  * [Data](#data-3)
  * [Methods](#methods-3)
    + [InitializeNodeGridMap](#initializenodegridmap)
      - [Signature](#signature-12)
- [Struct NodeGrid](#struct-nodegrid)
  * [Data](#data-4)
  * [Methods](#methods-4)
    + [InitializeNodeGrid](#initializenodegrid)
      - [Signature](#signature-13)
- [Struct Node](#struct-node)
  * [Data](#data-5)
  * [Methods](#methods-5)
    + [InitializeNode](#initializenode)
      - [Signature](#signature-14)
    + [ValidateNodeState](#validatenodestate)
      - [Signature](#signature-15)
    + [UpdateNodeState](#updatenodestate)
      - [Signature](#signature-16)
- [Struct NodeBlockMeta](#struct-nodeblockmeta)
  * [Data](#data-6)
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

## Struct CommunicationsDelayMatrix
The CommunicationsDelayMatrix struct tracks the tick delay when sending messages between two nodes
### Data
```
type CommunicationsDelayMatrix struct {
	matrix map[cipher.PubKey]map[cipher.PubKey]int
}
```
### Methods
Following are the methods supported by the CommunicationsDelayMatrix Struct required for this simulation
#### InitializeCommunicationsDelayMatrix
Initializes the CommunicationsDelayMatrix with random values (Right now all of them are set to 1s)
##### Signature
```
func (cdm *CommunicationsDelayMatrix) InitializeCommunicationsDelayMatrix(nodes []*Node) {}
```

## Struct NodeGridMap
The NodeGridMap struct Keeps a mapping for Each node's PubKey to it's NodeGrid
### Data
```
type NodeGridMap struct {
	gridMap map[cipher.PubKey]*NodeGrid
}
```
### Methods
Following are the methods supported by the NodeGridMap Struct required for this simulation
#### InitializeNodeGridMap
Initializes the NodeGrid for each node by randomly placing other nodes on the grid
##### Signature
```
func (ngm *NodeGridMap) InitializeNodeGridMap(nodes []*Node) {}
```

## Struct NodeGrid
The NodeGrid struct represents a grid of nodes showing their distances to a particular node
### Data
```
type NodeGrid struct {
	grid [][]*Node
}
```
### Methods
Following are the methods supported by the NodeGrid Struct required for this simulation
#### InitializeNodeGrid
Initializes the NodeGrid for a given node by randomly placing other nodes on the grid
##### Signature
```
func (ng *NodeGrid) InitializeNodeGrid(initNode *Node, nodes []*Node) {}
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
- Adds specified number of subscribers to each node based on the subscription assignment policy. Following two policies are supported:
  - NODE_INIT_RANDOM_SUBSCRIPTION_POLICY (by randomly choosing subscribers from the overall nodes available)
  - NODE_INIT_RANDOM_ANNUAL_RING_SUBSCRIPTION_POLICY *[Currently Active]* (via Annual Ring Algorithm using Node Grids)
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
<dir-Path>/obelisk$ ./simulation -block-record-count 3 -children-per-block 2 -nodes 3 -subscribers 2 -iterations 1000
```
### Sample Output
```console

$ ./simulation -block-record-count 3 -children-per-block 2 -nodes 3 -subscribers 2 -iterations 1000

#begin Simulation Initial State:

Node (id=1 seqNo=0) Details:
PubKey:8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025
Subscriptions:[2 2]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.49
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.51

Node (id=2 seqNo=0) Details:
PubKey:be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a
Subscriptions:[1 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.51
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.49

Node (id=3 seqNo=0) Details:
PubKey:469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9
Subscriptions:[1 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 0 | 0 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.49
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0 | 0 | 0.51

#end Simulation Initial State



#begin Simulation Final State:

Node (id=1 seqNo=6) Details:
PubKey:8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025
Subscriptions:[2 2]
Messages Received [Format: from(pubKey) | to(pubKey) | sentTick | arrivedTick | message]:
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 1 | 2 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 1 | 2 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 3 | 4 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 3 | 4 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 4 | 5 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 4 | 5 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 5 | 6 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 5 | 6 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 6 | 7 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 6 | 7 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 8 | 9 | Hello from Node 3 to Node 1
469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9 | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 8 | 9 | Hello from Node 3 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 9 | 10 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 9 | 10 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 11 | 12 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 11 | 12 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 12 | 13 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 12 | 13 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 13 | 14 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 13 | 14 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 16 | 17 | Hello from Node 2 to Node 1
be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | 16 | 17 | Hello from Node 2 to Node 1
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 5 | 17 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 5 | 17 | 1.00
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 5 | 17 | 0.00

Node (id=2 seqNo=5) Details:
PubKey:be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a
Subscriptions:[1 1]
Messages Received [Format: from(pubKey) | to(pubKey) | sentTick | arrivedTick | message]:
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 2 | 3 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 2 | 3 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 7 | 8 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 7 | 8 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 10 | 11 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 10 | 11 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 14 | 15 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 14 | 15 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 15 | 16 | Hello from Node 1 to Node 2
8da86849fb1e18e43cb472164538058f74280ac93645731d1f9f3b9989575ba025 | be587209261ee28c951fa22aca128ccc1c7fa68a3f29a3f2b0f67959033383ff7a | 15 | 16 | Hello from Node 1 to Node 2
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 5 | 16 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 5 | 16 | 1.00
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 5 | 16 | 0.00

Node (id=3 seqNo=7) Details:
PubKey:469a59a6ba9dbd86be3cb76ddd7070905650b846e54fdeb4b870c287ed7fe3c4c9
Subscriptions:[1 1]
State [Format: blockHash | parentHash | seqNo | ticks | weight]:
8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 0000000000000000000000000000000000000000000000000000000000000000 | 6 | 18 | 1.00
bef619d411d63d1323bfc74c2740d30e4a32e76b2d54e762504c076dfe875a9f | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 6 | 18 | 1.00
6b58c0b40df27ce0c14f6d9f94b7a78b25b95f5b832e8c7aa312c4993f52acd8 | 8da81b6972bff1a7c4a85f17c2b3cb6e86239c7c41ab3b6dae15491e2139e200 | 6 | 18 | 0.00

#end Simulation Final State

Iteration 17: Convergence Achieved!!!

```
