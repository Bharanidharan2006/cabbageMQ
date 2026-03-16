// TODO
// [ ] - Custom Error types
// [ ] - Changing % to AND since the size in the powers of two

package queue

import (
	"errors"
	"sync/atomic"

	"github.com/Bharanidharan2006/cabbageMQ/internal/core"
)

// Upgrade to a RWMutex after learning its advantages

type Slot struct {
	sequenceNo atomic.Uint64
	job *core.Job
}

type RingBuffer struct {
	buffer []Slot
	size uint64
	head atomic.Uint64
	_padding [56]byte
	// THis padding is essential, as the CPU, when it reads a value, it also reads its adjacent bytes(a total of 64 bytes) so to read our uint64, it 
	// reads the 8 bytes plus the remaining 56 bytes. So the head and tail or prevented from being the same cache line by adding padding.
	// What happens if they are in the same cache line, the moment head is changed and released, the processor invalidates the entire cache line including tail
	// But tail is still correct, so this padding prevents it
	tail atomic.Uint64
}

func NewRingBuffer(size uint64 ) *RingBuffer {
	buff := &RingBuffer{
		buffer : make([]Slot, size),
		size: size,
	}

	// Initialize head and tail pointers

	buff.head.Store(0)
	buff.tail.Store(0)

	// Allocate memory for each slot and then set the sequence no initially to the index

	for i := range buff.buffer {
		buff.buffer[i].sequenceNo.Store(uint64(i))
	}

	return buff
}

// Chose a fixed sized buffer because resize operation takes more time -> O(n)
// Instead return an err when the queue is full. Used by the queue to backpressure the producer

// Return an err or something? But it is guaranteed that it will be stored
func (rb *RingBuffer) Enqueue(data *core.Job) error  {

	for{
		claim := rb.head.Load()
		slotIndex := claim % rb.size // change to claim AND (size - 1) since the size is always in the power of two

		seqNo := rb.buffer[slotIndex].sequenceNo.Load()

		if seqNo == claim {
			// First claiming the slot before writing data into it. After writing data, marking the slot as ready by doing an atomic store
			swapped := rb.head.CompareAndSwap(claim, claim + 1)
			if !swapped {
				continue
			}
			rb.buffer[slotIndex].job = data
			rb.buffer[slotIndex].sequenceNo.Store(claim + 1)
			return nil
		} else if seqNo < claim { 
			// claim - seqNo < 0 => The buffer is full and return an error (backpressuring)
			return errors.New("The buffer is full!")
		} else {
			continue
		}
		// because if seqNo is ahead of the claim, it means that the slot is already acquired by another producer
			
	}
}

// Returns the following 
// 1. nil , false -> if the job is read successfully
// 2. error, false -> check the error first, so an error occured
// 3. nil, true -> queue is empty

func (rb *RingBuffer) Dequeue() *core.Job {

	for {
		claim := rb.tail.Load()
		slotIndex := claim % rb.size

		seqNo := rb.buffer[slotIndex].sequenceNo.Load()
		// First the slot is given to a particular consumer and then consume it

		if seqNo == (claim + 1) {
			// It fails if a another consumed has already consumed the job
			swapped := rb.tail.CompareAndSwap(claim, claim + 1)
			if !swapped {
				continue
			}
			job := rb.buffer[slotIndex].job
			rb.buffer[slotIndex].sequenceNo.Store(claim + rb.size)
			return job
		} else if seqNo < (claim + 1) { // if the seqNo is behind the tail, it means that the buffer is empty
			return nil
		}else {
			continue
		}
	}
}

