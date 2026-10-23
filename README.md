# Mo - api framework

A backend server template which is made on top of net/http and is minimalist

inspired heavily from [Echo](https://github.com/labstack/echo/)


## Features:
- Segment Trie based router, has parameter and wildcard paths. (:id,*)
- Single context which gives access to Request and ResponseWriter objects
- Has middlewares for rate limiting, 2 types as of now, token bucket algorithm and window counter
- Has in built validator which validates structs

## Get started:
```go
import "github.com/impl0x/mo"

func main() {
	m := mo.New()           // New instance of Mo
	m.GET("/", rootHandler) // Registering a handler function for path "/"
	m.Start(":8080")        // start and serve at port 8080
}

func rootHandler(c *mo.Context) error { // signature for a mo.HandlerFunc
	println(`Got request at "/"`)
	return c.JSON(200, map[string]any{"status": "success"}) // returns a json response with status code 200 and {"status":"success"}
}
```
## Documentation:  
Documentation is present [here](DOCS.md), reach out to me if you have any troubles at my email in my profile.

# 
Note about performance:  
I have tried my best to keep the balance between user experience and performance, because this framework is not meant to be the fastest of them all and be the most powerful and lowest latency in the first place. I started this project for myself to learn about golang better and learn the language itself, but then I added more and more and now I use this in all my backend related projects as its reliable to me and works. Although this is not at all thoroughly tested and bugs can exist, so if you find any then raise an issue, I will surely look into it.
#  
That's it, if there are any suggestions please reach out to me via the email in my profile.  
Thank you  
*~ Impl0x*
