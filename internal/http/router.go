package http

import (
	stdhttp "net/http"
)

func NewMux() *stdhttp.ServeMux {
	mux := stdhttp.NewServeMux()
	// later: register handlers
	return mux
}
