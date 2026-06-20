package main

import "github.com/labstack/echo/v5"

type Example struct {
	FieldOne string   `json:"field_one"`
	FieldTwo chan int `json:"field_two"` // bad field
}

func main() {
	test := echo.New()
	test.GET("/", handler)
	test.Start(":8080")
}

func handler(c *echo.Context) error {
	c.JSON(204, make(chan int)) // will return an error
	return nil
}
