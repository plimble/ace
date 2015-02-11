package ace

import (
	"net/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}

func (a *Ace) UseHtmlTemplate(render Renderer) {
	a.render = render
}
