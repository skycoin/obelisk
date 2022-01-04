package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
)

const DEFAULT_NODES = 3
const DEFAULT_SUBSCRIBERS = 2
const DEFAULT_ITERATIONS = 10
const DEFAULT_BLOCK_TREE_BLOCK_RECORD_COUNT = 5
const DEFAULT_BLOCK_TREE_CHILDREN_COUNT = 2

const MIN_NODES = 3
const MIN_SUBSCRIBERS = 1
const MIN_ITERATIONS = 1
const MIN_BLOCK_TREE_BLOCK_RECORD_COUNT = 3
const MIN_BLOCK_TREE_CHILDREN_COUNT = 2

func main() {

	// Required Arguments
	nodeCount := flag.Int("nodes", DEFAULT_NODES, fmt.Sprintf("[Required] Number of nodes to consider for simulation. Min Value: %d", MIN_NODES))
	subscriberCount := flag.Int("subscribers", DEFAULT_SUBSCRIBERS, fmt.Sprintf("[Required] Number of subscribers per node. Must be less than number of nodes. Min Value: %d", MIN_SUBSCRIBERS))

	// Optional Arguments
	showHelp := flag.Bool("help", false, "Show Help")
	verboseMode := flag.Bool("verbose", false, "Run in Verbose Mode")
	defaultSeed := int64(1); // time.Now().UTC().UnixNano();
	seed := flag.Int64("seed", defaultSeed, "Seed to use while running the simulation. Must be a valid integer > 0")
	iterations := flag.Int("iterations", DEFAULT_ITERATIONS, fmt.Sprintf("Number of iterations to run this simulation. Min Value: %d", MIN_ITERATIONS))
	blockRecordCount := flag.Int("block-record-count", DEFAULT_BLOCK_TREE_BLOCK_RECORD_COUNT, fmt.Sprintf("Total Number of Blocks in Root Block Tree. Min Value: %d", MIN_BLOCK_TREE_BLOCK_RECORD_COUNT))
	childrenPerBlock := flag.Int("children-per-block", DEFAULT_BLOCK_TREE_CHILDREN_COUNT, fmt.Sprintf("Max Number of Children Per Block in Root Block Tree. Min Value: %d", MIN_BLOCK_TREE_CHILDREN_COUNT))

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if flag.NFlag() < 2 {
		fmt.Printf("\nNot enough arguments!!\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *nodeCount < MIN_NODES {
		log.Printf("Invalid Value for nodes: %d (Must be a valid integer with minimum value: %d)", *nodeCount, MIN_NODES)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *subscriberCount < MIN_SUBSCRIBERS || *subscriberCount >= *nodeCount {
		log.Printf("Invalid Value for subscribers: %d (Must be a valid integer with minimum value: %d and must be less than nodes)", *subscriberCount, MIN_SUBSCRIBERS)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *blockRecordCount < MIN_BLOCK_TREE_BLOCK_RECORD_COUNT {
		log.Printf("Invalid Value for blockRecordCount: %d (Must be a valid integer with minimum value: %d)", *blockRecordCount, MIN_BLOCK_TREE_BLOCK_RECORD_COUNT)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *childrenPerBlock < MIN_BLOCK_TREE_CHILDREN_COUNT || *childrenPerBlock >= *blockRecordCount {
		log.Printf("Invalid Value for childrenPerBlock: %d (Must be a valid integer with minimum value: %d and must be less than %d)", *childrenPerBlock, MIN_BLOCK_TREE_CHILDREN_COUNT, *blockRecordCount)
		flag.PrintDefaults()
		os.Exit(1)
	}

	rand.Seed(*seed)

	simulation := GetSimulation()
	simulation.InitSimulation(*blockRecordCount, *childrenPerBlock, *nodeCount, *subscriberCount, *iterations, *verboseMode)
	if err := simulation.RunSimulation(); err != nil {
		fmt.Printf("\nSimulation Failed with Error: %s\n", err.Error())
		simulation.PrintAllNodes()
		os.Exit(1)
	}
}
