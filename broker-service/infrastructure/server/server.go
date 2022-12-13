package server

import (
	"fmt"
	"helper"
	"net/http"
	"time"
)

// takes a http.hanlder interface as an argument(mux) and server that handler to the client
func Serve(mux http.Handler) (err error) {

	srv := &http.Server{
		Handler:           mux,
		Addr:              fmt.Sprintf(":%d", helper.GetEnvInt("PORT")),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	err = srv.ListenAndServe()

	return
}
