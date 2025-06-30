package util

import "hot_key/model/key"

type BlockingQueue struct {
	queue chan key.HotKeyModel
}

var (
	BlockingQueueSize = 1000000
	BlQueue           *BlockingQueue
)

func init() {
	BlQueue = NewBlockingQueue(BlockingQueueSize)
}

func NewBlockingQueue(size int) *BlockingQueue {
	return &BlockingQueue{
		queue: make(chan key.HotKeyModel, size),
	}
}

func (q *BlockingQueue) Put(item key.HotKeyModel) {
	q.queue <- item
}

func (q *BlockingQueue) Take() key.HotKeyModel {
	return <-q.queue
}
