package main

import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)

func main() {
	m:=mo.New()
	m.Use(mw)
	r:=m.GET("/",anh)
	r.Use(mw)
	logger.Fatal(m.Start(":8080").Error())
}

func anh(c *mo.Context) error {
		logger.Info(c.Request().URL.Path)
		return c.JSON(200,map[string]any{"message":"ok"})
	}
func log(){
	logger.Info("Check")
}

func mw(next mo.HandlerFunc)mo.HandlerFunc{
	return func(c *mo.Context) error {
		log()
		return next(c)
	}
}