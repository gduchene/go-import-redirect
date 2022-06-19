// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !aws

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"go.awhk.org/core"
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
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{Handler: &redirector{regexp.MustCompile(*from), *to, *vcs}}
	go func() {
		if err := srv.Serve(core.Must(core.Listen(*addr))); err != nil && err != http.ErrServerClosed {
			log.Fatalln("server.Serve:", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalln("server.Shutdown:", err)
	}
}
