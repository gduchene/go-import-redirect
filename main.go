// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

// +build !aws

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"golang.org/x/sys/unix"
)

func redirect(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.Header().Set("Allow", http.MethodGet)
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	if req.URL.Query().Get("go-get") == "1" {
		resp.WriteHeader(http.StatusOK)
	} else {
		resp.Header().Set("Location", "https://pkg.go.dev/"+pkg)
		resp.WriteHeader(http.StatusFound)
	}
	if _, err := fmt.Fprint(resp, GetBody(pkg)); err != nil {
		log.Println("fmt.Fprint:", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirect)
	srv := http.Server{Handler: mux}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, unix.SIGTERM)

	go func() {
		var (
			addr = os.Getenv("ADDR")
			l    net.Listener
			err  error
		)
		if addr != "" && addr[0] == '/' {
			if l, err = net.Listen("unix", addr); err != nil {
				log.Fatalln("net.Listen:", err)
			}
			// We do not do any authorization anyway, so 0666 makes sense here.
			if err = os.Chmod(addr, 0666); err != nil {
				log.Fatalln("os.Chmod:", err)
			}
		} else {
			if addr == "" {
				addr = ":8080"
			}
			if l, err = net.Listen("tcp", addr); err != nil {
				log.Fatalln("net.Listen:", err)
			}
		}
		if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
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
