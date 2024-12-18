// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import "sync"



const MaxReaders = 32



type RWLock struct {
    mu       *sync.Mutex
    cond     *sync.Cond
    readers  int    
    writers  int    
    waitingWriters int 
}

func NewRWLock() *RWLock {
    rw := &RWLock{
        mu: &sync.Mutex{},
    }
    rw.cond = sync.NewCond(rw.mu)
    return rw
}

func (rw *RWLock) Lock() {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    rw.waitingWriters++ 

    for rw.readers > 0 || rw.writers > 0 {
        rw.cond.Wait()
    }

    rw.waitingWriters-- 
    rw.writers = 1      
}

func (rw *RWLock) Unlock() {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    rw.writers = 0
    rw.cond.Broadcast() 
}

func (rw *RWLock) RLock() {
    rw.mu.Lock()
    defer rw.mu.Unlock()

    for rw.writers > 0 || rw.waitingWriters > 0 || rw.readers >= MaxReaders {
        rw.cond.Wait()
    }

    rw.readers++
}

func (rw *RWLock) RUnlock() {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    rw.readers--
    if rw.readers == 0 {
        rw.cond.Broadcast() 
    }
}
