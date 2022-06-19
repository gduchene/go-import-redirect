// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !aws

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"go.awhk.org/core"
	"go.awhk.org/go-import-redirect/pkg/redirector"
)

var (
	addr = flag.String("addr", "localhost:8080", "address to listen on")
	cfg  = flag.String("c", "", "path to a configuration file")
	from = flag.String("from", "", "package prefix to remove")
	to   = flag.String("to", "", "repository prefix to add")
	vcs  = flag.String("vcs", "git", "version control system to signal")
)

func main() {
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	h := &redirector.Redirector{Template: redirector.DefaultTemplate}
	if *cfg != "" {
		if err := json.NewDecoder(core.Must(os.Open(*cfg))).Decode(&h.Transformers); err != nil {
			log.Fatalln(err)
		}
	} else {
		h.Transformers = append(h.Transformers, redirector.Transformer{
			Pattern:     &redirector.Pattern{regexp.MustCompile(strings.ReplaceAll(*from, `\\`, `\`))},
			Replacement: *to,
			VCS:         *vcs,
		})
	}

	srv := http.Server{Handler: h}
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
