package di

// Box always create a global container
var Box = New()

// Get service component from the global container
func Get(name string) (interface{}, error) {
	return Box.Get(name)
}

// Set a service component to the global container
func Set(name string, val interface{}, singleton ...bool) *Container {
	return Box.Set(name, val, singleton...)
}

// Has name in the global container
func Has(name string) bool {
	return Box.Has(name)
}
