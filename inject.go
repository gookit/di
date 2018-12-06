package di

import (
	"fmt"
	"reflect"
	"strings"
)

// Inject services by struct tags.
func (c *Container) Inject(ptr interface{}) error {
	elemVal := reflect.ValueOf(ptr).Elem()
	elemType := reflect.TypeOf(ptr).Elem()

	for i := 0; i < elemType.NumField(); i++ {
		// get tag info. eg: `DI:"request"`
		name := elemType.Field(i).Tag.Get("DI")
		name = strings.TrimSpace(name)
		if name == "" { // no tag info
			continue
		}

		diInstance, err := c.SafeGet(name)
		if err != nil {
			return fmt.Errorf("dependency '%s' not found in the container", name)
		}

		elemVal.Field(i).Set(reflect.ValueOf(diInstance))
	}

	return nil
}
