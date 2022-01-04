// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build aws && linux

package main

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	from  = os.Getenv("FROM")
	to    = os.Getenv("TO")
	vcs   = os.Getenv("VCS")
	redir = &redirector{from, to, vcs}
)

func redirect(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pkg := path.Join(req.Headers["Host"], req.Path)
	if v, ok := req.QueryStringParameters["go-get"]; !ok || v != "1" {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Location": "https://pkg.go.dev/" + pkg},
			StatusCode: http.StatusFound,
		}, nil
	}
	var buf strings.Builder
	if err := body.Execute(&buf, bodyData{pkg, redir.getRepo(pkg), vcs}); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       buf.String(),
		Headers:    map[string]string{"Content-Type": "text/html; charset=utf-8"},
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(redirect)
}
