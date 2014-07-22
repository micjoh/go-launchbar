package launchbar

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type CacheError string

func (c CacheError) Error() string {
	return string(c)
}

var (
	ErrCacheDoesNotExists = CacheError("the cache does not exists")
	ErrCacheIsCorrupted   = CacheError("the cache is corrupted")
	ErrCacheIsExpired     = CacheError("the cache is expired")
)

type Cache struct {
	path string
}

func NewCache(p string) Cache {
	return Cache{path: p}
}

type cacheItem struct {
	Time  *time.Time `json:"expiry"`
	Items []*item    `json:"items"`
}
type genericCache struct {
	Time *time.Time  `json:"expiry"`
	Data interface{} `json:"data"`
}

func (c *Cache) Clean() {
	c.Delete("my-recent")
	c.Delete("my-tags")
	c.Delete("my-posts")

}

func (c *Cache) Delete(key string) {
	p := path.Join(c.path, key)
	if stat, err := os.Stat(p); err == nil {
		if !stat.IsDir() {
			if filepath.Dir(p) == c.path {
				os.Remove(p)
			}
		}
	}
}

func (c *Cache) Set(key string, data interface{}, d time.Duration) {
	wd, err := os.Create(path.Join(c.path, key))
	if err != nil {
		log.Fatalln(err)
	}
	defer wd.Close()
	t := time.Now().Add(d)
	b, err := json.Marshal(genericCache{&t, data})
	if err != nil {
		log.Fatalln(err)
	}
	wd.Write(b)
}
func (c *Cache) Get(key string, v interface{}) (*time.Time, error) {
	if _, err := os.Stat(path.Join(c.path, key)); err != nil {
		return nil, ErrCacheDoesNotExists
	}
	rd, err := os.Open(path.Join(c.path, key))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, ErrCacheIsCorrupted
	}

	var e genericCache
	e.Data = &v
	// e := cacheItem{Items: &Items}
	// err = gob.NewDecoder(rd).Decode(&e)
	err = json.Unmarshal(data, &e)
	if err != nil {
		return nil, ErrCacheIsCorrupted
	}
	if time.Now().After(*e.Time) {
		return e.Time, ErrCacheIsExpired
	}
	return e.Time, nil
}

func (c Cache) SetItems(key string, items *Items, d time.Duration) {
	wd, err := os.Create(path.Join(c.path, key))
	if err != nil {
		log.Fatalln(err)
	}
	defer wd.Close()
	t := time.Now().Add(d)
	// err = gob.NewEncoder(wd).Encode([]interface{}{t.Unix(), e})
	b, err := json.Marshal(cacheItem{&t, items.getItems()})
	if err != nil {
		log.Fatalln(err)
	}
	wd.Write(b)
}

func (c *Cache) GetItemsWithInfo(key string) (*Items, *time.Time, error) {
	if _, err := os.Stat(path.Join(c.path, key)); err != nil {
		return nil, nil, ErrCacheDoesNotExists
	}
	rd, err := os.Open(path.Join(c.path, key))
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, nil, ErrCacheIsCorrupted
	}

	var e cacheItem
	// e := cacheItem{Items: &Items}
	// err = gob.NewDecoder(rd).Decode(&e)
	err = json.Unmarshal(data, &e)
	if err != nil {
		return nil, nil, ErrCacheIsCorrupted
	}
	items := &Items{}
	items.setItems(e.Items)

	if time.Now().After(*e.Time) {
		return items, e.Time, ErrCacheIsExpired
	}
	return items, e.Time, nil

}
func (c Cache) GetItems(key string) *Items {
	items, _, _ := c.GetItemsWithInfo(key)
	return items
}
