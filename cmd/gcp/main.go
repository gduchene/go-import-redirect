// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package main

import (
	"go.awhk.org/go-import-redirect/internal"
	"log"
	"net/http"
	"os"
	"path"
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", redirect)
	log.Fatal(http.ListenAndServe(os.Getenv("ADDR")+":"+port, nil))
}
