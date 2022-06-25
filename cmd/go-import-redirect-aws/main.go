// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package main

import (
	"context"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.awhk.org/go-import-redirect/pkg/redirector"
)

var transf = redirector.Transformer{
	Pattern:     &redirector.Pattern{regexp.MustCompile(strings.ReplaceAll(os.Getenv("FROM"), `\\`, `\`))},
	Replacement: os.Getenv("TO"),
	VCS:         os.Getenv("VCS"),
}

func redirect(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pkg := path.Join(req.Headers["Host"], req.Path)

	if !transf.Pattern.MatchString(pkg) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	if v, ok := req.QueryStringParameters["go-get"]; !ok || v != "1" {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Location": "https://pkg.go.dev/" + pkg},
			StatusCode: http.StatusFound,
		}, nil
	}

	data := redirector.TemplateData{
		Package:    pkg,
		Repository: transf.Pattern.ReplaceAllString(pkg, transf.Replacement),
		VCS:        transf.VCS,
	}
	var buf strings.Builder
	if err := redirector.DefaultTemplate.Execute(&buf, data); err != nil {
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
