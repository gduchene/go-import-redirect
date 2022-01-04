// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !aws

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

	"go.awhk.org/gosdd"
)

func redirect(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.Header().Set("Allow", http.MethodGet)
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	if req.URL.Query().Get("go-get") != "1" {
		resp.Header().Set("Location", "https://pkg.go.dev/"+pkg)
		resp.WriteHeader(http.StatusFound)
		return
	}
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
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
		ln, err := listenSD()
		if err != nil {
			log.Fatalln("listenSD:", err)
		}
		if ln == nil {
			ln = listenEnv()
		}
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
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

func listenEnv() net.Listener {
	addr := os.Getenv("ADDR")
	if addr == "" || addr[0] != '/' {
		if addr == "" {
			addr = ":8080"
		}
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalln("net.Listen:", err)
		}
		return ln
	}
	ln, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatalln("net.Listen:", err)
	}
	// We do not do any authorization anyway, so 0666 makes sense here.
	if err = os.Chmod(addr, 0666); err != nil {
		log.Fatalln("os.Chmod:", err)
	}
	return ln
}

func listenSD() (net.Listener, error) {
	fds, err := gosdd.SDListenFDs(true)
	if err != nil {
		if err == gosdd.ErrNoSDSupport {
			return nil, nil
		}
		return nil, err
	}
	if len(fds) == 0 {
		return nil, nil
	}
	return net.FileListener(fds[0])
}
