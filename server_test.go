package main

import "testing"

func TestSplitUrl(t *testing.T) {

	if parts := splitUrl("/xyz/abc/tgif/"); len(parts) != 3 {
		t.Errorf("should have returned 3 parts, was %d", len(parts))
	}
}

func TestGetBlob(t *testing.T) {

	if _, err := getBlob("this-hash-does-not-exist"); err.Error() != "No such blob" {
		t.Errorf("Error was '%s', should have been 'No such blob'", err.Error())
	}
}
