package database

import (
	"errors"
	"strings"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type KVStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string][]byte),
	}
}

func (k *KVStore) Set(key string, value []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.data[key] = value
}

func (k *KVStore) Get(key string) ([]byte, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	value, ok := k.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

func (k *KVStore) Delete(key string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(k.data, key)
	return nil
}

func (k *KVStore) List(prefix string) (map[string][]byte, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	result := make(map[string][]byte)
	for key, value := range k.data {
		if strings.HasPrefix(key, prefix) {
			result[key] = value
		}
	}
	return result, nil
}
