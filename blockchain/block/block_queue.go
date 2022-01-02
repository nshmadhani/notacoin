package block

import (
	"container/heap"
)

type BlockQueue []*Block

func (pq BlockQueue) Len() int { return len(pq) }

func (pq BlockQueue) Less(i, j int) bool {
	return pq[i].Less(pq[j])
}

func (pq BlockQueue) Swap(i, j int) {
	*pq[i], *pq[j] = *pq[j], *pq[i]
}

func (pq *BlockQueue) Push(x interface{}) {
	item := x.(*Block)
	*pq = append(*pq, item)
}

func (pq *BlockQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (pq *BlockQueue) PushBlock(block *Block) {
	heap.Push(pq, block)
}

func (pq *BlockQueue) PopBlock() *Block {
	if len(*pq) == 0 {
		return nil
	}
	return heap.Pop(pq).(*Block)
}

func (pq *BlockQueue) Peek() Block {

	blockPtr := pq.PopBlock()
	if blockPtr == nil {
		return Block{}
	}

	block := *blockPtr

	pq.PushBlock(blockPtr)

	return block
}

func NewBlockQueue() BlockQueue {
	var pq BlockQueue
	heap.Init(&pq)
	return pq
}
