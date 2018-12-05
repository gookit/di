package di

import (
	"errors"
	"fmt"
	"regexp"
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
	instances  map[string]interface{}

	// service names {name: initialized 0/1}
	names   map[string]uint8
	values  map[string]interface{}
	rawData map[string]interface{}
	// name aliases [alias => real name]
	aliases map[string]string
}

var goodNameReg = regexp.MustCompile(`\w[\w-.]+`)

var ErrFactoryNotFound = errors.New("container: factory is not found")

// New a container
func New() *Container {
	return &Container{
		names:   make(map[string]uint8),
		aliases: make(map[string]string),
	}
}

// Get
func (c *Container) Get(name string) (val interface{}, err error) {
	realName := c.getRealName(name)
	_, ok := c.names[realName]

	if !ok {
		return nil, fmt.Errorf("container: the '%s' service is not exist", name)
	}

	// in singletons
	if val, ok := c.singletons[realName]; ok {
		return val, nil
	}

	// in factories
	if cb, ok := c.factories[realName]; ok {
		return cb()
	}

	return nil, fmt.Errorf("container: the '%s' service is not exist", name)
}

// Add new service to container
func (c *Container) Add(name string, val interface{}, singleton bool) {
	if c.Has(name) {
		return
	}

	c.Set(name, val, singleton)
}

// Set a service to container by name
func (c *Container) Set(name string, val interface{}, singleton bool) {
	// check name
	name = goodName(name)

	c.Lock()
	defer c.Unlock()

	hasUsed, ok := c.names[name]
	if ok && hasUsed == 1 {
		panic(fmt.Errorf("container: cannot override the '%s', it's has been used", name))
	}

	// storage
	c.names[name] = 0

	if singleton {
		c.singletons[name] = val
	} else {
		c.factories[name] = val.(factory)
	}
}

// SetSingleton Set Singleton
func (c *Container) SetSingleton(name string, val interface{}) {
	c.Set(name, val, true)
}

// SetFactory Set Factory
func (c *Container) SetFactory(name string, factory factory) {
	c.Set(name, factory, false)
}

func (c *Container) GetSingleton(name string) interface{} {
	return c.singletons[name]
}

// GetInstance Get Instance from a factory
func (c *Container) GetInstance(name string) (interface{}, error) {
	factory, ok := c.factories[name]
	if !ok {
		return nil, ErrFactoryNotFound
	}

	return factory()
}

// Has service name in the container
func (c *Container) Has(name string) bool {
	name = c.getRealName(name)
	_, ok := c.names[name]
	return ok
}

// Del a service by name
func (c *Container) Del(name string) bool {
	name = c.getRealName(name)
	if _, ok := c.names[name]; !ok {
		return false
	}

	delete(c.names, name)

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
	for n := range c.names {
		names = append(names, n)
	}

	return
}

// Aliases get aliases names
func (c *Container) Aliases() map[string]string {
	return c.aliases
}

// SetAliases for a name
func (c *Container) SetAliases(name string, aliases ...string) {
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

func goodName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		panic("container: the added name cannot be empty")
	}

	if !goodNameReg.MatchString(name) {
		panic(`container: the added name is invalid, must match regex '\w[\w-.]+'`)
	}

	return name
}
