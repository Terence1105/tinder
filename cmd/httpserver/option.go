package httpserver

import (
	"net/http"
)

const (
	HttpServerPortFlag       = "port"
	HttpServerPortVal        = 80
	HttpServerhandlerTypeGin = "gin"
)

type Options struct {
	port int
	tls  bool

	handler http.Handler
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		port: HttpServerPortVal,
		tls:  false,
	}
}

func ServeMux(h http.Handler) Option {
	return func(o *Options) {
		o.handler = h
	}
}
