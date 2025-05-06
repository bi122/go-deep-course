package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}
type SingletonService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	store map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		store: make(map[string]interface{}),
	}
}

func (c *Container) exist(name string) bool {
	_, res := c.store[name]
	return res
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	if c.exist(name) {
		return
	}

	c.store[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, fn interface{}) {
	if c.exist(name) {
		return
	}

	constructor, ok := fn.(func() interface{})
	if !ok {
		return
	}

	c.store[name] = constructor()
}

func (c *Container) Resolve(name string) (interface{}, error) {
	fn, ok := c.store[name]
	if !ok {
		return nil, fmt.Errorf("%s not exist", name)
	}

	switch fn := fn.(type) {
	case func() interface{}:
		return fn(), nil
	default:
		return fn, nil
	}
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})
	container.RegisterSingletonType("SingletonService", func() interface{} {
		return &SingletonService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	singletonService1, err := container.Resolve("SingletonService")
	assert.NoError(t, err)
	singletonService2, err := container.Resolve("SingletonService")
	assert.NoError(t, err)

	s1 := singletonService1.(*SingletonService)
	s2 := singletonService2.(*SingletonService)
	assert.True(t, s1 == s2)
}
