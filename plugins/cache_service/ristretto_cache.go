package cache_service

import (
	"github.com/dgraph-io/ristretto"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/orcaman/concurrent-map/v2"
	"sync"
	"time"
)

// RistrettoStateBackend
//
//	@Description:
type RistrettoStateBackend struct {
	cache *ristretto.Cache
}

var localCache *RistrettoStateBackend

func init() {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		// num of keys to track frequency, usually 10*MaxCost
		NumCounters: 100000,
		// cache size(max num of items)
		MaxCost: 10000,
		// number of keys per Get buffer
		BufferItems: 64,
		// !important: always set true if not limiting memory
		IgnoreInternalCost: true,
	})
	localCache = NewRistrettoStateBackend(cache)
}

var rwMtx = sync.RWMutex{}

func NewRistrettoStateBackend(cache *ristretto.Cache) *RistrettoStateBackend {
	return &RistrettoStateBackend{cache: cache}
}

func ByteSet(key string, value []byte, ttl time.Duration) error {
	rwMtx.Lock()
	defer rwMtx.Unlock()
	localCache.cache.SetWithTTL(key, value, 0, ttl)
	localCache.cache.Wait()
	return nil
}

func ByteGet(key string, value *[]byte) error {
	rwMtx.RLock()
	defer rwMtx.RUnlock()
	if val, ok := localCache.cache.Get(key); !ok {
		return nil
	} else {
		if varByte, ok := val.([]byte); ok {
			value = &varByte
			return nil
		}
	}
	return nil
}

func StringSet(key string, value string, ttl time.Duration) error {
	rwMtx.Lock()
	logger.Info("本地缓存 StringSet key=" + key + " value=" + value)
	defer rwMtx.Unlock()
	localCache.cache.SetWithTTL(key, value, 0, ttl)
	localCache.cache.Wait()
	return nil
}

func StringSetEx(key string, value string, ttl time.Duration) error {
	rwMtx.Lock()
	logger.Info("本地缓存  StringSetEx key=" + key + " value=" + value)
	defer rwMtx.Unlock()
	if _, ok := localCache.cache.Get(key); !ok {
		return nil
	} else {
		if ttl != 0 {
			localCache.cache.SetWithTTL(key, value, 0, ttl)
		}
	}
	localCache.cache.Wait()
	return nil
}

func HashSet(key string, hashKey string, value string) error {
	rwMtx.Lock()
	defer rwMtx.Unlock()
	m := cmap.New[string]()
	m.Set(hashKey, value)
	//最长ttl设置为1month
	localCache.cache.SetWithTTL(key, m, 0, time.Duration(24*30)*time.Hour)
	localCache.cache.Wait()
	return nil
}

func Del(key string) error {
	rwMtx.Lock()
	defer rwMtx.Unlock()
	localCache.cache.Del(key)
	localCache.cache.Wait()
	return nil
}

func StringGet(key string) (string, error) {
	rwMtx.RLock()
	defer rwMtx.RUnlock()
	if val, ok := localCache.cache.Get(key); !ok {
		return "", nil
	} else {
		if valStr, transOk := val.(string); transOk {
			return valStr, nil
		}
	}
	return "", nil
}

func HasKey(key string) (bool, error) {
	rwMtx.RLock()
	defer rwMtx.RUnlock()
	if _, ok := localCache.cache.Get(key); !ok {
		return false, nil
	} else {
		return true, nil
	}
}

func Expire(key string, ttl time.Duration) error {
	rwMtx.Lock()
	defer rwMtx.Unlock()
	if val, ok := localCache.cache.Get(key); !ok {
		return nil
	} else {
		if ttl != 0 {
			localCache.cache.SetWithTTL(key, val, 0, ttl)
		}
	}
	localCache.cache.Wait()
	return nil
}
