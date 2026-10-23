
## Documentation:

### Context basics:  
The `*mo.Context` is provided on every handler as the only argument, and this is all you need, because it contains all necessary functions to access request, responses, write responses, write headers, basically everything.  

#### Store:  
A context also contains a Store, a Store basically is a map of values that can be passed down in the handler chain using middlewares, the usage of Store is demonstrated below.  
```go
// Middleware that adds a request id to the context Store
func AddRequestID(next mo.HandlerFunc) mo.HandlerFunc {
	return func(c *mo.Context) error {
		c.Add("req_id", "requestID123")
		return next(c)
	}
}

// Handler function that is wrapped with the middleware above
func Handler(c *mo.Context) error {
	reqIDAny, ok := c.Get("req_id")
	if !ok {
		panic("request id should be present in context store")
	}
	reqID := reqIDAny.(string)
	return c.TEXT(200, "request id: "+reqID)
}
```
There are a few methods of context Store, the most used ones are to Add and Get a value, Adding is pretty simple as shown above, and the above Get function is the most basic example. It returns a `any` type which you need to type cast.
There is also a other way to do this and skip the boilerplate.  
for example if you know that the value exists in the store and do not want to check the `ok` field required on a Get call everytime then you can use `c.GetMust`, this function wraps the Get call with the panic statement such that if the value does not exist it panics.  
Another better way is to use generics, there is another method called `c.GetTyped[T]` which returns a typed value of type `T`, and an `ok` field optionally. Another method `c.GetTypedMust[T]` exists to wrap this function with a panic statement such that the value must exist and match the type or it panics. 
In a real app store is used only by the app itself to communicate and send values down the chain, so if your code is correct then you can safely use the `c.GetTypedMust[T]` because it removes all the boilerplate.

An example is shown below of how getting values from the context store looks like
```go
func GetExamples(c *mo.Context) error {
	anyValue, ok := c.Get("key")
	anyValue = c.GetMust("key")
	typedValue, ok := c.GetTyped[Type]("key")
	typedValue = c.GetTypedMust[Type]("key")
}
```

#### Writing responses:  
Multiple methods are present in `*mo.Context` to write responses, given below are some examples.  
```go
c.BLOB(200, "text/plain", []byte("BLOB response example"))
c.TEXT(200, "Text response example")
c.JSON(200, map[string]string{"key": "json response example"})
```

#### Getting url parameters:
The router supports "*" wildcard paths and ":" parameter paths.
for example a parameter path while registering could look like - "/api/v1/users/:id" and the router will match this to anything that looks like this "api/v1/users/123" and store the value of "id" in this case to the context Store's param map.

In summary, if you register a path which has "/users/:id", and want the id in the handler, use the method `c.Param("id")`, passing "id" in this scenario, the `c.Param` method returns a string and a bool, to signify if it exists. `c.ParamMust` also exists wrapping the previous function, to panic if the value doesn't exist.

 Information about wildcard and others given in detail on the method documentation of `mo.Context.Param` so check that out.
### Middleware usage
using the inbuilt logger middleware
```go
import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares"
)

func main() {
	m := mo.New()
	// m.Use() takes a parameter of type mo.Middleware
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
the signature is *`func(HandlerFunc) HandlerFunc`* same as the type *`mo.Middleware`*.  

**Note**: Middlewares are precomputed to the handlers, which means at runtime middlewares aren't composed, but when a middleware is added, it is stored and after every handler add that handler is wrapped with the middlewares that are stored.  
Which means, adding a middleware after registering a handler will lead to the handler not using that middleware, so middlewares must always be added before adding a handler. This also means that as the middlewares are ran at startup to wrap the respective handlers the code inside the middleware function is ran. So it is safe to put code that needs to run only once inside the middleware function and can be used as closure inside the handler function.  

In summary,  
- Add middlewares before registering paths
- Use middleware functions to run code that only needs to run once and be used as closures in the handler function.  

Examples are given in the next topic Grouping routes as well as below for further clarification.  
```go
func AuthorizationMiddleware(next mo.HandlerFunc) mo.HandlerFunc {
	errUnauthorized := errors.New("Unauthorized") // this will run only once
	return func(c *mo.Context) error { // this function will run on every call
		auth := c.Request().Header.Get("authorization")
		if auth == "" {
			// not practical to return a errors.New here 
			// just for example purposes, use 
			// proper error handling in a real app
			return errUnauthorized
		}
		return next(c)
	}
} 
```
the above example shows how closures are used in middlewares, the returning handler function holds a closure variable to `errUnauthorized`.
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
Same with grouped paths as-well
```go
customGroup:=m.Group("/custom")
customGroup.GET("/1",customOneHandler)

