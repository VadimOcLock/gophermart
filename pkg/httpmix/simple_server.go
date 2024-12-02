package httpmix

import (
	"fmt"
	"net/http"
	"time"
)

const (
	readTimeout  = 15 * time.Second
	writeTimeout = 15 * time.Second
)

func NewSimpleServer(handler http.Handler, port int) *http.Server {
	return &http.Server{ //nolint: exhaustruct
		Handler:      handler,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: readTimeout,
		ReadTimeout:  writeTimeout,
	}
}
