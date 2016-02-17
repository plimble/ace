ACE [![godoc badge](http://godoc.org/github.com/plimble/ace?status.svg)](http://godoc.org/github.com/plimble/ace)   [![gocover badge](http://gocover.io/_badge/github.com/plimble/ace?t=3)](http://gocover.io/github.com/plimble/ace) [![Build Status](https://api.travis-ci.org/plimble/ace.svg?branch=master&t=3)](https://travis-ci.org/plimble/ace) [![Go Report Card](http://goreportcard.com/badge/plimble/ace?t=3)](http:/goreportcard.com/report/plimble/ace)
========

Blazing fast Go Web Framework

![image](http://image.free.in.th/v/2013/id/150218064526.jpg)

####Installation

```
go get github.com/plimble/ace
```

#### Import

```go
import "github.com/plimble/ace"
```

## Performance
Ace is very fast you can see on [this](https://gist.github.com/witooh/1c05c71d9510b2020e48)

## Usage

#### Quick Start

```go
a := ace.New()
a.GET("/:name", func(c *ace.C) {
	name := c.Param("name")
	c.JSON(200, map[string]string{"hello": name})
})
a.Run(":8080")
```

Default Middleware (Logger, Recovery)

```go
a := ace.Default()
a.GET("/", func(c *ace.C) {
	c.String(200,"Hello ACE")
})
a.Run(":8080")
```

### Router

```go
a.DELETE("/", HandlerFunc)
a.HEAD("/", HandlerFunc)
a.OPTIONS("/", HandlerFunc)
a.PATCH("/", HandlerFunc)
a.PUT("/", HandlerFunc)
a.POST("/", HandlerFunc)
a.GET("/", HandlerFunc)
```

##### Example

```go
a := ace.New()

a.GET("/", func(c *ace.C){
	c.String(200, "Hello world")
})

a.POST("/form/:id/:name", func(c *ace.C){
	id := c.Param("id")
	name := c.Param("name")
	age := c.Request.PostFormValue("age")
})
```

## Response

##### JSON

```go
data := struct{
	Name string `json:"name"`
}{
	Name: "John Doe",
}
c.JSON(200, data)
```

##### String

```go
c.String(200, "Hello Ace")
```

##### Download

```go
//application/octet-stream
c.Download(200, []byte("Hello Ace"))
```

##### HTML

```go
c.HTML("index.html")
```

##### Redirect

```go
c.Redirect("/home")
```

## Group Router

```go
g := a.Group("/people", func(c *ace.C) {
	fmt.Println("before route func")
	c.Next()
})

// /people/:name
g.GET("/:name", func(c *ace.C) {
	c.JSON(200, map[string]string{"TEST": "GET METHOD"})
})

// /people/:name
g.POST("/:name", func(c *ace.C) {
	c.JSON(200, map[string]string{"TEST": "POST METHOD"})
})
```

## Data

Set/Get data in any HandlerFunc

```go
a.Use(func(c *ace.C){
	c.SetData("isLogin", true)
})

a.Get("/", func(c *ace.C){
	isLogin := c.GetData("isLogin").(bool)
	//or get all data
	//c.GetAllData()
})
```

## Get Post Form and Query

```go
a.Get("/", func(c *ace.C){
	name := c.MustPostString(key, default_value)
	age := c.MustPostInt(key, d)

	q := c.MustQueryString(key, default_value)
	score := c.MustQueryFloat64(key, default_value)
})
```

## Get data From JSON Request

```go
a.Get("/", func(c *ace.C){
	user := struct{
		Name string `json:"name"`
	}{}

	c.ParseJSON(&user)
})
```

## Panic Response

Use panic instead of `if err != nil` for response error

```go
a.Get("/save", func(c *ace.C){
	user := &User{}

	c.ParseJSON(user)

	//this func return error
	//if error go to panic handler
	c.Panic(doSomething1(user))
	c.Panic(doSomething2(user))

	c.String(201, "created")
}

a.Get("/get", func(c *ace.C){
	id := c.Param("id")

	user, err := doSomeThing()
	//if error go to panic handler
	c.Panic(err)

	c.JSON(200, user)
}
```

#### Custom panic response

```go
a := ace.New()
a.Panic(func(c *ace.C, rcv interface{}){
	switch err := rcv.(type) {
		case error:
			c.String(500, "%s\n%s", err, ace.Stack())
		case CustomError:
			log.Printf("%s\n%s", err, ace.Stack())
			c.JSON(err.Code, err.Msg)
	}
})
```


## Middlewares

Ace middleware is implemented by custom handler

```go
type HandlerFunc func(c *C)
```

##### Example

```go
a := ace.New()
a.Use(ace.Logger())
```

### Built-in Middlewares

##### Serve Static

```go
a.Static("/assets", "./img")
```

##### Session

You can use store from [sessions](https://github.com/plimble/sessions)

```go
import github.com/plimble/sessions/store/cookie

a := ace.New()

store := cookie.NewCookieStore()
a.Use(ace.Session(store, nil))

a.GET("/", func(c *ace.C) {
	//get session name
	session1 := c.Sessions("test")
	session1.Set("test1", "123")
	session1.Set("test2", 123)

	session2 := c.Sessions("foo")
	session2.Set("baz1", "123")
	session2.Set("baz2", 123)

	c.String(200, "")
})

a.GET("/test", func(c *C) {
	session := c.Sessions("test")
	//get value from key test1 if not found default value ""
	test1 := session.GetString("test1", "")
	test2 := session.GetInt("test2", 0)

	c.String(200, "")
})
```

##### Logger

```go
a.Use(ace.Logger())
```

## HTML Template Engine

Ace built on renderer interface. So you can use any template engine

```go
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}
```

### ACE Middlewares

| Name                                                	| Description                                 	|
|-----------------------------------------------------	|---------------------------------------------	|
| [gzip](https://github.com/plimble/ace-contrib/tree/master/gzip)         	| GZIP compress                               	|
| [cors](https://github.com/plimble/ace-contrib/tree/master/cors)         	| Enable Cross-origin resource sharing (CORS) 	|
| [sessions](https://github.com/plimble/sessions) 													| Sessions      				                      	|
| [pongo2](https://github.com/plimble/ace-contrib/tree/master/pongo2)     	| Pongo2 Template Engine                      	|
| [csrf](https://github.com/plimble/ace-contrib/tree/master/csrf)         	| Cross Site Request Forgery protection       	|

### Contributing

If you'd like to help out with the project. You can put up a Pull Request.
