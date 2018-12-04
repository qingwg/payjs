package context

import (
	"net/http"
)

// Context struct
type Context struct {
	Key       string
	MchID     string
	NotifyURL string

	Writer  http.ResponseWriter
	Request *http.Request
}

// Query returns the keyed url query value if it exists
func (ctx *Context) Query(key string) string {
	value, _ := ctx.GetQuery(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) GetQuery(key string) (string, bool) {
	req := ctx.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
