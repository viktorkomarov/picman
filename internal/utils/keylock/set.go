package keylock

import (
	"hash/fnv"
	"sync"
)

type HashedMutex interface {
	Lock(key string)
	Unlock(key string)
}

type hashedMutex struct {
	mutexes []sync.Mutex
}

func NewHashedMutex(n int) HashedMutex {
	return &hashedMutex{
		mutexes: make([]sync.Mutex, max(n, 1)),
	}
}

func (h *hashedMutex) Lock(key string) {
	h.mutexes[idByKey(key, len(h.mutexes))].Lock()
}

func (h *hashedMutex) Unlock(key string) {
	h.mutexes[idByKey(key, len(h.mutexes))].Unlock()
}

func idByKey(key string, count int) int {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % count
}
