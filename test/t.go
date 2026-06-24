package main

import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares"
	"github.com/impl0x/mo/middlewares/ratelimiters"
)

func main() {
	m := mo.New()
	newRl:=ratelimiters.NewTokenBucket(2, 1)
	m.GET("/",func(c *mo.Context) error {
		return c.JSON(200, map[string]any{"status":"ok"})
	})
	m.Start(":8080")
}
