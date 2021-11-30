# Obelisk Sim Readme
This Script performs a simulation of the obelisk consensus algorithm

## Global Public Methods
Following are the public methods that will be used by the cmd script
#### GetRandomInteger
Generates a random int32 based within the given range [min,max]
##### Signature
```
func GetRandomInteger(min int32, max int64) int32 {
}
```
#### GetRandomBlockHash
Generates a random sha256 hash based on the seed provided via command line. Throughout the code we use this to generate hashes
##### Signature
```
func GetRandomBlockHash() cipher.SHA256 {
}
```
#### GetSimulation
This method returns the globally active singleton simulation object
##### Signature
```
func GetSimulation() *Simulation {
}
```
#### InitSimulation
This method initializes the global simulation object based on command line arguments
##### Signature
```
func InitSimulation(sim *Simulation, numberOfNodes int, numberOfSubscribers int, seed int64) {
}
```

## Struct Simulation
The Simulation struct will hold all the data required for the running instance of the obelisk simulation. This struct will be maintained as singleton
### Data
```
struct Simulation {
    Ticks uint64                        // Running ticks of the application
    Nodes map[cipher.SHA256]*Node       // Node Public key to Node data structure map
    Root  *BlockRecord                  // BlockRecord hierarchy of the block records in the application
}
```
### Methods
Following are the methods supported by the Simulation Struct
#### AdvanceTicks
This method simply increments ticks held by the running simulation application. It tracks the total number of updates done on any of the nodes.
##### Signature
```
func (sim *Simulation) AdvanceTicks() {
}
```

#### PrintTotalState
This method prints the state of each of the nodes as csv
##### Signature
```
func (sim *Simulation) PrintTotalState() {
}
```

#### GenerateRandomBlockTree
This method recursively generates a random block tree for a given number of nodes and a given number of children of each node 
##### Signature
```
func (sim *Simulation) GenerateRandomBlockTree(seqNo uint64, totalNodes int64, childrenPerNode int64) *BlockRecord {
}
```

## Struct BlockRecord
The BlockRecord struct will hold data to simulate a Block
### Data
```
struct BlockRecord {
    Hash            cipher.SHA256       // Hash of the Block
    SeqNo           uint64              // SeqNo of the block
    Parent          *BlockRecord        // Pointer to the parent of the block record
    Children        []*BlockRecord      // List of children of the block record
}
```
### Methods
Following are the methods supported by the BlockRecord Struct required for this simulation
#### InitializeRandomBlock
Initializes a block based on Hash Values
##### Signature
```
func (b *BlockRecord) InitializeRandomBlock() {
}
```

## struct Node
The Node struct holds the Node information for the running simulation
### Data
```
struct Node {
    subscriptions []*Node                                // List of Nodes subscribed by the current Node
    state         map[cipher.SHA256]*NodeBlockMeta       // A mapping from BlockRecord Hash to current Node's separate copy of NodeBlockMeta
    seqNo         uint64                                 // Node's sequence number tracking the number of updates done on the node
}
```
### Methods
Following are the methods supported by the Node Struct required for this simulation
#### IntroduceBlock
- Given a block record. Add a new entry in state with key being the block record hash and value being a new NodeBlockMeta object
- Set the sequence number && ticks to 0 for NodeBlockMeta 

##### Signature
```
func (n *Node) IntroduceBlock(b *blockRecord) {
}
```

#### InitializeNodeState
Initializes the current node's state:
- Iterate through the global block record tree held by Simulation struct 
- Foreach of the block record call n.IntroduceBlock and then initialize the weight = (weight of parent) / (number of children of parent)
##### Signature
```
func (n *Node) InitializeNodeState() {
}
```


#### ValidateNodeState
  - For NodeBlockMeta entry in current node's state:
  - Call VerifyNodeBlockMeta 
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
## struct NodeBlockMeta
The NodeBlockMeta struct holds each node's individual copy of block record details.
### Data
```
struct NodeBlockMeta {
    blockRecord   *NodeBlockRecord    // The corresponding BlockRecord from the tree
    seqNo         uint64              // Here we maintain the highest seqNo among the nodes considered while syncing the states  
    ticks         uint64              // capture the ticks from the global simulation at the time when the NodeBlockMeta was synced
    weight        int32               // weight of NodeBlockMeta. This weight will be sum of the weight of the children of the blockRecord  
}
```
### Methods
Following are the methods supported by the NodeBlockMeta Struct required for this simulation
#### VerifyNodeBlockMeta
- For the given block verify that it's weight is equal to the sum of the weights of the children

##### Signature
```
func (n *NodeBlockMeta) VerifyNodeBlockMeta() {
}
```


## Simulation Run 
The simulation will be run as a command line script

### Sample Run
```console
foo@bar:~$ go run obelisk-sim.go --nodes 3 --subscribers 2 --iterations 100 --seed 123
```
### Simulation Flow
- Print each node id along with public key
- Print each node along with it's subscribers 
- For the number of iterations provided:
  - Based on the given seed, generate a random number in range (1, number of nodes)
  - Get the corresponding node based on the index of the node at above generated random number and call node.UpdateState.
