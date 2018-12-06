# DI

[![GoDoc](https://godoc.org/github.com/gookit/di?status.svg)](https://godoc.org/github.com/gookit/di)
[![Build Status](https://travis-ci.org/gookit/di.svg?branch=master)](https://travis-ci.org/gookit/di)
[![Coverage Status](https://coveralls.io/repos/github/gookit/di/badge.svg?branch=master)](https://coveralls.io/github/gookit/di?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/di)](https://goreportcard.com/report/github.com/gookit/di)

> **[EN README](README.md)**

Golang实现依赖注入容器，提供内部服务实例管理。

## GoDoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/di.v1)
- [godoc for github](https://godoc.org/github.com/gookit/di)

## 安装

```bash
go get github.com/gookit/di
```

## 快速使用

```go
package main

import (
    "github.com/gookit/di"
)

func main() {
    box := di.New("my-services")
    
    // 添加简单的值
    box.Set("service0", "val1")
    box.Set("service1", &MyService1{})
    
    // 注册一个回调函数。
    // - 它会在第一次获取时执行。
    // - 执行返回的结果，会被保存起来
    // - 后续获取时都是拿到的第一次的结果，不会再次执行这个函数
    box.Set("service2", func() (interface, error) {
    	return &MyApp{}, nil
    })
    
    // 注册一个工厂函数
    // - 每次获取时，都会执行这个函数。不会保存结果，每次都返回新的。
    box.Set("service3", func() (interface, error) {
    	return &MyObject{}, nil
    }, true)
    
    // 获取 
    v1 := box.Get("service1") // "val1"
    
    // 是一个闭包，会自动执行它. 注意: v2 == v3
    v2 := box.Get("service2").(*MyApp)
    v3 := box.Get("service2").(*MyApp)
    
    // 是一个工厂方法. 注意: v4 != v5
    v4 := box.Get("service3").(*MyObject)
    v5 := box.Get("service3").(*MyObject)
}
```

## API 方法

- `func (c *Container) Set(name string, val interface{}, isFactory ...bool) *Container`
- `func (c *Container) Get(name string) interface{}`
- `func (c *Container) SafeGet(name string) (val interface{}, err error)`
- `func (c *Container) Inject(ptr interface{}) (err error)`

## 参考

- https://github.com/sarulabs/di
- https://github.com/codegangsta/inject
- https://github.com/go-ozzo/ozzo-di
- https://github.com/xialeistudio/di-demo

## License

**[MIT](LICENSE)**
