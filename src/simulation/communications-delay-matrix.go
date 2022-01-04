package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type CommunicationsDelayMatrix struct {
	matrix map[cipher.PubKey]map[cipher.PubKey]int
}

func (cdm *CommunicationsDelayMatrix) InitializeCommuncationsDelayMatrix(nodes []*Node) {
	cdm.matrix = map[cipher.PubKey]map[cipher.PubKey]int{}
	for _, nodeI := range nodes {
		cdm.matrix[nodeI.pubKey] = map[cipher.PubKey]int{};
		for _, nodeJ := range nodes {
			// Right now we assume communication delay to be 1 for each pair of nodes
			cdm.matrix[nodeI.pubKey][nodeJ.pubKey] = 1
		}
	}
}
