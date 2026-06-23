package main

import (
	// "strings"

	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)

func main() {
	m:=mo.New()
	users:=mo.NewGroup("/users")
	users.Use(logUsers)
	users.GET("/1",func(c *mo.Context) error {
		logger.Mo("Test")
		return nil
	})
	m.Start(":8080")
}

func logUsers(next mo.HandlerFunc)mo.HandlerFunc{
	return func(c *mo.Context) error {
		logger.Info("User log check")
		return next(c)
	}
}