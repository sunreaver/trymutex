package trymutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	lockedFlag   int32 = 1
	unlockedFlag int32 = 0
)

// TryMutex 具有探测功能的sync.Mutex
type TryMutex struct {
	in *sync.Mutex
}

// NewTryMutex NewTryMutex
func NewTryMutex() *TryMutex {
	return NewTryMutexWithSyncMutex(&sync.Mutex{})
}

// NewTryMutexWithSyncMutex NewTryMutexWithSyncMutex
func NewTryMutexWithSyncMutex(m *sync.Mutex) *TryMutex {
	return &TryMutex{
		in: m,
	}
}

// Lock 同sync.Lock
func (m *TryMutex) Lock() {
	m.in.Lock()
}

// Unlock 同sync.Unlock
func (m *TryMutex) Unlock() {
	m.in.Unlock()
}

// TryLock 探测是否可lock，如果可以则lock并返回true；反之相反
func (m *TryMutex) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m.in)), unlockedFlag, lockedFlag) {
		return true
	}
	return false
}

// TryUnLock 探测是否可unlock，如果可以则unlock并返回true；反之相反
func (m *TryMutex) TryUnLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m.in)), lockedFlag, unlockedFlag) {
		return true
	}
	return false
}

// IsLocked 返回是否在lock状态
func (m *TryMutex) IsLocked() bool {
	if atomic.LoadInt32((*int32)(unsafe.Pointer(m.in))) == lockedFlag {
		return true
	}
	return false
}
