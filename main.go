// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"path"
	"strings"
)

func getDest(dest, repo string) string {
	dest = strings.TrimRight(dest, "/")
	return fmt.Sprintf("%s/%s", dest, repo)
}

func getRepo(pkg, prefix string) string {
	prefix = strings.TrimRight(prefix, "/")
	path := strings.TrimLeft(strings.TrimPrefix(pkg, prefix), "/")
	return strings.Split(path, "/")[0]
}

func redirect(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		pkg  = path.Join(req.Headers["Host"], req.Path)
		dest = getDest(os.Getenv("DEST"), getRepo(pkg, os.Getenv("PREFIX")))
		doc  = fmt.Sprintf("https://godoc.org/%s", pkg)
		vcs  = os.Getenv("VCS")
		body = fmt.Sprintf(`<!doctype html>
<meta name="go-import" content="%s %s %s">
<title>go-import-redirect</title>
`, pkg, vcs, dest)
		resp = events.APIGatewayProxyResponse{
			Body:    body,
			Headers: map[string]string{"Content-Type": "text/html; charset=utf-8"},
		}
	)
	if v, ok := req.QueryStringParameters["go-get"]; ok && v == "1" {
		resp.StatusCode = 200
	} else {
		resp.Headers["Location"] = doc
		resp.StatusCode = 302
	}
	return resp, nil
}

func main() {
	lambda.Start(redirect)
}
