package ace

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/plimble/sessions"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	contentType    = "Content-Type"
	acceptLanguage = "Accept-Language"
	abortIndex     = math.MaxInt8 / 2
)

//C is context for every goroutine
type C struct {
	writercache responseWriter
	params      httprouter.Params
	Request     *http.Request
	Writer      ResponseWriter
	index       int8
	handlers    []HandlerFunc
	data        map[string]interface{}
	sessions    *sessions.Sessions
	render      Renderer
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
	c.Writer.Header().Set(contentType, "application/json; charset=UTF-8")
	c.Writer.WriteHeader(status)
	if v == nil {
		return
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		panic(err)
	}

	c.Writer.Write(buf.Bytes())
}

//String response with text/html; charset=UTF-8 Content type
func (c *C) String(status int, format string, val ...interface{}) {
	c.Writer.Header().Set(contentType, "text/html; charset=UTF-8")
	c.Writer.WriteHeader(status)

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if len(val) == 0 {
		buf.WriteString(format)
	} else {
		buf.WriteString(fmt.Sprintf(format, val...))
	}

	c.Writer.Write(buf.Bytes())
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

//Param get param from route
func (c *C) Param(name string) string {
	return c.params.ByName(name)
}

//ParseJSON decode json to interface{}
func (c *C) ParseJSON(v interface{}) {
	defer c.Request.Body.Close()
	c.Panic(json.NewDecoder(c.Request.Body).Decode(v))
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
func (c *C) Set(key string, v interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = v
}

//SetAll data
func (c *C) SetAll(data map[string]interface{}) {
	c.data = data
}

//Get data
func (c *C) Get(key string) interface{} {
	return c.data[key]
}

//GetAllData return all data
func (c *C) GetAll() map[string]interface{} {
	return c.data
}

func (c *C) Sessions(name string) *sessions.Session {
	return c.sessions.Get(name)
}

func (c *C) MustQueryInt(key string, d int) int {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}

	return i
}

func (c *C) MustQueryFloat64(key string, d float64) float64 {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(c.Request.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (c *C) MustQueryString(key, d string) string {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}

	return val
}

func (c *C) MustQueryStrings(key string, d []string) []string {
	val := c.Request.URL.Query()[key]
	if len(val) == 0 {
		return d
	}

	return val
}

func (c *C) MustQueryTime(key string, layout string, d time.Time) time.Time {
	val := c.Request.URL.Query().Get(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, c.Request.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}

	return t
}

/////////////////////////

func (c *C) MustPostInt(key string, d int) int {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err.Error())
	}

	return i
}

func (c *C) MustPostFloat64(key string, d float64) float64 {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	f, err := strconv.ParseFloat(c.Request.URL.Query().Get(key), 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (c *C) MustPostString(key, d string) string {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}

	return val
}

func (c *C) MustPostStrings(key string, d []string) []string {
	if c.Request.PostForm == nil {
		c.Request.ParseForm()
	}

	val := c.Request.PostForm[key]
	if len(val) == 0 {
		return d
	}

	return val
}

func (c *C) MustPostTime(key string, layout string, d time.Time) time.Time {
	val := c.Request.PostFormValue(key)
	if val == "" {
		return d
	}
	t, err := time.Parse(layout, c.Request.URL.Query().Get(key))
	if err != nil {
		panic(err)
	}

	return t
}

func (c *C) Panic(err error) {
	if err != nil {
		panic(err)
	}
}
