package di

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type factory func() (interface{}, error)
type diFactory func(c *Container) (interface{}, error)

// Container struct definition
type Container struct {
	sync.Mutex

	factories  map[string]factory
	singletons map[string]interface{}

	keys    map[string]uint8
	values  map[string]interface{}
	rawData map[string]interface{}
	// name aliases [alias => real name]
	aliases map[string]string
}

var ErrFactoryNotFound = errors.New("container: factory is not found")

// New a container
func New() *Container {
	return &Container{
		keys:    make(map[string]uint8),
		aliases: make(map[string]string),
	}
}

// Add new service to container
func (c *Container) Add(name string, val interface{}, aliases ...string) {

}

// Set
func (c *Container) Set(name string, val interface{}, aliases ...string) {

}

// Get
func (c *Container) Get(name string) (val interface{}, err error) {
	name = c.getRealName(name)

	// in singletons
	if val, ok := c.singletons[name]; ok {
		return val, nil
	}

	// in factories
	if cb, ok := c.factories[name]; ok {
		return cb()
	}

	return
}

// Singleton
func (c *Container) SetSingleton(name string, val interface{}) *Container {
	c.Lock()

	c.keys[name] = 1
	c.singletons[name] = val

	c.Unlock()

	return c
}

// Factory
func (c *Container) SetFactory(name string, factory factory) *Container {
	c.Lock()

	c.keys[name] = 1
	c.factories[name] = factory

	c.Unlock()
	return c
}

func (c *Container) GetSingleton(name string) interface{} {
	return c.singletons[name]
}

// 获取实例对象
func (c *Container) GetPrototype(name string) (interface{}, error) {
	factory, ok := c.factories[name]
	if !ok {
		return nil, ErrFactoryNotFound
	}

	return factory()
}

// Has service name in the container
func (c *Container) Has(name string) bool {
	name = c.getRealName(name)
	_, ok := c.keys[name]

	return ok
}

// Del a service by name
func (c *Container) Del(name string) bool {
	name = c.getRealName(name)

	if _, ok := c.keys[name]; !ok {
		return false
	}

	delete(c.keys, name)

	// delete aliases
	for a, r := range c.aliases {
		if r == name {
			delete(c.aliases, a)
		}
	}

	if _, ok := c.values[name]; ok {
		delete(c.values, name)
	}

	return true
}

// Names get all registered service names
func (c *Container) Names() (names []string) {
	for n := range c.keys {
		names = append(names, n)
	}

	return
}

// AddAliases for a name
func (c *Container) AddAliases(name string, aliases ...string) {
	for _, alias := range aliases {
		c.aliases[alias] = name
	}
}

// get real name
func (c *Container) getRealName(name string) string {
	if realName, ok := c.aliases[name]; ok {
		return realName
	}

	return name
}
