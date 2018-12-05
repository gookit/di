package di

// global container
var globalBox = New()

// Box get Container
func Box() *Container {
	return globalBox
}

// Get component by name
func Get(name string) (interface{}, error) {
	return globalBox.Get(name)
}
