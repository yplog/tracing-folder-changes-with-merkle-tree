package main

import (
	"encoding/hex"
	"fmt"
)

type MerkleNode struct {
	Left   *MerkleNode
	Right  *MerkleNode
	Hash   []byte
	Data   []byte
	IsLeaf bool
}

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleProofStep struct {
	Hash   []byte
	IsLeft bool
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	var nodeHash []byte
	var isLeaf bool

	if left == nil && right == nil {
		nodeHash = hash(data)
		isLeaf = true
	} else {
		combined := append(left.Hash, right.Hash...)
		nodeHash = hash(combined)
	}

	return &MerkleNode{
		Left:   left,
		Right:  right,
		Hash:   nodeHash,
		Data:   data,
		IsLeaf: isLeaf,
	}
}

func NewMerkleTree(dataBlocks [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	for _, data := range dataBlocks {
		node := NewMerkleNode(nil, nil, data)
		nodes = append(nodes, node)
	}

	for len(nodes) > 1 {
		var newLevel []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				newLevel = append(newLevel, NewMerkleNode(nodes[i], nodes[i], nil))
			} else {
				newLevel = append(newLevel, NewMerkleNode(nodes[i], nodes[i+1], nil))
			}
		}

		nodes = newLevel
	}

	return &MerkleTree{Root: nodes[0]}
}

func (t *MerkleTree) GenerateProof(data []byte) ([]MerkleProofStep, bool) {
	var path []MerkleProofStep
	var found bool

	var walk func(node *MerkleNode) bool
	walk = func(node *MerkleNode) bool {
		if node == nil {
			return false
		}
		if node.IsLeaf && string(node.Data) == string(data) {
			return true
		}
		if walk(node.Left) {
			path = append(path, MerkleProofStep{Hash: node.Right.Hash, IsLeft: false})
			return true
		}
		if walk(node.Right) {
			path = append(path, MerkleProofStep{Hash: node.Left.Hash, IsLeft: true})
			return true
		}
		return false
	}

	found = walk(t.Root)
	return path, found
}

func VerifyProof(data []byte, proof []MerkleProofStep, rootHash []byte) bool {
	current := hash(data)
	for _, step := range proof {
		if step.IsLeft {
			current = hash(append(step.Hash, current...))
		} else {
			current = hash(append(current, step.Hash...))
		}
	}
	return string(current) == string(rootHash)
}

func printMerkleTree(node *MerkleNode, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	hashStr := hex.EncodeToString(node.Hash)
	label := "Internal"
	if node.IsLeaf {
		label = fmt.Sprintf("Leaf: %q", string(node.Data))
	}

	branch := "├──"
	if !isLeft {
		branch = "└──"
	}

	fmt.Printf("%s%s [%s] %s\n", prefix, branch, label, hashStr)

	childPrefix := prefix + "│   "
	if !isLeft {
		childPrefix = prefix + "    "
	}

	printMerkleTree(node.Left, childPrefix, true)
	printMerkleTree(node.Right, childPrefix, false)
}
