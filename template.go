package ace

import (
	"github.com/plimble/copter"
)

type TemplateOptions struct {
	Directory     string
	Extensions    []string
	IsDevelopment bool
}

func (a *Ace) UseTemplate(options *TemplateOptions) {
	a.render = copter.New(&copter.Options{
		Directory:     options.Directory,
		Extensions:    options.Extensions,
		IsDevelopment: options.IsDevelopment,
	})
}
