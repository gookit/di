package di

// Object struct definition
type Object struct {
	Name string
	instance interface{}
	callback func(c *Container) (interface{}, error)
	protected bool
}

// Instance get
func (o *Object) Instance() interface{} {
	return o.instance
}

// NewObject new object
func NewObject() *Object {
	return &Object{}
}


