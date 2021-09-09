package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	srv *http.Server
}

// Get returns an instance of Server which has a pointer to an http server
func Get() *Server {
	return &Server{
		srv: &http.Server{
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
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

// WithRouter adds a mux router to the server handler
func (s *Server) WithRouter(r *mux.Router) *Server {
	s.srv.Handler = r
	return s
}

// Listen starts the ListenAndServe method of the http server
func (s *Server) Listen() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("server is missing an address port")
	}

	if s.srv.Handler == nil {
		return errors.New("server is missing a mux router")
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Close() error {
	return s.srv.Close()
}
