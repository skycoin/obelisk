package main

import (
	"github.com/skycoin/skycoin/src/cipher"
)

func GetRandomSHA256() cipher.SHA256 {
	return cipher.MustSHA256FromBytes(cipher.RandByte(32))
}

// Simulate a random cipher.PubKey
func GetRandomPubKey() cipher.PubKey {
	b := cipher.RandByte(33)
	p := cipher.PubKey{}
	copy(p[:], b[:])
	return p
}
