package store

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	storeMap map[string]Value
	mtx      sync.Mutex
	walChan  chan<- string
}

func (s *Store) Start(ctx context.Context) {

	ticker := time.NewTicker(DEFAULT_DELETE_CRON)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.deleteCronTask()
		case <-ctx.Done():
			return
		}
	}

}

func (s *Store) deleteCronTask() {
	sampleSize := 20
	deleted := 0

	for k, v := range s.storeMap {
		if v.CheckIsExpired() {
			delete(s.storeMap, k)
			deleted++
		}
		if deleted >= sampleSize {
			break
			println("Deleted " + strconv.FormatInt(int64(deleted), 10) + " Keys")
		}
	}
}

func (s *Store) Get(key string) *Value {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	val, ok := s.storeMap[key]
	if !ok {
		return nil
	}
	if val.CheckIsExpired() {
		s.storeMap[key] = val
		return nil
	}
	s.walChan <- "GET " + key
	return &val
}

func (s *Store) Set(key string, value Value) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.storeMap[key] = value
}

func (s *Store) Delete(key string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	val, ok := s.storeMap[key]
	if !ok {
		return false
	}
	if val.CheckIsExpired() {
		return false
	}
	val.SetIsExpired()
	return true
}

func (s *Store) GetAllKeys() []string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var key []string = make([]string, len(s.storeMap))
	for k, v := range s.storeMap {
		if v.CheckIsExpired() {
			continue
		}
		key = append(key, k)
	}
	return key
}

func NewStore(walChan chan<- string) *Store {
	return &Store{
		storeMap: make(map[string]Value),
		walChan:  walChan,
	}
}
