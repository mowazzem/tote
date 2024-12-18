package router

import "net/http"

type middleware func(http.HandlerFunc) http.HandlerFunc

type router struct {
	*http.ServeMux
	middlewares []middleware
}

func New() *router {
	return &router{
		ServeMux: http.NewServeMux(),
	}
}

func (r *router) Sub() *router {
	return &router{
		ServeMux: r.ServeMux,
	}
}

func (r *router) Use(m ...middleware) {
	if r.middlewares == nil {
		r.middlewares = make([]middleware, 0, len(m))
	}
	r.middlewares = append(r.middlewares, m...)
}

func (r *router) HandleFunc(pattern string, handler http.HandlerFunc) {
	for _, m := range r.middlewares {
		handler = m(handler)
	}
	r.ServeMux.HandleFunc(pattern, handler)
}
