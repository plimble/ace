package ace

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"math"
	"net/http"
	"strings"
)

const (
	ContentType    = "Content-Type"
	AcceptLanguage = "Accept-Language"
	AbortIndex     = math.MaxInt8 / 2
)

type C struct {
	writercache      responseWriter
	Params           httprouter.Params
	Request          *http.Request
	Writer           ResponseWriter
	index            int8
	handlers         []HandlerFunc
	errorHandlerFunc ErrorHandlerFunc
	//recovery
	context map[string]interface{}
	err     error
}

func (a *Ace) CreateContext(w http.ResponseWriter, r *http.Request) *C {
	c := a.pool.Get().(*C)
	c.writercache.reset(w)
	c.Request = r
	c.context = nil
	c.index = -1

	return c
}

func (c *C) JSON(status int, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(status)
	c.Writer.Header().Set(ContentType, "application/json; charset=UTF-8")
	c.Writer.Write(result)
}

func (c *C) String(status int, format string, val ...interface{}) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Set(ContentType, "text/html; charset=UTF-8")
	if len(val) == 0 {
		c.Writer.Write([]byte(format))
	} else {
		c.Writer.Write([]byte(fmt.Sprintf(format, val...)))
	}
}

func (c *C) Data(status int, v []byte) {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Set(ContentType, "application/octet-stream; charset=UTF-8")
	c.Writer.Write(v)
}

func (c *C) ParseJSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

func (c *C) ParseForm(v interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	return decoder.Decode(v, c.Request.PostForm)
}

func (c *C) ParseMultipartForm(v interface{}, maxMemory int64) error {
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	return decoder.Decode(v, c.Request.PostForm)
}

func (c *C) ParseQuery(v interface{}) error {
	decoder := schema.NewDecoder()
	return decoder.Decode(v, c.Request.URL.Query())
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

func (c *C) Error(err error) {
	c.err = err
	c.errorHandlerFunc(c, err)
	c.index = AbortIndex
}

func (c *C) GetError() error {
	return c.err
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

func (c *C) Set(key string, v interface{}) {
	if c.context == nil {
		c.context = make(map[string]interface{})
	}
	c.context[key] = v
}

func (c *C) Get(key string) interface{} {
	return c.context[key]
}
