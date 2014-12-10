package ace

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"math"
	"net/http"
	"strings"
)

const (
	ContentType    = "Content-Type"
	AcceptLanguage = "Accept-Language"
)

const (
	AbortIndex   = math.MaxInt8 / 2
	MIMEJSON     = "application/json"
	MIMEHTML     = "text/html"
	MIMEXML      = "application/xml"
	MIMEXML2     = "text/xml"
	MIMEPlain    = "text/plain"
	MIMEPOSTForm = "application/x-www-form-urlencoded"
)

type HTMLOptions struct {
	Layout string
}

type C struct {
	writercache         responseWriter
	Params              httprouter.Params
	Request             *http.Request
	Writer              ResponseWriter
	render              *render.Render
	index               int8
	handlers            []HandlerFunc
	notfoundHandlerFunc HandlerFunc
	failHandlerFunc     HandlerFunc
	//recovery
	Recovery interface{}
	context  map[string]interface{}
}

func (a *Ace) CreateContext(w http.ResponseWriter, r *http.Request) *C {
	context := a.pool.Get().(*C)
	context.writercache.reset(w)
	context.Writer = &context.writercache
	context.Request = r
	context.context = make(map[string]interface{})

	return context
}

func (c *C) header(status int, ct string) {
	c.Writer.Header().Set(ContentType, "application/json")
	c.Writer.WriteHeader(status)
}

func (c *C) JSON(status int, v interface{}) {
	c.render.JSON(c.Writer, status, v)
}

func (c *C) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	if len(htmlOpt) == 0 {
		c.render.HTML(c.Writer, status, name, binding)
	} else {
		c.render.HTML(c.Writer, status, name, binding, render.HTMLOptions(htmlOpt[0]))
	}
}

func (c *C) XML(status int, v interface{}) {
	c.render.XML(c.Writer, status, v)
}

func (c *C) Data(status int, v []byte) {
	c.render.Data(c.Writer, status, v)
}

func (c *C) String(status int, format string, val ...interface{}) {
	c.header(status, "text/plain")
	if len(val) == 0 {
		c.Writer.Write([]byte(format))
	} else {
		c.Writer.Write([]byte(fmt.Sprintf(format, val...)))
	}
}

func (c *C) HTTPLang() string {
	langStr := c.Request.Header.Get(AcceptLanguage)
	return strings.Split(langStr, ",")[0]
}

func (c *C) Redirect(url string) {
	http.Redirect(c.Writer, c.Request, url, 302)
}

func (c *C) Abort(status int) {
	c.Writer.WriteHeader(status)
	c.index = AbortIndex
}

func (c *C) NotFound() {
	c.notfoundHandlerFunc(c)
}

func (c *C) Fail() {
	c.failHandlerFunc(c)
}

func (c *C) Next() {
	c.index++
	s := int8(len(c.handlers))
	if c.index < s {
		c.handlers[c.index](c)
	}
}

func (c *C) ClientIP() string {
	clientIP := c.Request.Header.Get("X-Real-IP")
	if len(clientIP) == 0 {
		clientIP = c.Request.Header.Get("X-Forwarded-For")
	}
	if len(clientIP) == 0 {
		clientIP = c.Request.RemoteAddr
	}
	return clientIP
}

func (c *C) Bind(obj interface{}) bool {
	var b binding.Binding
	ctype := c.Request.Header.Get("Content-Type")
	switch {
	case c.Request.Method == "GET" || ctype == MIMEPOSTForm:
		b = binding.Form
	case ctype == MIMEJSON:
		b = binding.JSON
	case ctype == MIMEXML || ctype == MIMEXML2:
		b = binding.XML
	default:
		c.String(400, "unknown content-type: "+ctype)
		return false
	}
	return c.BindWith(obj, b)
}

func (c *C) BindWith(obj interface{}, b binding.Binding) bool {
	if err := b.Bind(c.Request, obj); err != nil {
		c.String(400, err.Error())
		return false
	}
	return true
}

func (c *C) Set(key string, v interface{}) {
	c.context[key] = v
}

func (c *C) Get(key string) interface{} {
	return c.context[key]
}
