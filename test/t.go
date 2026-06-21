package main

import (
	"net/http"

	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)



func main() {
	m:=mo.New()

	m.GET("/",func(c *mo.Context) error {
		return mo.NewHTTPError(http.StatusBadRequest,"Bad request")
	})
	logger.Fatal(m.Start(":8080").Error())
}

