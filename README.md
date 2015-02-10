ACE [![godoc badge](http://godoc.org/github.com/plimble/ace?status.png)](http://godoc.org/github.com/plimble/ace)   [![gocover badge](http://gocover.io/_badge/github.com/plimble/ace?t=1)](http://gocover.io/github.com/plimble/ace) [![Build Status](https://api.travis-ci.org/plimble/ace.svg?branch=master&t=1)](https://travis-ci.org/plimble/ace) [![Go Report Card](http://goreportcard.com/badge/plimble/ace?t=1)](http:/goreportcard.com/report/plimble/ace)
========

Blazing fast Go Web Framework

####Installation

```
go get github.com/plimble/ace
```

#### Import

```
import "github.com/plimble/ace"
```

##Usage



#### Quick Start

```
a := ace.New()
a.GET("/", func(c *ace.C) {
	name := c.Params.ByName("name")
	c.JSON(200, map[string]string{"hello": name})
}
a.Run(":8080")
```

Default Middleware (Logger, Recovery)
```
a := ace.Default()
a.GET("/", func(c *ace.C) {
	c.String(200,"Hello ACE")
}
a.Run(":8080")
```

### Router
```
a.DELETE("/", HandlerFunc)
a.HEAD("/", HandlerFunc)
a.OPTIONS("/", HandlerFunc)
a.PATCH("/", HandlerFunc)
a.PUT("/", HandlerFunc)
a.POST("/", HandlerFunc)
a.GET("/", HandlerFunc)
```
##### Example
```
	a := ace.New()

	a.Get("/", func(c *ace.C){
		c.String(200, "Hello world")
	})

	a.POST("/form/:id/:name", func(c *ace.C){
		id := c.Params.ByName("id")
		name := c.Params.ByName("name")
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
```
	//application/octet-stream
	c.Download(200, []byte("Hello Ace"))
```
##### HTML
```
	c.HTML("index.html")
```
##### Redirect
```
	c.Redirect("/home")
```

## Group Router

```
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


## Middlewares
Ace middleware is implemented by custom handler
```
type HandlerFunc func(c *C)
```
#####Example
```
a := ace.New()
a.Use(ace.Logger())
```

### Built-in Middlewares

##### Serve Static
```
a.Static("/assets", "./img")
```

##### Session with Gorilla sessions

```
var store = sessions.NewCookieStore([]byte("something-very-secret"))
a.UseSession("cookie", store, nil)

```

```
a.GET("/hello", func(c *ace.C) {
	c.Session.SetString("name", "John Doe")
	fmt.Println(c.Session.GetString("name"))
}
```
##### Logger
```
a.Use(ace.Logger())
```

##### Recovery
```
a.Use(ace.Recovery())
```

##### CSRF
```
options := &ace.CSRFOptions{
	FailedHandler: func(c *ace.C){
		c.String(500, "Invalid CSRF Token")
	}
}
a.Use(ace.CSRF(options))
```

## HTML Template Engine
Ace use [Pongo2](https://github.com/flosch/pongo2) as a default template engine

```
options := &ace.TemplateOptions{
	Directory: "./web",
	IsDevelopment: true,
	Extensions: []string{},
}
a.UseTemplate(at)

a.GET("/", func(c *ace.C) {
	c.Data["fname"] = "John"
	c.Data["lname"] = "Doe"
	c.HTML("index.html")
})
```

#####index.html

```
Hello ACE {{fname}} {{lname}}
```







###Contributing

If you'd like to help out with the project. You can put up a Pull Request.