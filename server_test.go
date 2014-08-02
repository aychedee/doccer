package main

import "testing"
import "fmt"

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

func TestParseCommit(t *testing.T) {

	committer := "ashy david <a@b.com>"
	parent := "parenthash12345"
	content := "contenthash12345"
	cMap := parseCommit(
		fmt.Sprintf(
			"committer %s\nparent %s\ncontent %s\n", committer, parent, content))

	if cMap["committer"] != committer {
		t.Errorf("committer was '%s', should have been '%s'", cMap["committer"], committer)
	}

	if cMap["parent"] != parent {
		t.Errorf("parent was '%s', should have been '%s'", cMap["parent"], parent)
	}

	if cMap["content"] != content {
		t.Errorf("content was '%s', should have been '%s'", cMap["content"], content)
	}
}
