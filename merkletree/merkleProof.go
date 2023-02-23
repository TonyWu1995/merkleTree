package merkletree

import (
	"bytes"
	"crypto/sha256"
)

type MerkleProof struct {
	index int
	Val   [][32]byte
}

func (proof *MerkleProof) Verify(rootHash [32]byte) bool {

	index := proof.index
	var parentIndex int
	for i := 0; i < len(proof.Val)-1; i = i + 2 {
		concatHash := append(proof.Val[i][:], proof.Val[i+1][:]...)
		digest := sha256.Sum256(concatHash)
		if parentIndex/2 == 1 {
			parentIndex = i + 2
		} else {
			parentIndex = i + 2 + (index/2)%2
		}
		if bytes.Compare(proof.Val[parentIndex][:], digest[:]) != 0 {
			return false
		}
		index = index / 2
	}

	return bytes.Compare(proof.Val[len(proof.Val)-1][:], rootHash[:]) == 0
}
