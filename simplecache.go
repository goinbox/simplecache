package simplecache

import (
	"github.com/goinbox/shardmap"

	"time"
)

const (
	DefTickInterval = 30 * time.Second

	NoExpire = 0
)

type cacheItem struct {
	value  interface{}
	expire int64

	setTime    time.Time
	accessTime time.Time
}

type SimpleCache struct {
	shardMap *shardmap.ShardMap
}

func NewSimpleCache() *SimpleCache {
	return New(shardmap.DefShardCnt, DefTickInterval)
}

func New(shardCnt uint8, tickInterval time.Duration) *SimpleCache {
	s := &SimpleCache{
		shardMap: shardmap.New(shardCnt),
	}

	go s.runJanitor(tickInterval)

	return s
}

func (s *SimpleCache) Set(key string, value interface{}, expire time.Duration) {
	s.shardMap.Set(key, newCacheItem(value, expire))
}

func (s *SimpleCache) Get(key string) (interface{}, bool) {
	value, ok := s.shardMap.Get(key)
	if !ok {
		return nil, false
	}

	ci, ok := value.(*cacheItem)

	now := time.Now()
	if expired(now.UnixNano(), ci) {
		return nil, false
	}

	ci.accessTime = now

	return ci.value, true
}

func (s *SimpleCache) runJanitor(tickInterval time.Duration) {
	ticker := time.NewTicker(tickInterval)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			s.shardMap.Walk(func(k string, v interface{}) {
				ci, ok := v.(*cacheItem)
				if ok {
					if expired(now, ci) {
						s.shardMap.Del(k)
					}
				}
			})
		}
	}
}

func expired(now int64, ci *cacheItem) bool {
	if ci.expire == NoExpire {
		return false
	}

	if (now - ci.expire) >= ci.setTime.UnixNano() {
		return true
	}

	return false
}

func newCacheItem(value interface{}, expire time.Duration) *cacheItem {
	now := time.Now()

	return &cacheItem{
		value:  value,
		expire: int64(expire),

		setTime:    now,
		accessTime: now,
	}
}
