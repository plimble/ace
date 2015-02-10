# ACE

web framework by golang

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
	c.String(200,"Hello ACE")
}
a.Run(":8080")
```


#### Use Middlewares
support middlewares logger , recovery , crsf

```
a := ace.New()
a.Use(ace.Logger())
```


#### Create ACE default builtin middlewares logger & recovery

```
a := ace.Default()
```
#### Use GET, POST, PUT, PATCH, DELETE and OPTIONS

```
a.DELETE("/", myFunc)
a.HEAD("/", myFunc)
a.OPTIONS("/", myFunc)
a.PATCH("/", myFunc)
a.PUT("/", myFunc)
a.POST("/", myFunc)
a.GET("/", myFunc)
```

#### Get param & Response json

```
a.GET("/hello/:name", func(c *ace.C) { 
	name := c.Params.ByName("name")
	c.JSON(200, map[string]string{"hello": name})
}
```

#### Render html template

```
at := ace.TemplateOptions{Directory: "./web", IsDevelopment: true}
a.UseTemplate(&at)

a.GET("/", func(c *ace.C) {
	c.Data = map[string]interface{}{"name": "john"}
	c.HTML("index.html")
})
```

index.html

```
Hello ACE {{name}}
```

#### Static file

```
a.Static("/assets", "./img")
```

#### Group router

```
mygroup := a.Group("/people", func(c *ace.C) {
	fmt.Println("before route func")
	c.Next()
})

mygroup.GET("/", func(c *ace.C) {
	c.JSON(200, map[string]string{"TEST": "GET METHOD"})
})

mygroup.POST("/", func(c *ace.C) {
	c.JSON(200, map[string]string{"TEST": "POST METHOD"})
})
```

#### Session with gorilla

```
var store = sessions.NewCookieStore([]byte("something-very-secret"))
a.UseSession("cookie", store, nil)

```

Example Get , Set

```
a.GET("/cookie", func(c *ace.C) {
	c.Session.SetString("imString", "TESTSTRING")
	fmt.Println(c.Session.GetString("imString"))
}
```