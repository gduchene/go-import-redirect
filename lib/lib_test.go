// Copyright (c) 2019, Grégoire Duchêne <gduchene@awhk.org>
//
// Use of this source code is governed by the ISC license that can be
// found in the LICENSE file.

package lib

import "testing"

func TestGetDest(t *testing.T) {
	cs := []struct{ srcPrefix, destPrefix, pkg, expected string }{
		{"src.example.com/x/", "https://example.com/git/", "src.example.com/x/foo", "https://example.com/git/foo"},
		{"src.example.com/x/", "https://example.com/git/", "src.example.com/x/foo/bar", "https://example.com/git/foo"},
		{"src.example.com/x", "https://example.com/git", "src.example.com/x/foo", "https://example.com/git/foo"},
		{"src.example.com/x", "https://example.com/git", "src.example.com/x/foo/bar", "https://example.com/git/foo"},
	}
	for _, c := range cs {
		actual := GetDest(c.srcPrefix, c.destPrefix, c.pkg)
		if actual != c.expected {
			t.Errorf("expected %s, got %s", c.expected, actual)
		}
	}
}
