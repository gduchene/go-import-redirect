// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package lib

import (
	"fmt"
	"os"
	"strings"
)

func GetBody(pkg string) string {
	dest := strings.TrimRight(os.Getenv("DEST"), "/") + "/" + getRepo(pkg)
	return fmt.Sprintf(`<!doctype html>
<meta name="go-import" content="%s %s %s">
<title>go-import-redirect</title>
`, pkg, os.Getenv("VCS"), dest)
}

func getRepo(pkg string) string {
	prefix := strings.TrimRight(os.Getenv("PREFIX"), "/")
	path := strings.TrimLeft(strings.TrimPrefix(pkg, prefix), "/")
	return strings.Split(path, "/")[0]
}
