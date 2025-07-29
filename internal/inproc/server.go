package inproc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
)

type inMemoryServer struct {
	server   *httptest.Server
	listener *memoryListener
}

type memoryListener struct {
	conns  chan net.Conn
	once   sync.Once
	closed chan struct{}
}

// NewInMemoryServer constructs and starts an inMemoryServer.
func NewInMemoryServer(handler http.Handler) *inMemoryServer {
	lis := &memoryListener{
		conns:  make(chan net.Conn),
		closed: make(chan struct{}),
	}
	server := httptest.NewUnstartedServer(handler)
	server.Listener = lis
	server.EnableHTTP2 = true
	server.StartTLS()
	return &inMemoryServer{
		server:   server,
		listener: lis,
	}
}

// Client returns an HTTP client configured to trust the server's TLS
// certificate and use HTTP/2 over an in-memory pipe. Automatic HTTP-level gzip
// compression is disabled. It closes its idle connections when the server is
// closed.
func (s *inMemoryServer) Client() *http.Client {
	client := s.server.Client()
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.DialContext = s.listener.DialContext
		transport.DisableCompression = true
	}
	return client
}

// URL is the server's URL.
func (s *inMemoryServer) URL() string {
	return s.server.URL
}

// Close shuts down the server, blocking until all outstanding requests have
// completed.
func (s *inMemoryServer) Close() {
	s.server.Close()
}

// Accept implements net.Listener.
func (l *memoryListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.conns:
		return conn, nil
	case <-l.closed:
		return nil, errors.New("listener closed")
	}
}

func (l *memoryListener) Close() error {
	l.once.Do(func() {
		close(l.closed)
	})
	return nil
}

func (l *memoryListener) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	select {
	case <-l.closed:
		return nil, errors.New("listener closed")
	default:
	}
	server, client := net.Pipe()
	l.conns <- server
	return client, nil
}

// Addr implements net.Listener.
func (l *memoryListener) Addr() net.Addr {
	return &memoryAddr{}
}

type memoryAddr struct{}

// Network implements net.Addr.
func (*memoryAddr) Network() string { return "memory" }

// String implements io.Stringer, returning a value that matches the
// certificates used by net/http/httptest.
func (*memoryAddr) String() string { return "example.com" }
