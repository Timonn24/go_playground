package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type KeyValueStore struct {
	data sync.Map
}

// sync.Map doesn't support Marshal interface directly
func (kv *KeyValueStore) Bytes() ([]byte, error) {
	tempMap := make(map[string]any)
	kv.data.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			tempMap[k] = value
		}
		return true
	})

	jsonData, err := json.Marshal(tempMap)
	if err != nil {
		return nil, fmt.Errorf("error marshalling: %v", err)
	}
	return jsonData, nil
}

func (kv *KeyValueStore) FromBytes(data []byte) error {
	tempMap := make(map[string]any)
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	for k, v := range tempMap {
		kv.data.Store(k, v)
	}
	return nil
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data: sync.Map{},
	}
}

func (s *KeyValueStore) Set(key string, value string) {
	s.data.Store(key, value)
}

func (s *KeyValueStore) Get(key string) (string, bool) {
	value, ok := s.data.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (s *KeyValueStore) Delete(key string) {
	s.data.Delete(key)
}
