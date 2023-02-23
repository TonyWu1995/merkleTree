package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type node struct {
	Id    int
	Val   [32]byte
	Left  *node
	Right *node
}

type MerkleTree struct {
	Node   []node
	Height int
	Index  int
}

func Build(height int) (*MerkleTree, error) {
	if height < 1 {
		return nil, errors.New("height < 1")
	}
	nodes := make([]node, 1<<height)
	for i := 1<<height - 1; i > 0; i-- {
		if i >= (1 << (height - 1)) {
			nodes[i] = node{
				i, [32]byte{}, nil, nil,
			}
		} else {
			nodes[i] = node{i, [32]byte{}, &nodes[i*2], &nodes[i*2+1]}
		}
	}
	return &MerkleTree{
		Height: height,
		Node:   nodes,
		Index:  0,
	}, nil
}

func (mk *MerkleTree) GetRootHash() [32]byte {
	return mk.Node[1].Val
}
func (mk *MerkleTree) GetRootHashString() string {
	rootHash := mk.GetRootHash()
	return hex.EncodeToString(rootHash[:])
}

func (node *node) put(hash [32]byte) {
	node.Val = hash
}

func (mk *MerkleTree) Put(hash [32]byte) int {
	i := mk.Index + (1 << (mk.Height - 1))
	mk.Node[i].put(hash)
	mk.Index++
	return i
}

func (mk *MerkleTree) CalcAndGetRootHash() [32]byte {
	for i := 1 << (mk.Height - 1); i < 1<<mk.Height-1; i = i + 2 {
		res := append(mk.Node[i].Val[:], mk.Node[i+1].Val[:]...)
		sha256.Sum256(res)
		mk.Node[i/2].Val = sha256.Sum256(res)
	}
	mk.Node[1].Val = sha256.Sum256(append(mk.Node[2].Val[:], mk.Node[3].Val[:]...))
	return mk.GetRootHash()
}

func (mk *MerkleTree) GetProof(leafNodeIndex int) (*MerkleProof, error) {
	if leafNodeIndex < (1<<mk.Height - 1) {
		return nil, errors.New("leafNode less than setting")
	}
	var proof [][32]byte
	index := leafNodeIndex
	if index%2 == 1 {
		index--
	}
	for i := index; i > 1; i = i / 2 {
		if i%2 == 0 {
			proof = append(proof, mk.Node[i].Val)
			proof = append(proof, mk.Node[i+1].Val)
		} else {
			proof = append(proof, mk.Node[i-1].Val)
			proof = append(proof, mk.Node[i].Val)

		}
	}

	proof = append(proof, mk.Node[1].Val)
	return &MerkleProof{
		index: leafNodeIndex,
		Val:   proof,
	}, nil
}
