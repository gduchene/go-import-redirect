// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !aws

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"golang.org/x/sys/unix"

	"go.awhk.org/gosdd"
)

var (
	addr = flag.String("addr", "localhost:8080", "address to listen on")
	from = flag.String("from", "", "package prefix to remove")
	to   = flag.String("to", "", "repository prefix to add")
	vcs  = flag.String("vcs", "git", "version control system to signal")
)

func main() {
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, unix.SIGTERM)

	srv := http.Server{Handler: &redirector{regexp.MustCompile(*from), *to, *vcs}}
	go func() {
		ln, err := listenSD()
		if err != nil {
			log.Fatalln("listenSD:", err)
		}
		if ln == nil {
			ln = listenFlag()
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

func listenFlag() net.Listener {
	if (*addr)[0] != '/' {
		ln, err := net.Listen("tcp", *addr)
		if err != nil {
			log.Fatalln("net.Listen:", err)
		}
		return ln
	}
	ln, err := net.Listen("unix", *addr)
	if err != nil {
		log.Fatalln("net.Listen:", err)
	}
	// We do not do any authorization anyway, so 0666 makes sense here.
	if err = os.Chmod(*addr, 0666); err != nil {
		log.Println("Failed to set permissions on UNIX socket:", err)
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
