package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_Add(t *testing.T) {
	ris := assert.New(t)

	s1 := func() (interface{}, error){
		return "ABC", nil
	}

	// reflect.TypeOf(s1)

	ris.IsType(new(factory), &s1)
}
