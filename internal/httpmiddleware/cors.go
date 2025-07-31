package httpmiddleware

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
)

// CORS adds CORS support to a Connect HTTP handler.
func CORS(h http.Handler) http.Handler {
	connectHeaders := connectcors.AllowedHeaders()
	var headers []string
	headers = append(headers, connectHeaders...)

	middleware := cors.New(cors.Options{
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: headers,
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
