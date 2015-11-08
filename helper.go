package cmdns

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs template.FuncMap = template.FuncMap{
	"trim":           strings.TrimSpace,
	"trimRightSpace": trimRightSpace,
	"rpad":           rpad,
	"gt":             cobra.Gt,
	"eq":             cobra.Eq,
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

//rpad adds padding to the right of a string
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
