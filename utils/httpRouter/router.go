package httpRouter

import "net/http"

type Router interface {
	GET(uri string, f http.HandlerFunc)
	POST(uri string, f http.HandlerFunc)
	PUT(uri string, f http.HandlerFunc)
	DELETE(uri string, f http.HandlerFunc)
	PATCH(uri string, f http.HandlerFunc)
	SERVE(port string)
}
