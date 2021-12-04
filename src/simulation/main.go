package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const DEFAULT_NODES = 3
const DEFAULT_SUBSCRIBERS = 2
const MIN_NODES = 3
const MIN_SUBSCRIBERS = 1
const DEFAULT_ITERATIONS = 3
const MIN_ITERATIONS = 1
const BLOCK_TREE_BLOCK_RECORD_COUNT = 5
const BLOCK_TREE_CHILDREN_COUNT = 2

func main() {

	showHelp := flag.Bool("help", false, "Show Help")

	// Required Arguments
	nodeCount := flag.Int("nodes", DEFAULT_NODES, fmt.Sprintf("[Required] Number of nodes to consider for simulation. Min Value: %d", MIN_NODES))
	subscriberCount := flag.Int("subcribers", DEFAULT_SUBSCRIBERS, fmt.Sprintf("[Required] Number of subscribers per node. Must be less than nodes. Min Value: %d", MIN_SUBSCRIBERS))

	// Optional Arguments
	verboseMode := flag.Bool("verbose", false, "Run in Verbose Mode")
	seed := flag.Int64("seed", time.Now().UTC().UnixNano(), "Seed to use while running the simulation. Must be a valid integer > 0")
	iterations := flag.Int("iterations", DEFAULT_ITERATIONS, fmt.Sprintf("Number of iterations to run this simulation. Min Value: %d", MIN_ITERATIONS))

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

	rand.Seed(*seed)

	simulation := GetSimulation()
	simulation.InitSimulation(BLOCK_TREE_BLOCK_RECORD_COUNT, BLOCK_TREE_CHILDREN_COUNT, *nodeCount, *subscriberCount, *iterations, *verboseMode)
	simulation.RunSimulation()
}
