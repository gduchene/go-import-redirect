// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package internal

import (
	"fmt"
	"os"
	"strings"
)

func GetBody(pkg string) string {
	dest := GetDest(os.Getenv("PREFIX"), os.Getenv("DEST"), pkg)
	return fmt.Sprintf(`<!doctype html>
<meta name="go-import" content="%s %s %s">
<title>go-import-redirect</title>
`, pkg, os.Getenv("VCS"), dest)
}

func GetDest(srcPrefix, destPrefix, pkg string) string {
	srcPrefix = strings.TrimRight(srcPrefix, "/")
	destPrefix = strings.TrimRight(destPrefix, "/")
	path := strings.TrimLeft(strings.TrimPrefix(pkg, srcPrefix), "/")
	return destPrefix + "/" + strings.Split(path, "/")[0]
}
