package store

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/cespare/xxhash/v2"
)

type StoreManager struct {
	stores  []*Store
	active  bool
	walChan chan<- string
}

func NewStoreManger(numShards uint8, walChan chan<- string) *StoreManager {
	stores := make([]*Store, numShards)

	for i := uint8(0); i < numShards; i++ {
		store := NewStore(walChan)
		stores[i] = store
	}
	return &StoreManager{
		stores:  stores,
		active:  false,
		walChan: walChan,
	}
}

func (storeManager *StoreManager) Run(ctx context.Context) {
	go func() {
		storeManager.start(ctx)
	}()
}

func (storeManager *StoreManager) start(ctx context.Context) {
	sharedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, store := range storeManager.stores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Start(sharedCtx)
		}()
	}
	println("Started StoreManager")
	storeManager.active = true
	wg.Wait()
}

func (storeManager *StoreManager) GetShard(key string) (*Store, error) {
	if !storeManager.active {
		println("Please start the Shard Manager First!")
		return nil, fmt.Errorf("ERR Please start the shard manager first")
	}
	index := xxhash.Sum64String(key) % uint64((len(storeManager.stores)))
	println("Operation in Shard number : " + strconv.FormatUint(index, 10))
	return storeManager.stores[index], nil
}
