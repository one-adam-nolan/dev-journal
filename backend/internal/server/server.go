package server

import (
	"net/http"

	"connectrpc.com/connect"

	"dj/backend/internal/app"
	djv1connect "dj/gen/go/dj/v1/djv1connect"
)

// NewHandler returns an HTTP handler that serves the JournalService Connect API.
func NewHandler(application *app.App) http.Handler {
	mux := http.NewServeMux()

	path, handler := djv1connect.NewJournalServiceHandler(
		application.Connect,
		connect.WithCompressMinBytes(1024),
	)
	mux.Handle(path, handler)

	return mux
}
