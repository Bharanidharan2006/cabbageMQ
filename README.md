# CabbageMQ

CabbageMQ is a lightweight, RabbitMQ-inspired distributed message broker built entirely from scratch in Go. It handles high-throughput task scheduling using a custom length-prefixed TCP protocol, a lock-free ring buffer for O(1) in-memory routing, and an append-only Write-Ahead Log (WAL) to guarantee zero job loss during server crashes.

## Lock Free Ring Buffer

Using a lock to synchronize among multiple producers and multiple consumers is to slow. So implemented a simple lock-free algorithm using CAS.

### Struct

```go
type Slot struct {
    job *core.Job
    sequenceNo atomic.Uint64
}
type RingBuffer struct {
    buffer []Slot
    head atomic.Uint64
    tail atomic.Uint64
    _padding [56]byte // Why a 56 byte padding? Explained below
    size uint64 // must be a power of two
}
```

### Methods

```go
 RingBuffer.Enqueue(data *core.Job) error
```

- Takes in a pointer to a job and returns nil, if successfully enqueued and an error if the buffer is full

```go
 RingBuffer.Dequeue() *core.Job
```

- Returns nil if the buffer is empty or else returns the job

### Why the 56 byte padding?

- Modern CPUs, whenever they load some data in memory, they also load the adjacent bytes (in the assumption that the adjacent bytes will be needed in the subsequent operations like in the case of looping through an error).
- Without the padding, when the CPU loads the head, it also loads the tail and both of them are put together in the same cache line.
- Now when the head(or tail) is updated and release store is called, the cache line gets invalid and the tail(or head) is read again on acquire even though its value is not changed.
- This causes undetectable performance issues later down the line so the padding is necessary
