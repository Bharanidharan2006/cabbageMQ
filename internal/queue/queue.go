package queue

import "github.com/Bharanidharan2006/cabbageMQ/internal/core"

type TopicConfig struct {
	Topic string
	Size int
	MaxBufferSize int
}

var QueuesMap map[string]*RingBuffer[core.Job] = make(map[string]*RingBuffer[core.Job])

func CreateQueues(configs []TopicConfig) {
	for _, config := range configs {
		QueuesMap[config.Topic] = NewRingBuffer[core.Job](config.Size, config.MaxBufferSize)
	}
}

