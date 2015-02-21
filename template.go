package ace

import (
	"net/http"
)

//Renderer html render interface
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}

//HtmlTemplate use html template middleware
func (a *Ace) HtmlTemplate(render Renderer) {
	a.render = render
}
