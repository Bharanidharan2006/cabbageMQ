package queue

import (
	"errors"
	"sync"
)

// Upgrade to a RWMutex after learning its advantages

type RingBuffer[T any] struct {
	buffer []T
	size int
	mu sync.Mutex
	lastInsert int
	nextRead int
	maxBufferSize int
	currentSize int
}

func NewRingBuffer[T any](size int, maxBufferSize int) *RingBuffer[T] {
	if size == 0 || maxBufferSize == 0 {
		return nil
	}
	return &RingBuffer[T]{
		buffer : make([]T, size),
		size: size,
		lastInsert: -1,
		nextRead: -1,
		maxBufferSize: maxBufferSize,
		currentSize: 0,
	}
}

// Return an err or something? But it is guaranteed that it will be stored
func (rb *RingBuffer[T]) Enqueue(data T) error  {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.currentSize == rb.size {
		err := rb.resize()
		if err != nil {
			return err
		}
	}

	if rb.lastInsert == -1 {
		rb.lastInsert = 0
		rb.nextRead = 0
	}else {
		rb.lastInsert = (rb.lastInsert + 1) % rb.size
	}
	
	rb.buffer[rb.lastInsert] = data

	rb.currentSize++

	return nil
}

func (rb *RingBuffer[T]) Dequeue() (T, bool) {

	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.nextRead == -1 {
		var zero T
		return  zero, false
	}

	rb.currentSize--
	data := rb.buffer[rb.nextRead]
	
	if rb.currentSize == 0 {
		rb.lastInsert = -1
		rb.nextRead = -1
	} else {
		rb.nextRead = (rb.nextRead + 1) % rb.size
	}

	
	return data, true
}

func (rb *RingBuffer[T]) resize() error {
	newSize := rb.size * 2

	if newSize > rb.maxBufferSize {
		return errors.New("Job queue has reached the maximum buffer length")
	}

	newBuffer := make([]T, newSize)
	newLastInsert := -1

	for {
		newLastInsert++
		newBuffer[newLastInsert] = rb.buffer[rb.nextRead]
		rb.nextRead = (rb.nextRead + 1) % rb.size
		if rb.nextRead == (rb.lastInsert + 1) % rb.size {
			break;
		}
	}

	rb.nextRead = 0
	rb.lastInsert = newLastInsert
	rb.size = newSize

	rb.buffer = newBuffer

	return nil
}