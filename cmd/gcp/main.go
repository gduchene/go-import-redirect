// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"go.awhk.org/go-import-redirect/internal"
)

func redirect(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.Header()["Allow"] = []string{"GET"}
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	resp.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	if v, ok := req.URL.Query()["go-get"]; ok && len(v) > 0 && v[0] == "1" {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.Header()["Location"] = []string{"https://godoc.org/" + pkg}
		resp.WriteHeader(http.StatusFound)
	}
	resp.Write([]byte(internal.GetBody(pkg)))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirect)
	srv := http.Server{Handler: mux}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		srv.Addr = os.Getenv("ADDR") + ":" + port
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server.ListenAndServe:", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalln("server.Shutdown:", err)
	}
}
