package ace

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math"
	"net/http"
	"strings"
)

const (
	contentType    = "Content-Type"
	acceptLanguage = "Accept-Language"
	abortIndex     = math.MaxInt8 / 2
)

//C is context for every goroutine
type C struct {
	writercache responseWriter
	Params      httprouter.Params
	Request     *http.Request
	Writer      ResponseWriter
	index       int8
	handlers    []HandlerFunc
	//recovery
	data    map[string]interface{}
	Session *session
	render  Renderer
}

func (a *Ace) createContext(w http.ResponseWriter, r *http.Request) *C {
	c := a.pool.Get().(*C)
	c.writercache.reset(w)
	c.Request = r
	c.index = -1
	c.data = nil
	c.render = a.render

	return c
}

//JSON response with application/json; charset=UTF-8 Content type
func (c *C) JSON(status int, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Writer.Header().Set(contentType, "application/json; charset=UTF-8")
	c.Writer.WriteHeader(status)
	c.Writer.Write(result)
}

//String response with text/html; charset=UTF-8 Content type
func (c *C) String(status int, format string, val ...interface{}) {
	c.Writer.Header().Set(contentType, "text/html; charset=UTF-8")
	c.Writer.WriteHeader(status)
	if len(val) == 0 {
		c.Writer.Write([]byte(format))
	} else {
		c.Writer.Write([]byte(fmt.Sprintf(format, val...)))
	}
}

//Download response with application/octet-stream; charset=UTF-8 Content type
func (c *C) Download(status int, v []byte) {
	c.Writer.Header().Set(contentType, "application/octet-stream; charset=UTF-8")
	c.Writer.WriteHeader(status)
	c.Writer.Write(v)
}

//HTML render template engine
func (c *C) HTML(name string, data interface{}) {
	c.render.Render(c.Writer, name, data)
}

//ParseJSON decode json to interface{}
func (c *C) ParseJSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

//HTTPLang get first language from HTTP Header
func (c *C) HTTPLang() string {
	langStr := c.Request.Header.Get(acceptLanguage)
	return strings.Split(langStr, ",")[0]
}

//Redirect 302 response
func (c *C) Redirect(url string) {
	http.Redirect(c.Writer, c.Request, url, 302)
}

//Abort stop maddileware
func (c *C) Abort() {
	c.index = abortIndex
}

//AbortWithStatus stop maddileware and return http status code
func (c *C) AbortWithStatus(status int) {
	c.Writer.WriteHeader(status)
	c.Abort()
}

//Next next middleware
func (c *C) Next() {
	c.index++
	s := int8(len(c.handlers))
	if c.index < s {
		c.handlers[c.index](c)
	}
}

//ClientIP get ip from RemoteAddr
func (c *C) ClientIP() string {
	return c.Request.RemoteAddr
}

//Set data
func (c *C) SetData(key string, v interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = v
}

//Get data
func (c *C) GetData(key string) interface{} {
	return c.data[key]
}

//GetAllData return all data
func (c *C) GetAllData() map[string]interface{} {
	return c.data
}
