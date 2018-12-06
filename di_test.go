package di

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestContainer(t *testing.T) {
	ris := assert.New(t)

	// var s1 interface{}
	s1 := func() (interface{}, error) {
		return "ABC", nil
	}

	// Set
	Set("s1", s1).AddAlias("s1", "the-s1")
	Set("s2", "val2")
	Box.Add("s2", "new-val")
	Box.Add("s3", "val3")

	// invalid name
	ris.Panics(func() {
		Set("", "val")
	})
	ris.Panics(func() {
		Set("+sf", "val")
	})

	Box.AddAlias("not-exist", "n1")

	// exist
	ris.True(Has("s1"))
	ris.True(Has("s2"))
	ris.True(Has("the-s1"))
	ris.False(Box.Has("not-exist"))
	ris.Contains(Box.Aliases(), "the-s1")

	// Get value
	val, err := Get("s1")
	ris.Nil(err)
	ris.Equal("ABC", val)

	// get by alias name
	val, err = Get("the-s1")
	ris.Nil(err)
	ris.Equal("ABC", val)

	val, err = Get("s2")
	ris.Nil(err)
	ris.Equal("val2", val)

	val, err = Get("not-exist")
	ris.Error(err)
	ris.Nil(val)

	// Raw value
	val, err = Box.Raw("s1")
	ris.Nil(err)
	ris.Equal(reflect.TypeOf(s1), reflect.TypeOf(val))

	val, err = Box.Raw("not-exist")
	ris.Error(err)
	ris.Nil(val)

	// Value
	val, ok := Box.Value("s2")
	ris.True(ok)
	ris.Equal("val2", val)

	// Names
	ris.Contains(Box.Names(), "s1")
	ris.Contains(Box.Names(), "s2")

	// Del
	ris.True(Box.Del("s1"))
	ris.False(Box.Has("s1"))
	ris.True(Box.Del("s2"))
	ris.False(Box.Has("s2"))
	ris.False(Box.Del("not-exist"))

	// Clear
	Box.Clear()
	ris.Len(Box.Names(), 0)
}

func TestContainer_SetSingleton(t *testing.T) {
	ris := assert.New(t)

	s1 := func() (interface{}, error) {
		val := fmt.Sprint(time.Now().Nanosecond())
		return val, nil
	}

	Box.SetSingleton("st1", s1)

	ris.True(Box.Has("st1"))

	// Value
	val, ok := Box.Value("st1")
	ris.True(ok)
	ris.Equal(reflect.TypeOf(s1), reflect.TypeOf(val))

	// get
	val1, err := Get("st1")
	ris.Nil(err)

	// get again
	val2, err := Get("st1")
	ris.Nil(err)

	// value always is equals.
	ris.Equal(val1, val2)

	// cannot override the 'st1', it's has been used
	ris.Panics(func() {
		Box.SetSingleton("st1", "val")
	})

}

func TestContainer_SetFactory(t *testing.T) {
	ris := assert.New(t)

	f1 := func() (interface{}, error) {
		val := fmt.Sprint(time.Now().Nanosecond())
		return val, nil
	}

	// add factory func
	Box.SetFactory("f1", f1)

	ris.True(Box.Has("f1"))

	// Value
	_, ok := Box.Value("s2")
	ris.False(ok)

	// get
	val1, err := Get("f1")
	ris.Nil(err)

	// get again
	val2, err := Get("f1")
	ris.Nil(err)

	// value always is not equals.
	ris.NotEqual(val1, val2)

	raw, err := Box.Raw("f1")
	ris.Nil(err)
	ris.Equal(reflect.TypeOf(FactoryFunc(f1)), reflect.TypeOf(raw))

	fn, ok := Box.Factory("f1")
	ris.True(ok)
	ris.Equal(reflect.TypeOf(FactoryFunc(f1)), reflect.TypeOf(fn))

	// del
	ris.True(Box.Del("f1"))
	ris.False(Box.Has("f1"))
}

func TestContainer_Inject(t *testing.T) {
	//
}