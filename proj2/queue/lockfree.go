package queue

import (
    "sync/atomic"
    "unsafe"
)

type Request struct {
    Command   string
    ID        int
    Body      string
    Timestamp float64
}

type node struct {
    task *Request
    next unsafe.Pointer
}

// LockFreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
    head unsafe.Pointer
    tail unsafe.Pointer
}

// NewLockFreeQueue creates and initializes a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
    dummy := &node{}
    queue := &LockFreeQueue{
        head: unsafe.Pointer(dummy),
        tail: unsafe.Pointer(dummy),
    }
    return queue
}

// Enqueue adds a series of Request to the queue
func (queue *LockFreeQueue) Enqueue(task *Request) {
    newNode := &node{task: task}
    for {
        tail := (*node)(atomic.LoadPointer(&queue.tail))
        next := (*node)(atomic.LoadPointer(&tail.next))
        if tail == (*node)(atomic.LoadPointer(&queue.tail)) {
            if next == nil {
                if atomic.CompareAndSwapPointer(&tail.next, nil, unsafe.Pointer(newNode)) {
                    atomic.CompareAndSwapPointer(&queue.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
                    return
                }
            } else {
                atomic.CompareAndSwapPointer(&queue.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
            }
        }
    }
}

// Dequeue removes a Request from the queue
func (queue *LockFreeQueue) Dequeue() *Request {
    for {
        head := (*node)(atomic.LoadPointer(&queue.head))
        tail := (*node)(atomic.LoadPointer(&queue.tail))
        next := (*node)(atomic.LoadPointer(&head.next))
        if head == (*node)(atomic.LoadPointer(&queue.head)) {
            if head == tail {
                if next == nil {
                    return nil
                }
                atomic.CompareAndSwapPointer(&queue.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
            } else {
                task := next.task
                if atomic.CompareAndSwapPointer(&queue.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
                    return task
                }
            }
        }
    }
}

func (queue *LockFreeQueue) IsEmpty() bool {
    head := (*node)(atomic.LoadPointer(&queue.head))
    next := (*node)(atomic.LoadPointer(&head.next))
    return next == nil
}