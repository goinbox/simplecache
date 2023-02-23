package simplecache

import (
	"github.com/goinbox/crypto"

	"github.com/stretchr/testify/assert"

	"strconv"
	"testing"
	"time"
)

var sc *SimpleCache

func init() {
	sc = NewSimpleCache()
}

func TestSetGet(t *testing.T) {
	for i := 0; i < 10000; i++ {
		key := crypto.Md5String([]byte(strconv.Itoa(i)))
		sc.Set(key, i, 10*time.Second)

		v, ok := sc.Get(key)
		if !ok || v != i {
			t.Error(v, ok)
		}
	}

	time.Sleep(16 * time.Second)

	for i := 0; i < 10000; i++ {
		key := crypto.Md5String([]byte(strconv.Itoa(i)))

		v, ok := sc.Get(key)
		if ok || v == i {
			t.Error(v, ok)
		}
	}
}

func TestSetNX(t *testing.T) {
	key := "abc"
	ok := sc.SetNX(key, 1, time.Second*1)
	assert.True(t, ok)

	time.Sleep(time.Second * 2)
	_, ok = sc.Get(key)
	assert.False(t, ok)

	ok = sc.SetNX(key, 1, time.Second*1)
	assert.True(t, ok)
	ok = sc.SetNX(key, 1, time.Second*1)
	assert.False(t, ok)
}

func TestDel(t *testing.T) {
	key := "abc"
	sc.Set(key, 1, NoExpire)
	_, ok := sc.Get(key)
	assert.True(t, ok)

	sc.Del(key)
	_, ok = sc.Get(key)
	assert.False(t, ok)
}
