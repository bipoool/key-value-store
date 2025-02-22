package kvpApi

import (
	"key-value-store/store"
	"time"
)

type kvpHelper struct {
	storeManager *store.StoreManager
}

func NewKvp(numShards uint8, storeManager *store.StoreManager) *kvpHelper {
	return &kvpHelper{
		storeManager: storeManager,
	}
}

func (k *kvpHelper) Set(key string, value string, ttl uint64) error {
	currTime := time.Now().Unix()
	val := store.NewValue()
	val.SetValue(value)
	if ttl != 0 {
		val.SetTtl(uint64(currTime) + ttl)
	} else {
		val.SetTtl(0)
	}
	store, err := k.storeManager.GetShard(key)
	if err != nil {
		return err
	}
	store.Set(key, *val)
	return nil
}

func (k *kvpHelper) Get(key string) (*store.Value, error) {
	store, err := k.storeManager.GetShard(key)
	if err != nil {
		return nil, err
	}
	return store.Get(key), nil
}

func (k *kvpHelper) Delete(key string) error {
	store, err := k.storeManager.GetShard(key)
	if err != nil {
		return err
	}
	store.Delete(key)
	return nil
}
