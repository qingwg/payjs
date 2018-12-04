package context

import "net/http"

// Context struct
type Context struct {
	Key       string
	MchID     string
	NotifyURL string

	Writer  http.ResponseWriter
	Request *http.Request
}
