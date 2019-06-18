package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s := &Server{
		server: httpServer,
	}
	return s
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("cannot start http server %s", err.Error())
		} else {
			log.Print(err.Error())
		}
	}()
}

func (s *Server) Stop() {
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
