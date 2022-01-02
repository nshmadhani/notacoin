package tx

import (
	"crypto/sha256"
	"encoding/hex"
)

type MerkleTree struct {
	Root MerkleNode
}

type MerkleNode struct {
	Right *MerkleNode
	Left  *MerkleNode
	Data  []byte
}

func (tree *MerkleTree) String() string {

	return preorder(&tree.Root)

}

func preorder(node *MerkleNode) string {
	if len(node.Data) == 0 {
		return ""
	} else {
		return preorder(node.Left) + preorder(node.Right) + node.String()
	}
}

func (node *MerkleNode) String() string {
	return hex.EncodeToString(node.Data) + "\n"
}
func CreateMerkleTree(txs []*Transaction) MerkleTree {

	if len(txs)%2 != 0 {
		txs = append(txs, txs[len(txs)-1])
	}

	var nodes []*MerkleNode

	for _, tx := range txs {
		nodes = append(nodes, createMerkleNode(&MerkleNode{}, &MerkleNode{}, tx.Serialize()))
	}

	root := createTree(nodes)

	return MerkleTree{root}
}

func createTree(nodes []*MerkleNode) MerkleNode {

	if len(nodes) == 1 {
		return *nodes[0]
	}
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	var newNodes []*MerkleNode

	for i := 0; i < len(nodes); i += 2 {
		newNodes = append(newNodes, createMerkleNode(nodes[i], nodes[i+1], []byte{}))
	}
	return createTree(newNodes)
}

func createMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {

	var hash [32]byte
	if len(left.Data) == 0 && len(right.Data) == 0 {
		hash = sha256.Sum256(data)
	} else {
		data = append(left.Data, right.Data...)
		hash = sha256.Sum256(data)
	}

	return &MerkleNode{right, left, hash[:]}
}
