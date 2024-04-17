package httpjson

import (
	"net/http/httptest"
	"sync"
	"testing"
)

// Write a struct that represents a Key Value store that is safe for concurrent use.
type KV struct {
	lock sync.RWMutex
	data map[string]any
}

func NewKV() *KV {
	return &KV{
		lock: sync.RWMutex{},
		data: make(map[string]any),
	}
}

func (kv *KV) Get(key string) any {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	return kv.data[key]
}

func (kv *KV) Set(key string, value any) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data[key] = value
}

type Step struct {
	RequestTemplatePath string
	PreHook             func(t *testing.T, kv *KV) (vars map[string]any)
	PostHook            func(t *testing.T, kv *KV, w *httptest.ResponseRecorder)
}
