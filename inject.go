package di

import (
	"fmt"
	"reflect"
	"strings"
)

// Inject services by struct tags.
func (c *Container) Inject(ptr interface{}) (err error) {
	elemType := reflect.TypeOf(ptr).Elem()
	elemVal := reflect.ValueOf(ptr).Elem()

	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		fieldType := elemType.Field(i)
		// get tag. eg: `DI:"request"`
		tag := fieldType.Tag.Get("DI")
		name := c.getInjectName(tag)
		if name == "" {
			continue
		}

		diInstance, err := c.Get(name)
		if err != nil {
			return fmt.Errorf("dependency not found: %s", name)
		}

		elemVal.Field(i).Set(reflect.ValueOf(diInstance))
	}

	return
}

// get inject name
func (c *Container) getInjectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}

	return strings.TrimSpace(tags[0])
}
