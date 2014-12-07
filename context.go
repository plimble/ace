package copter

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"strings"
)

const (
	ContentType    = "Content-Type"
	AcceptLanguage = "Accept-Language"
)

type HTMLOptions struct {
	Layout string
}

type C struct {
	Params   httprouter.Params
	Req      *http.Request
	Res      http.ResponseWriter
	render   *render.Render
	index    int8
	handlers []HandlerFunc
}

func (c *C) header(status int, ct string) {
	c.Res.Header().Set(ContentType, "application/json")
	c.Res.WriteHeader(status)
}

func (c *C) JSON(status int, v interface{}) {
	c.render.JSON(c.Res, status, v)
}

func (c *C) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	if len(htmlOpt) == 0 {
		c.render.HTML(c.Res, status, name, binding)
	} else {
		c.render.HTML(c.Res, status, name, binding, render.HTMLOptions(htmlOpt[0]))
	}
}

func (c *C) XML(status int, v interface{}) {
	c.render.XML(c.Res, status, v)
}

func (c *C) Data(status int, v []byte) {
	c.render.Data(c.Res, status, v)
}

func (c *C) String(status int, format string, val ...interface{}) {
	c.header(status, "text/plain")
	c.Res.Write([]byte(fmt.Sprintf(format, val)))
}

func (c *C) HTTPLang() string {
	langStr := c.Req.Header.Get(AcceptLanguage)
	return strings.Split(langStr, ",")[0]
}

func (c *C) Redirect(url string) {
	http.Redirect(c.Res, c.Req, url, 302)
}

func (c *C) About(status int) {
	c.Res.WriteHeader(status)
	c.index = 127
}

func (c *C) Next() {
	c.index++
	s := int8(len(c.handlers))
	if c.index < s {
		c.handlers[c.index](c)
	}
}
