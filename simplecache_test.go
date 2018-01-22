package simplecache

import (
	"github.com/goinbox/crypto"
	"github.com/goinbox/gomisc"

	"strconv"
	"testing"
	"time"
)

var sc *SimpleCache

func init() {
	sc = NewSimpleCache()
}

func TestSetGet(t *testing.T) {
	gomisc.PrintCallerFuncNameForTest()

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
