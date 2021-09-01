package server

import (
	"errors"
	"log"
	"net/http"
)

type Server struct {
	srv *http.Server
}

// Get returns an instance of Server which has a pointer to an http server
func Get() *Server {
	return &Server{
		srv: &http.Server{},
	}
}

// WithErrorLogger adds a custom logger to use as http server the error log
func (s *Server) WithErrorLogger(l *log.Logger) *Server {
	s.srv.ErrorLog = l
	return s
}

// WithAddr sets the port for the server to listen on
func (s *Server) WithAddr(address string) *Server {
	s.srv.Addr = address
	return s
}

// Listen starts the ListenAndServe method of the http server
func (s *Server) Listen() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("server is missing an address port")
	}
	// TODO: Add a check to see if a router is passed and throw error if not
	return s.srv.ListenAndServe()
}
