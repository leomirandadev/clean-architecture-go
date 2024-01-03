package httprouter

import (
	"context"
	"net/http"
)

type Router interface {
	GET(uri string, f HandlerFunc)
	POST(uri string, f HandlerFunc)
	PUT(uri string, f HandlerFunc)
	DELETE(uri string, f HandlerFunc)
	PATCH(uri string, f HandlerFunc)
	ParseHandler(h http.HandlerFunc) HandlerFunc
	SERVE(port string)
}

type HandlerFunc func(c Context) error

type Context interface {
	Context() context.Context
	JSON(status int, data any) error
	Decode(v any) error
	GetPathParam(param string) string
	GetQueryParam(param string) string
	GetFromHeader(param string) string
	Headers() http.Header
	GetResponseWriter() http.ResponseWriter
	GetRequestReader() *http.Request
}
