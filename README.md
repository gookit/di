# DI

[![GoDoc](https://godoc.org/github.com/gookit/di?status.svg)](https://godoc.org/github.com/gookit/di)
[![Build Status](https://travis-ci.org/gookit/di.svg?branch=master)](https://travis-ci.org/gookit/di)
[![Coverage Status](https://coveralls.io/repos/github/gookit/di/badge.svg?branch=master)](https://coveralls.io/github/gookit/di?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/di)](https://goreportcard.com/report/github.com/gookit/di)

Golang实现依赖注入容器，提供内部服务实例管理。

## GoDoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/di.v1)
- [godoc for github](https://godoc.org/github.com/gookit/di)

## Install

```bash
go get github.com/gookit/di
```

## Usage

```go
import (
    "github.com/gookit/di"
)

func main() {
    box := di.New("my-services")
    
    // add a simple value
    box.Set("service1", "val1")
    
    // register by callback func.
    box.Set("service2", func() (interface, error) (
    	return &MyApp{}, nil
    ))
    
    // register a factory func.
    box.Set("service3", func() (interface, error) (
    	return &MyObject{}, nil
    ), true)
    
    // get 
    v1 := box.Get("service1") // "val1"
    
    // is a singleton value. Notice: v2 == v3
    v2 := box.Get("service2").(*MyApp)
    v3 := box.Get("service2").(*MyApp)
    
    // is factory func. Notice: v4 != v5
    v4 := box.Get("service3").(*MyObject)
    v5 := box.Get("service3").(*MyObject)
}
```

## API Methods

- `func (c *Container) Set(name string, val interface{}, isFactory ...bool) *Container`
- `func (c *Container) Get(name string) interface{}`
- `func (c *Container) SafeGet(name string) (val interface{}, err error)`
- `func (c *Container) Inject(ptr interface{}) (err error)`

## Refer

- https://github.com/sarulabs/di
- https://github.com/codegangsta/inject
- https://github.com/go-ozzo/ozzo-di
- https://github.com/xialeistudio/di-demo

## License

**[MIT](LICENSE)**
