package main

import "testing"

func TestSplitUrl(t *testing.T) {

	if parts := splitUrl("/xyz/abc/tgif/"); len(parts) != 3 {
		t.Errorf("should have returned 3 parts, was %d", len(parts))
	}
}
