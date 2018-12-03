package context

import (
	"net/http"
	"sync"
)

// Context struct
type Context struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	PayMchID       string
	PayNotifyURL   string
	PayKey         string

	Cache cache.Cache

	Writer  http.ResponseWriter
	Request *http.Request

	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex

	//jsAPITicket 读写锁 同一个AppID一个
	jsAPITicketLock *sync.RWMutex
}
