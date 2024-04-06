package cache

import (
  "sync"
  "time"
)

const needCap = 100

type Cache[K comparable, V any] struct {
  mu  sync.Mutex
  s   []V
  pos map[K]int
}

func NewCache[K comparable, V any]() *Cache[K, V] {
  c := &Cache[K, V]{
    s:   make([]V, 0, needCap),
    pos: make(map[K]int),
  }
  go c.bloat()
  return c
}

func (c *Cache[K, V]) Get(key K) V {
  c.mu.Lock()
  defer c.mu.Unlock()

  idx, ok := c.pos[key]
  if !ok {
    return *new(V)
  }
  return c.s[idx]
}

func (c *Cache[K, V]) Put(key K, value V) {
  c.mu.Lock()
  defer c.mu.Unlock()

  idx, ok := c.pos[key]
  if !ok {
    c.s = append(c.s, value)
    c.pos[key] = len(c.s) - 1
    return
  }
  c.s[idx] = value
}

func (c *Cache[K, V]) Values() []V {
  c.mu.Lock()
  defer c.mu.Unlock()

  return c.s
}

func (c *Cache[K, V]) bloat() {
  t := time.NewTicker(300 * time.Millisecond)

  for {
    select {
    case <-t.C:
      c.mu.Lock()

      if cap(c.s) < needCap {
        c.mu.Unlock()
        continue
      }
      copyS := make([]V, len(c.s), needCap)
      copy(copyS, c.s)

      c.s = copyS
      c.mu.Unlock()
    }
  }
}
