package main

import (
	"fmt"
	"strings"
)

// import (
// 	// "net/http"

// 	"github.com/impl0x/mo"
// 	"github.com/impl0x/mo/modules/logger"
// )

// func main() {
// 	m:=mo.New()

// 	m.GET("/",anh)
// 	m.GET("/health", anh)
// 	logger.Fatal(m.Start(":8080").Error())
// }

// func anh(c *mo.Context) error {
// 		logger.Info(c.Request().URL.Path)
// 		return c.JSON(200,map[string]any{"message":"ok"})
// 	}

func main() {
	d:="/1"
	v:=strings.Split(d, "/")
	fmt.Printf("%T",v)
	println(v[1])
}