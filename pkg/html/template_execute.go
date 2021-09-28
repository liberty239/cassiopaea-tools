package html

import (
	"io"
	"text/template"
)

func TemplateExecute(tpl *template.Template, w io.Writer, data interface{}) error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					err = r
				default:
					panic(r)
				}
			}
		}()
		return tpl.Execute(w, data)
	}()
}
