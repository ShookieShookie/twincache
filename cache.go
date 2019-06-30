package twincache

import (
	"fmt"
	"github.com/ShookieShookie/ringslice"
	"sync"
	"time"
)

type Cache struct {
	capacity   int
	expiries   *ringslice.Slice
	values     map[interface{}]interface{}
	expiryTime time.Duration
	sync.RWMutex
}

type Unit struct {
	key    interface{}
	expiry int64
}

func New(cap int, dur time.Duration) *Cache {
	c := &Cache{capacity: cap, expiries: ringslice.NewSlice(cap, false, wipe), values: make(map[interface{}]interface{}, cap), expiryTime: dur}
	t := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-t.C
			c.expiries.Stats()
		}
	}()
	return c
}

// TODO have to enforce some type safety i imagine?
func (c *Cache) Add(key interface{}, val interface{}) error {
	expiry := time.Now().Add(c.expiryTime).Unix()
	u := &Unit{key: key, expiry: expiry}
	c.Lock()
	c.values[key] = val
	c.expiries.Append(u)
	c.Unlock()
	return nil
}
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	found, ok := c.values[key]
	c.RUnlock()
	if !ok {
		return nil, false
	}
	return found, true
}

func (c *Cache) ExpiredIndex() int {
	ind := c.expiries.FindClosestBelowOrEqual(time.Now().Unix(), func(i interface{}) int64 { return i.(*Unit).expiry })
	return ind
}

func wipe(ind int, l []interface{}) {
	l[ind].(*Unit).expiry = 0
}

func expiryFromUnit(i interface{}) int64 {
	return i.(*Unit).expiry
}

func (c *Cache) Expire() int {

	c.Lock()
	ind := c.expiries.FindClosestBelowOrEqual(time.Now().Unix(), expiryFromUnit)
	fmt.Println(c.expiries.Values(expiryFromUnit))
	c.expiries.Purge(time.Now().Unix(), expiryFromUnit)
	c.Unlock()
	return ind
	// get all the expired, give back the keys, wipe
}
