package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handlers http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handlers,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1024,
	}
	fmt.Printf("Running at :%s\n", port)
	err := s.httpServer.ListenAndServe()
	return err

}
