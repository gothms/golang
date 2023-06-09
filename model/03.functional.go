package model

import (
	"crypto/tls"
	"time"
)

/*
	Builder模式：使用 ServerBuilder 结构体，包装 Server
	Functional Options：
		省略包装结构体 ServerBuilder
*/

// Server https://time.geekbang.com/column/article/330218?utm_source=pc_cp&utm_term=pc_interstitial_1346
type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}
type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}
func Timeout(t time.Duration) Option {
	return func(s *Server) {
		s.Timeout = t
	}
}
func MaxConns(mc int) Option {
	return func(s *Server) {
		s.MaxConns = mc
	}
}
func TLS(tls *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}
func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {
	srv := &Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}
	for _, option := range options {
		option(srv)
	}
	return srv, nil
}
