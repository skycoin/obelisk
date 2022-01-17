package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

// CommunicationsDelayMatrix: tracks the tick delay when sending messages between two nodes
type CommunicationsDelayMatrix struct {
	matrix map[cipher.PubKey]map[cipher.PubKey]int
}

// InitializeCommunicationsDelayMatrix: Initializes the CommunicationsDelayMatrix with ideally random values
// Right now we have kept all the values as 1
func (cdm *CommunicationsDelayMatrix) InitializeCommunicationsDelayMatrix(nodes []*Node) {
	cdm.matrix = map[cipher.PubKey]map[cipher.PubKey]int{}
	for _, nodeI := range nodes {
		cdm.matrix[nodeI.pubKey] = map[cipher.PubKey]int{};
		for _, nodeJ := range nodes {
			// Right now we assume communication delay to be 1 for each pair of nodes
			cdm.matrix[nodeI.pubKey][nodeJ.pubKey] = 1
		}
	}
}
