package httprouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/leomirandadev/clean-architecture-go/pkg/customerr"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
)

type chiRouter struct {
	router *chi.Mux
	server *http.Server
}

func NewChiRouter() Router {
	return &chiRouter{
		router: chi.NewRouter(),
	}
}

func (r chiRouter) POST(uri string, f HandlerFunc) {
	r.router.Post(uri, func(w http.ResponseWriter, r *http.Request) {
		doer(f, w, r)
	})
}

func (r chiRouter) GET(uri string, f HandlerFunc) {
	r.router.Get(uri, func(w http.ResponseWriter, r *http.Request) {
		doer(f, w, r)
	})
}

func (r chiRouter) PUT(uri string, f HandlerFunc) {
	r.router.Put(uri, func(w http.ResponseWriter, r *http.Request) {
		doer(f, w, r)
	})
}

func (r chiRouter) PATCH(uri string, f HandlerFunc) {
	r.router.Patch(uri, func(w http.ResponseWriter, r *http.Request) {
		doer(f, w, r)
	})
}

func (r chiRouter) DELETE(uri string, f HandlerFunc) {
	r.router.Delete(uri, func(w http.ResponseWriter, r *http.Request) {
		doer(f, w, r)
	})
}

func doer(f HandlerFunc, w http.ResponseWriter, r *http.Request) {
	spanName := fmt.Sprintf("[%s] %s", r.Method, r.RequestURI)
	ctx, tr := tracer.Span(r.Context(), spanName)
	defer tr.End()

	// update request context
	r = r.WithContext(ctx)

	// default header json
	w.Header().Set("Content-Type", "application/json")

	// execute
	chiCtx := newChiContext(w, r)
	if err := f(chiCtx); err != nil {
		status := customerr.StatusCode(err)
		chiCtx.JSON(status, err)
		return
	}
}

func (r *chiRouter) Serve(port string) {
	fmt.Println("Online now in: http://localhost" + port)

	r.server = &http.Server{
		Addr:    port,
		Handler: r.router,
	}

	log.Fatal(r.server.ListenAndServe())
}

func (r chiRouter) Shutdown(ctx context.Context) error {
	if r.server == nil {
		return errors.New("server not initialized")
	}

	return r.server.Shutdown(ctx)
}

func (m chiRouter) ParseHandler(h http.HandlerFunc) HandlerFunc {
	return func(c Context) error {
		h(c.GetResponseWriter(), c.GetRequestReader())
		return nil
	}
}

type chiContext struct {
	w http.ResponseWriter
	r *http.Request
}

func newChiContext(w http.ResponseWriter, r *http.Request) Context {
	return chiContext{w, r}
}

func (c chiContext) Context() context.Context {
	return c.r.Context()
}

func (c chiContext) JSON(status int, data any) error {
	c.w.WriteHeader(status)
	return json.NewEncoder(c.w).Encode(data)
}

func (c chiContext) GetPathParam(param string) string {
	return chi.URLParam(c.r, param)
}

func (c chiContext) GetQueryParam(param string) string {
	paramResult := c.r.URL.Query()[param]
	if len(paramResult) == 0 {
		return ""
	}
	return paramResult[0]
}

func (c chiContext) GetFromHeader(param string) string {
	return c.r.Header.Get(param)
}

func (c chiContext) Decode(v any) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c chiContext) Headers() http.Header {
	return c.r.Header
}

func (c chiContext) GetResponseWriter() http.ResponseWriter {
	return c.w
}

func (c chiContext) GetRequestReader() *http.Request {
	return c.r
}
