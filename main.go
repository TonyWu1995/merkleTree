package main

import (
	"crypto/sha256"
	"fmt"
	"merkle-tree/merkletree"
	"merkle-tree/storage"
)

func main() {
	mk, err := merkletree.Build(3)
	if err != nil {
		fmt.Println(err)
	}
	mk.Put(sha256.Sum256([]byte("123")))
	mk.Put(sha256.Sum256([]byte("1234")))
	mk.Put(sha256.Sum256([]byte("12345")))
	mk.CalcAndGetRootHash()

	proof, _ := mk.GetProof(7)

	fmt.Println(proof.Verify(mk.CalcAndGetRootHash()))
	localStorage := storage.LocalStorage{}

	localStorage.Save("./test.json", mk)
	mk1, _ := merkletree.Build(3)
	localStorage.Load("./test.json", mk1)
	p1, _ := mk1.GetProof(7)
	fmt.Println(p1.Verify(mk.CalcAndGetRootHash()))

}
