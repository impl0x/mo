# Mo - api framework

A backend server template which is made on top of net/http and is minimalist

inspired heavily from [Echo](https://github.com/labstack/echo/)


## Features:
- Compact Trie based router, has parameter and wildcard paths. (:id,*)
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
	m.Use(middlewares.Logger()) // <--
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

### Grouping routes
You can group routes with the prefix of the route.  
example usage:  
```go
import "github.com/impl0x/mo"

func main() {
	m := mo.New()
	v1Group := m.Group("/api/v1") // creates a group with the prefix of /api/v1
	authGroup := v1Group.Group("/auth") // creates a sub group with /auth prefix
	authGroup.POST("/login", loginHandler) // registers handler with the final path of POST /api/v1/auth/login 
	m.Start(":8080") //starts the server
}
```
*Important*: Middlewares **must** be added before making a sub group or registering a path  
**Wrong way** ❌ :
```go
v1Group := m.Group("/api/v1")
authGroup := v1Group.Group("/auth")

v1Group.Use(v1Middleware) // will NOT be registered for authGroup ❌
authGroup.GET("/login",loginHandler) // the middleware will not run in this case 
```
**Right way** ✅:
```go
v1Group := m.Group("/api/v1")
v1Group.Use(v1Middleware)

authGroup := v1Group.Group("/auth") // we register the subgroup after registering the middleware ✅
authGroup.GET("/login",loginHandler) // the middleware will work as expected in this case
```
Same with grouped paths aswell
```go
customGroup:=m.Group("/custom")
customGroup.GET("/1",customOneHandler)

customGroup.Use(middlewares.Logger()) // This will NOT be registered for the above path 

customGroup.GET("/2",customTwoHandler) // But this will have the middleware registered for this path
```
I hope it is clear and I feel this is pretty intuitive.  
*Rule of thumb*: register middleware and groups first. then register paths.
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
Works like the default struct validation from [Go's Validator package](https://github.com/go-playground/validator)  
This is a experiment and some features might not work as expected and it isn't as extensive as the original validator package. I wrote my own validator just for the sake of learning and gaining experience of how it works under the hood, completely up to you if you want to use this package or the original one, I would recommend stick with the battle tested one. As I am sure it is way more optimized and faster than my implementation.
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
- optional: will skip if its a zero value
- dive: used in arrays/slices to validate every element against the provided rules
- email: matches against a regex
- url: same as above
- ipv4, ipv6
- alpha, alphanum
- e164: for phone numbers
- uuid: validates the uuid syntax
- min: if string, length need to satisfy this, if number then need to be more than this
- max: same logic as above
- oneof: need to be one of the valid options
- lte: less than or equal to
- lt: lesser than
- gte: greater than or equal to
- gt: greater than
- len: must be exactly this long, applies to string, arrays, slices, maps.

#### Error format for validator
The `validator.Validate` function returns a `validator.GroupedValidationError` type.  
Whose underlying type is of `[]ValidationError`. Where `ValidationError` is a interface embedding `error` and having a single function of `NameSpace() string`.    
Two structs satisfy this interface, `FieldValidateError` and `UserError`, the names are pretty self explanatory, `UserError` is only present if the user used the validation tags wrong.  
You can iterate over it and type assert for `FieldValidateError` struct to get access to the methods which return  
almost everything you need to make a custom error message with it.  
Else you can directly call the `.ToJsonStructList()` method on the `GroupedValidationError` that was initially returned, this returns a slice of structs which are json compatible (so you can directly pass it to json encoder),  
the default format of errors which is in this format:  
example: if we validated a wrong email and a wrong url, and then called the `.ToJsonStructList()` on that error we get this as a `[]ValidationErrorJson`.  
```json
[
	{
		"field":"email",
		"message":"Not a valid email"
	},
	{
		"field":"url",
		"message":"Not a valid url"
	}
]
```
It can also handle nested structs validation, and there are configs which you can manipulate to tweak the error json according to your likings. For example a nested error would look like this, if user field was a struct containing a email field.
```json
"field":"user.email"
```
Diving into slices also returns a similar response, `users.3`, indicating the third index   
The configs are present in   
- `validator.ErrorConfig` : Has 2 fields, `ReturnUserErrors` and `LogUserErrors`, both bool, pretty self explanatory.
- `validator.DefaultNameSpaceSettings`: Has 3 fields, take a look at the `NameSpaceSettings` struct in [here](validator/v2.go#L47) to understand what each field does. This config is used to modify the error message for nested structs
### Header management  
#### *Response headers* 
#### Default headers:  
These get set on every request globally.  
Example:  
```go
m := mo.New()
m.Headers.Add("x-test", "test")
```
#### Request specific headers:  
These get set on requests if you set them using context
```go
func ExampleHandler(c *mo.Context) error {
	c.ResponseHeaders.Add("x-test","test")
}
```
Both are of the type `map[string]string`. 

To bind headers from a struct use the .bind() method on headers object.  
for example,
```go
type DefaultHeaders struct {
	XtraceId   string `header:"x-trace-id"`
	XpoweredBy string `header:"x-powered-by"`
}

m := mo.New()
headers := DefaultHeaders{
	"abcd",
	"mo",
}
m.Headers.Bind(&headers)
```  
The struct that you are binding must contain the `header` tag.  

### Request body parsing  
Just call the `c.Bind` method and pass a struct with json tags
```go
func handler(c *mo.Context) error {
	c.DecodeBody(target) // just decodes
	c.DecodeBodyAndValidate(target) // decodes and validates the body
}
```
you can also validate it at the same place using the other function.

# 
Some features that I am aware of but unwilling to add:
- Pre compiling the middlewares and creating one compiled handler and storing that in the router instead of chaining everything on ServeHTTP function when a request arrives. Yes this saves heap allocations and makes it a bit faster, but this also means changing how middlewares are registered and how paths are stored fundamentally, because then I will have to figure out how to layer middlewares depending upon how they are registered and possibly have to enforce more rules upon the user to register middlewares a certain way and paths a certain way to make them work properly. Which feels tedious and will make me change a lot of code in the codebase as of now, and even though it has performance benefits I feel having a better user experience is more important.  

I have tried my best to keep the balance between user experience and performance, because after-all we all have relatively powerful systems which can handle a good amount of requests a second and this framework is not meant to be the fastest of them all and be the most powerful and lowest latency in the first place. I started this project for myself to have a reuseable backend template so being able to have comfort is better than having a tiny bit of extra performance.
#  
Thats it, If you find any bugs please raise an issue. And if there are any suggestions please reach out to me via gmail in my profile.  
Thank you  
*~ Impl0x*