// TODOS
// [ ] - Convert map to array as the topics are generated only at the startup. Consequently, remove CreateQueue Function
// [ ] - To make CreateQueues callabe only once, add some sort of guard

package queue

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Bharanidharan2006/cabbageMQ/internal/core"
)

type QueueManager struct {
	queues map[string]*RingBuffer
	mu sync.RWMutex
}

type TopicConfig struct {
	Topic string
	Size uint64
}


func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[string]*RingBuffer),
	}
}

// Add Queues only at the startup

func (qm *QueueManager)CreateQueues(configs []TopicConfig) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	for _, config := range configs {
		if (config.Size & (config.Size - 1)) == 0 ||  config.Size == 0  {
			return errors.New(fmt.Sprintf("the size of the queue for the topic %s is not a power of 2", config.Topic))
		}
	}
	for _, config := range configs {	
		qm.queues[config.Topic] = NewRingBuffer(config.Size)
	}
	return nil
}

// func (qm *QueueManager)CreateQueue(config TopicConfig) {
// 	qm.mu.Lock()
// 	defer qm.mu.Unlock()
// 	qm.queues[config.Topic] = NewRingBuffer(config.Size)
// }

func (qm *QueueManager)Push(topic string, job *core.Job) error {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	tqueue := qm.queues[topic]

	if tqueue == nil {
		return errors.New(fmt.Sprintf("the queue for the given topic: %s doesn't exist", topic));
	}

	return  tqueue.Enqueue(job)
}

// Here the bool indicates whether there is some job to read.
// If it is false, it means there is nothing to read and the queue is empty


// nil, nil means the job is empty

func (qm *QueueManager)Read(topic string) (*core.Job, error) {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	tqueue := qm.queues[topic]

	if tqueue == nil {
		return nil, errors.New(fmt.Sprintf("the queue for the given topic: %s doesn't exist", topic));
	}	

	data := tqueue.Dequeue()

	if data == nil {
		return nil, nil // no error occured and at the same time, job is nil. So the buffer is empty
	}

	return data, nil
}

