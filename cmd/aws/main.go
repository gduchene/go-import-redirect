// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

// +build linux

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.awhk.org/go-import-redirect/internal"
	"net/http"
	"path"
)

func redirect(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		pkg  = path.Join(req.Headers["Host"], req.Path)
		resp = events.APIGatewayProxyResponse{
			Body:    internal.GetBody(pkg),
			Headers: map[string]string{"Content-Type": "text/html; charset=utf-8"},
		}
	)
	if v, ok := req.QueryStringParameters["go-get"]; ok && v == "1" {
		resp.StatusCode = http.StatusOK
	} else {
		resp.Headers["Location"] = "https://godoc.org/" + pkg
		resp.StatusCode = http.StatusFound
	}
	return resp, nil
}

func main() {
	lambda.Start(redirect)
}