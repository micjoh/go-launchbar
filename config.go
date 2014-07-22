package launchbar

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"time"
)

type ConfigValues map[string]interface{}
type Config struct {
	path string
	data map[string]interface{}
}

func NewConfig(p string) Config {
	return loadConfig(p)
}

func NewConfigDefaults(p string, defaults ConfigValues) Config {
	config := loadConfig(p)
	for k, v := range defaults {
		if _, found := config.data[k]; !found {
			config.data[k] = v
		}
	}
	config.save()
	return config
}

func (c *Config) Delete(key string) {
	delete(c.data, key)
	c.save()
}

func (c Config) Set(key string, val interface{}) {
	c.data[key] = val
	c.save()
}

func (c Config) Get(key string) interface{} {
	return c.data[key]
}
func (c Config) GetString(key string) string {
	// TODO: fail check
	s, ok := c.data[key].(string)
	if !ok {
		return ""
	}
	return s
}
func (c Config) GetInt(key string) int {
	i, ok := c.data[key].(float64)
	if !ok {
		return 0
	}
	return int(i)
}
func (c Config) GetBool(key string) bool {
	b, ok := c.data[key].(bool)
	if !ok {
		return false
	}
	return b
}
func (c *Config) GetTimeDuration(key string) time.Duration {
	d, ok := c.data[key].(float64)
	if !ok {
		return 0
	}
	return time.Duration(d)
}

func loadConfig(p string) Config {
	p = path.Join(p, "config.json")
	config := Config{path: p, data: make(ConfigValues)}

	if data, err := ioutil.ReadFile(p); err == nil {
		json.Unmarshal(data, &config.data)
	}
	return config
}

func (c Config) save() {
	if data, err := json.Marshal(&c.data); err == nil {
		ioutil.WriteFile(c.path, data, 0664)
	}
}
