// Package di is an dependency injection container implements.
//
// Source code and other details for the project are available at GitHub:
//
// 	https://github.com/gookit/config
//
package di

// Box always create a global container
var Box = New("global")

// Get service component from the global container
func Get(name string) interface{} {
	return Box.Get(name)
}

// SafeGet service component from the global container
func SafeGet(name string) (interface{}, error) {
	return Box.SafeGet(name)
}

// Set a service component to the global container
func Set(name string, val interface{}, isFactory ...bool) *Container {
	return Box.Set(name, val, isFactory...)
}

// Has name in the global container
func Has(name string) bool {
	return Box.Has(name)
}
