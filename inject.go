package di

import (
	"fmt"
	"reflect"
	"strings"
)

// Inject component
func (c *Container) Inject(ptr interface{}) (err error) {
	elemType := reflect.TypeOf(ptr).Elem()
	ele := reflect.ValueOf(ptr).Elem()

	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		fieldType := elemType.Field(i)
		// 获取tag `DI:"request"`
		tag := fieldType.Tag.Get("DI")

		name := c.getInjectName(tag)
		if name == "" {
			continue
		}

		var diInstance interface{}

		name = c.getRealName(name)
		if val, ok := c.singletons[name]; ok {
			diInstance = val
		}

		// in factories
		if cb, ok := c.factories[name]; ok {
			diInstance, err = cb()
			if err != nil {
				return err
			}
		}

		if diInstance == nil {
			return fmt.Errorf("dependency not found: %s", name)
		}

		ele.Field(i).Set(reflect.ValueOf(diInstance))
	}

	return
}

// 获取需要注入的依赖名称
func (c *Container) getInjectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}

	return strings.TrimSpace(tags[0])
}