customGroup.Use(middlewares.Logger()) // This will NOT be registered for the above path 

customGroup.GET("/2",customTwoHandler) // But this will have the middleware registered for this path
```
I hope it is clear and I feel this is pretty intuitive.  
*Rule of thumb*: register middleware and groups first. then register paths.
### Rate limiter usage
``` go
import (
	"time"
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares"
	"github.com/impl0x/mo/middlewares/ratelimiters"
)

func main() {
	m := mo.New()
	wc := ratelimiters.NewWindowCounter(500, time.Second)
	m.Use(middlewares.Ratelimit(wc))
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
This is a IP based rate limiter  
The parameters are `maxCapacity` and `refillRate`.  
The way this works is by giving every visiting IP a maximum limit of requests they can perform, i.e. the capacity of the bucket.  
And refills the bucket every 1 second by whatever the refill rate is specified to be.  
if the user exceeds the bucket capacity then they are blocked from visiting temporarily.  
Currently I haven't devised a solution to implement permanent blocking algorithms, maybe in the future I will work on that.  

### Validation
Works like the default struct validation from [Go's Validator package](https://github.com/go-playground/validator)  
This is a experiment and some features might not work as expected and it isn't as extensive as the original validator package. I wrote my own validator just for the sake of learning and gaining experience of how it works under the hood, completely up to you if you want to use this package or the original one, I would recommend stick with the original battle tested and most popular go-playground/validator instead of mine. As I am sure it is way more optimized and faster than my implementation.

There are currently 2 versions working which I have developed, v2 and v3. Both work but v3 is the newest and has caching and dynamic function generation to save lots of time in reflection, there is documentation for both, they are pretty similar but with some small differences.

#### Docs for v2
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
- startswith: must start with a prefix string
- endswith: must end with a suffix string

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
- `validator.DefaultNameSpaceSettings`: Has 3 fields, take a look at the `NameSpaceSettings` struct in [here](validator/v2/v2.go#L47) to understand what each field does. This config is used to modify the error message for nested structs

#### Docs for v3
Although mostly same as v2 on the user side, so most things will be the same such as config. Error handling is also the same but a bit different, which is explained further down below.  
 **rules**:  
 rules are the string keys that are used to apply validations, such as "required", "optional", "len=3", etc. There are two types of rules non equal and equal.
 
 - non equal: the rules that don't require a second argument, example: required, optional, email, uuid, dive, etc.
 - equal: the rules that require two or more arguments separated by equal to sign, example: min=2, max=5, oneof=one two three, etc.

 non equal rule names are present [here](validator/v3/rules/noneq.go#L27)  
equal rule names are present [here](validator/v3/rules/eq.go#L28)

every rule has types that it applies to, for example the rule "required" can apply to any type in golang, but the rule "email" only applies to a string. This allowed types for a rule are determined and maintained in these variable declarations [here for non equal](validator/v3/rules/noneq.go#L43) and [here for equal](validator/v3/rules/eq.go#L41). Equal rules go a bit further by also having types that it accepts for the parameter, as equal rules need a parameter to be passed along side the rule itself that parameter also has a few allowed types, for example the rule "min" cannot have a parameter type of string, it is only allowed a type of number (floats,uints and ints).

**inner workings**:  
before we explain the error handling here, you need to understand a bit about how the validation engine works in the first place.

When a struct is passed to be validated for the first time it is read and all the rules are read, the rule types are validated as stated above and if anything is wrong with how the user uses the validation tags then a `*UserError` is returned, but if everything is fine then functions are composed according to the rules for each field and stored, then validation runs using those functions stored for each field. When that struct is passed to be validated the second time those functions stored are ran in order to each field respectively again just like before and this happens instead of reading rules every time, the functions are composed once and used throughout, and if validation fails, that is the value does not satisfy the rule function then a `*FieldValidateError` is returned.

**errors**:  
The return type for `Validate` is `error`, and that error is either a `*UserError` returned if struct tags are used incorrectly, (remember this is returned continuously if not fixed so make sure struct tags are used correctly). Else if the first run for validation succeeds without a `*UserError` then it is always guaranteed that the `Validate` function will always return a `GroupedValidationError` which is a underlying type of `[]*FieldValidateError`, so just a slice of all the field errors.

So in a real app scenario where you know that struct tags are not placed incorrectly you can confidently type cast the error to `GroupedValidationError` and use the `ToJsonStructList` function to have a `[]ValidationErrorJson` which can be used to format the errors prettily for a json response, or you could look over the slice of `FieldValidateError`s and use the methods present to make your own json formatter.
 
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

