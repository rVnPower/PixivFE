package kmutex

import "sync"

// Map, Refencence-counted by ID (any).
type Kmutex struct {
	Map sync.Map
}

// Create new Kmutex
func New() *Kmutex {
	return &Kmutex{}
}

// decrement ID ref count
// Returns: ref count after
func (km *Kmutex) Unlock(key any) uint64 {
	for {
		actual, ok := km.Map.Load(key)
		if !ok {
			panic("impossible! memory corruption?")
		}
		if actual.(uint64) == 1 {
			deleted := km.Map.CompareAndDelete(key, actual.(uint64))
			if deleted {
				return 0
			}
		} else {
			after := actual.(uint64) - 1
			swapped := km.Map.CompareAndSwap(key, actual.(uint64), after)
			if swapped {
				return after
			}
		}
	}
}

// increment ID ref count
// Returns: ref count after
func (km *Kmutex) Lock(key any) uint64 {
	for {
		actual, loaded := km.Map.LoadOrStore(key, uint64(1))
		if !loaded {
			return 1
		}
		after := actual.(uint64) + 1
		swapped := km.Map.CompareAndSwap(key, actual.(uint64), after)
		if swapped {
			return after
		}
	}
}
