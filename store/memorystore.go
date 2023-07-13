package store

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/wdantuma/signalk-server-go/signalk"
)

type memoryStore struct {
	values   map[string]*Value
	keyIndex []string
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{values: make(map[string]*Value), keyIndex: make([]string, 0)}
}

func (s *memoryStore) Put(key string, timestamp int64, vessel string, path string, source *signalk.Source, value interface{}) {
	storeValue := &Value{Vessel: vessel, Path: path, Value: value, Source: source, LastChange: timestamp}
	_, valueExists := s.values[key]
	s.values[key] = storeValue
	if !valueExists {
		insertIndex := sort.SearchStrings(s.keyIndex, key)
		if len(s.keyIndex) == insertIndex {
			s.keyIndex = append(s.keyIndex, key)
		} else {
			s.keyIndex = append(s.keyIndex[:insertIndex+1], s.keyIndex[insertIndex:]...)
			s.keyIndex[insertIndex] = key
		}
	}
}

func (s *memoryStore) Get(key string) (*Value, bool) {
	v, ok := s.values[key]
	return v, ok
}

func (s *memoryStore) GetList(key string) []*Value {
	values := make([]*Value, 0)
	startIndex := sort.SearchStrings(s.keyIndex, key)
	for i := startIndex; i < len(s.keyIndex); i++ {
		if strings.Index(s.keyIndex[i], key) == 0 {
			k := s.keyIndex[i]
			v := s.values[k]
			values = append(values, v)
		}
	}
	return values
}

func (s *memoryStore) Store(input <-chan signalk.DeltaJson) <-chan signalk.DeltaJson {
	output := make(chan signalk.DeltaJson)
	go func() {
		for delta := range input {
			for _, update := range delta.Updates {
				for _, value := range update.Values {
					key := fmt.Sprintf("%s/%s", *delta.Context, value.Path)
					timeStamp, err := time.Parse(signalk.TIME_FORMAT, string(*update.Timestamp))
					if err != nil {
						timeStamp = time.Now()
					}
					s.Put(key, timeStamp.Unix(), *delta.Context, value.Path, update.Source, value.Value)
				}
			}
			output <- delta
		}
	}()

	return output
}
