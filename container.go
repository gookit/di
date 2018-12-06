package di

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// FactoryFunc service build func
type FactoryFunc func() (interface{}, error)

// Container struct definition
type Container struct {
	sync.Mutex
	name string

	// service names
	// {name: has used? 0/1}
	names map[string]uint8
	// service values, not contains FactoryFunc services.
	// {name: value}
	values map[string]interface{}
	// name aliases
	// {alias: real name}
	aliases map[string]string

	// factory func services.
	factories map[string]FactoryFunc
	// if values's value is a callback func, it created instance will storage to the instances map.
	instances map[string]interface{}
}

var goodNameReg = regexp.MustCompile(`^[a-zA-Z][\w-.]*$`)

// New a container
func New(name string) *Container {
	return &Container{
		name: name,

		names:   make(map[string]uint8),
		values:  make(map[string]interface{}),
		aliases: make(map[string]string),

		factories: make(map[string]FactoryFunc),
		instances: make(map[string]interface{}),
	}
}

/*************************************************************
 * Get service
 *************************************************************/

// Get a service by name
func (c *Container) Get(name string) interface{} {
	val, err := c.SafeGet(name)
	if err != nil {
		panic(err)
	}

	return val
}

// SafeGet safe get a service by name
func (c *Container) SafeGet(name string) (val interface{}, err error) {
	realName := c.RealName(name)

	// check exist.
	if _, ok := c.names[realName]; !ok {
		return nil, fmt.Errorf("container: the service '%s' is not exist", name)
	}

	// mark value is used.
	c.names[realName] = 1

	// in factories
	if cb, ok := c.factories[realName]; ok {
		return cb()
	}

	// in instances
	if val, ok := c.instances[realName]; ok {
		return val, nil
	}

	// in values
	val = c.values[realName]

	// if val is an callback func
	// cb, ok := val.(FactoryFunc) // ERROR: Can't check correctly.
	if cb, ok := val.(func() (interface{}, error)); ok {
		val, err = cb()
		// storage to instances
		c.instances[realName] = val
	}

	return
}

// Raw get raw value by name
func (c *Container) Raw(name string) (interface{}, error) {
	realName := c.RealName(name)

	// check exist.
	if _, ok := c.names[realName]; !ok {
		return nil, fmt.Errorf("container: the service '%s' is not exist", name)
	}

	// in factories
	if cb, ok := c.factories[realName]; ok {
		return cb, nil
	}

	// in values
	return c.values[realName], nil
}

// Value get value by name
func (c *Container) Value(name string) (val interface{}, ok bool) {
	name = c.RealName(name)
	val, ok = c.values[name]
	return
}

// Factory get factory func from factories
func (c *Container) Factory(name string) (fn FactoryFunc, ok bool) {
	name = c.RealName(name)
	fn, ok = c.factories[name]
	return
}

/*************************************************************
 * Set service
 *************************************************************/

// Set a service to container by name.
// Usage:
// 	c.Set("service1", ...)
// 	c.Set("service1", ..., true)
func (c *Container) Set(name string, val interface{}, isFactory ...bool) *Container {
	// check name
	name = goodName(name)
	hasUsed, ok := c.names[name]
	if ok && hasUsed == 1 {
		panic(fmt.Errorf("container: cannot override the '%s', it's has been used", name))
	}

	c.Lock()
	defer c.Unlock()

	// storage name
	c.names[name] = 0

	// is factory?
	if len(isFactory) > 0 && isFactory[0] {
		c.factories[name] = val.(FactoryFunc)
	} else {
		c.values[name] = val
	}

	return c
}

// Add new service to container
func (c *Container) Add(name string, val interface{}, isFactory ...bool) *Container {
	if c.Has(name) {
		return c
	}

	return c.Set(name, val, isFactory...)
}

// SetSingleton set a singleton value
func (c *Container) SetSingleton(name string, val interface{}) *Container {
	return c.Set(name, val, false)
}

// SetFactory Set Factory
func (c *Container) SetFactory(name string, factory FactoryFunc) *Container {
	return c.Set(name, factory, true)
}

/*************************************************************
 * helper methods
 *************************************************************/

// Has service name in the container
func (c *Container) Has(name string) bool {
	name = c.RealName(name)
	_, ok := c.names[name]
	return ok
}

// Del a service by name
func (c *Container) Del(name string) bool {
	name = c.RealName(name)
	if _, ok := c.names[name]; !ok { // not exist
		return false
	}

	// del name
	delete(c.names, name)

	// del aliases
	for a, r := range c.aliases {
		if r == name {
			delete(c.aliases, a)
		}
	}

	// del from factories
	if _, ok := c.factories[name]; ok {
		delete(c.values, name)
		return true
	}

	// del from values
	if _, ok := c.values[name]; ok {
		delete(c.values, name)

		if _, ok := c.instances[name]; ok {
			delete(c.instances, name)
		}
	}

	return true
}

// Clear the container
func (c *Container) Clear() {
	c.names = make(map[string]uint8)
	c.values = make(map[string]interface{})
	c.aliases = make(map[string]string)

	c.factories = make(map[string]FactoryFunc)
	c.instances = make(map[string]interface{})
}

// AddAlias set aliases for a name. alias of the method: SetAliases().
func (c *Container) AddAlias(name string, aliases ...string) {
	c.SetAlias(name, aliases...)
}

// SetAlias set aliases for a name
func (c *Container) SetAlias(name string, aliases ...string) {
	// not exist
	if _, ok := c.names[name]; !ok {
		return
	}

	for _, alias := range aliases {
		c.aliases[alias] = name
	}
}

// Aliases get aliases names
func (c *Container) Aliases() map[string]string {
	return c.aliases
}

// Name get container name
func (c *Container) Name() string {
	return c.name
}

// Names get all registered service names
func (c *Container) Names() (names []string) {
	for n := range c.names {
		names = append(names, n)
	}

	return
}

// RealName get real name
func (c *Container) RealName(name string) string {
	name = strings.TrimSpace(name)

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
