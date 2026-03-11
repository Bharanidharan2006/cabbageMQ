package queue

import (
	"errors"
	"sync"

	"github.com/Bharanidharan2006/cabbageMQ/internal/core"
)

type QueueManager struct {
	queues map[string]*RingBuffer[*core.Job]
	mu sync.RWMutex
}

type TopicConfig struct {
	Topic string
	Size int
	MaxBufferSize int
}


func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[string]*RingBuffer[*core.Job]),
	}
}

func (qm *QueueManager)CreateQueues(configs []TopicConfig) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	for _, config := range configs {
		qm.queues[config.Topic] = NewRingBuffer[*core.Job](config.Size, config.MaxBufferSize)
	}
}

func (qm *QueueManager)CreateQueue(config TopicConfig) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.queues[config.Topic] = NewRingBuffer[*core.Job](config.Size, config.MaxBufferSize)
}

func (qm *QueueManager)Push(topic string, job *core.Job) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	tqueue := qm.queues[topic]

	if tqueue == nil {
		return errors.New("The queue for the given topic doesn't exist. Create a Queue First");
	}

	if err := tqueue.Enqueue(job); err != nil {
		return err
	} else {
		return nil
	}
}

// Here the bool indicates whether there is some job to read.
// If it is false, it means there is nothing to read and the queue is empty


// The ok parameter is false also when we encounter an error. So in the calling function we first check whether there is an error
// and only then we check whether the queue is empty denoted by the ok parameter

func (qm *QueueManager)Read(topic string) (*core.Job, bool, error) {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	tqueue := qm.queues[topic]

	if tqueue == nil {
		return nil, false, errors.New("The queue for the given topic doesn't exist. Create a Queue First");
	}	

	data, ok := tqueue.Dequeue()

	if !ok {
		return nil, false, nil
	}

	return data, true, nil
}

