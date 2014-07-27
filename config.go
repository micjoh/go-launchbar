package launchbar

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"time"
)

// ConfigValues represents a Config values
type ConfigValues map[string]interface{}

// Config provides permanent config utils for the action.
type Config struct {
	path string
	data map[string]interface{}
}

// NewConfig initializes an new Config object with the specified path and returns it.
func NewConfig(p string) Config {
	return loadConfig(p)
}

// NewConfigDefaults initializes a new Config object with the specified path
// and default values and returns it.
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

// Delete removes the key from config file.
func (c *Config) Delete(key string) {
	delete(c.data, key)
	c.save()
}

// Set sets the key, val and saves the config to the disk.
func (c Config) Set(key string, val interface{}) {
	c.data[key] = val
	c.save()
}

// Get gets the value from config for the key
func (c Config) Get(key string) interface{} {
	return c.data[key]
}

// Get gets the value from config for the key as string
func (c Config) GetString(key string) string {
	// TODO: fail check
	s, ok := c.data[key].(string)
	if !ok {
		return ""
	}
	return s
}

// Get gets the value from config for the key as int
func (c Config) GetInt(key string) int {
	i, ok := c.data[key].(float64)
	if !ok {
		return 0
	}
	return int(i)
}

// Get gets the value from config for the key as bool
func (c Config) GetBool(key string) bool {
	b, ok := c.data[key].(bool)
	if !ok {
		return false
	}
	return b
}

// Get gets the value from config for the key as time.Duration
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
