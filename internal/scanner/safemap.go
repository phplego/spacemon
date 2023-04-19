package scanner

import (
	"encoding/json"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: make(map[K]V)}
}

func (sm *SafeMap[K, V]) Get(key K) (value V, ok bool) {
	sm.RLock()
	defer sm.RUnlock()
	value, ok = sm.m[key]
	return
}

func (sm *SafeMap[K, V]) Set(key K, value V) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap[K, V]) Items() map[K]V {
	sm.RLock()
	defer sm.RUnlock()
	copied := make(map[K]V, len(sm.m))
	for k, v := range sm.m {
		copied[k] = v
	}
	return copied
}
func (sm *SafeMap[K, V]) MarshalJSON() ([]byte, error) {
	sm.RLock()
	defer sm.RUnlock()
	return json.Marshal(sm.m)
}

func (sm *SafeMap[K, V]) UnmarshalJSON(data []byte) error {
	sm.Lock()
	defer sm.Unlock()
	return json.Unmarshal(data, &sm.m)
}
