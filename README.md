# Mo - api framework

A backend server template which is made on top of net/http and is lightweight

inspired heavily from Echo

### Finished making !

## Features:
- Single context which gives access to Request and ResponseWriter objects
- Has middlewares for ratelimiting, 2 types as of now, token bucket algorithm and window counter
- Has in built validator which validates structs


## Documentation:
### Demo usage
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
### Middleware usage
using the inbuilt logger middleware
```go
import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares"
)

func main() {
	m := mo.New()
	// m.Use() takes a parameter of type mo.MiddlewareFunc
	m.Use(middlewares.Logger) // <--
	m.GET("/", func(c *mo.Context) error { return nil })
	m.Start(":8080")
}
```
defining our own middleware
```go
// You can pass this function inside m.Use()
func CustomMiddleware(next mo.HandlerFunc) mo.HandlerFunc {
	return func(c *mo.Context) error {
		// do middleware stuff here
		// ...
		next(c) // call next with the same context
	}
}
```
middlewares package contains all the middlewares you can use  
you can create your own middlewares as well.  
the signature is *`func(HandlerFunc) HandlerFunc`*

### Ratelimiter usage
``` go
import (
	"time"
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares"
	"github.com/impl0x/mo/middlewares/ratelimiters"
)

func main() {
	m := mo.New()
	newWcRl := ratelimiters.NewWindowCounter(500, time.Second)
	m.Use(middlewares.Ratelimit(newWcRl))
	m.GET("/", func(c *mo.Context) error { return nil })
	m.Start(":8080")
}
```
`middlewares.Ratelimit()` takes in a type of `ratelimiters.Ratelimiter`  
which is an interface with a single function signature of func *`func Allow(*http.Request) bool`*

### Types of ratelimiters
1. **Window counter Algorithm**  
This is a Non-IP based method of ratelimiting, So it applies globally to whatever its children down the middleware chain are.  
The config for this has 2 parameters of maxRequests and duration.  
*it only allows `maxRequests` amount of requests per `duration` of time.*  
it will block every other request with a **429**.

2. **Token Bucket Algorithm**  
This is a IP based ratelimiter  
The parameters are `maxCapacity` and `refillRate`.  
The way this works is by giving every visiting IP a maximum limit of requests they can perform, i.e. the capacity of the bucket.  
And refills the bucket every 1 second by whatever the refill rate is specified to be.  
if the user exceeds the bucket capacity then they are blocked from visiting temporarily.  
Currently I haven't devised a solution to implement permanent blocking algorithms, maybe in the future I will work on that.  

### Validation
Works like the default struct validation from https://github.com/go-playground/validator  
```go
import "github.com/impl0x/mo/validator"

type User struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Age      int    `validate:"min=18,max=120"`
}

func main() {
	badUser := User{ // contains invalid invalid information
		Email:    "not_a_valid_email",
		Password: "passw",
		Age:      2,
	}
	errs := validator.Validate(&badUser)
	if errs != nil {
		for _, err := range errs {
			println(err.Error())
		}
	}
}

```
has a few rules here
- required: will fail if the field is zero value
- email: matches against a regex
- url: same as above
- min: if string, length need to satisfy this, if number then need to be more than this
- max: same logic as above
- oneof: need to be one of the valid options
### Header management  
#### *Response headers* 
#### Default headers:  
These get set on every request globally.  
Example:  
```go
m := mo.New()
m.Headers["x-test"] = "test"
```
#### Request specific headers:  
These get set on requests if you set them using context
```go
func ExampleHandler(c *mo.Context) error {
	c.Headers["x-test"] = "test"
}
```
Both are of the type `map[string]string`. 

*Request headers*
