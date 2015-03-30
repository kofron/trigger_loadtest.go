package main

import (
	"strings"
	"testing"
)

// Test that a trailing slash is added if none is provided.
func TestDirCanonicalizationNoTS(t *testing.T) {
	test_str := "/abc/def"
	test_str = CanonicalizeDirName(test_str)
	if strings.HasSuffix(test_str, "/") == false {
		t.Fail()
	}
}

// Test that if there already is a trailing slash, that
// canonicalization doesn't change the name at all
func TestDirCanonicalizationHasTS(t *testing.T) {
	test_str := "/abc/def/"
	test_str_2 := CanonicalizeDirName(test_str)
	if test_str != test_str_2 {
		t.Fail()
	}
}

// Test that absolute paths are detected as absolute
func TestAbsPathIsAbs(t *testing.T) {
	test_str := "/abc/def"
	if IsAbsolutePath(test_str) != true {
		t.Fail()
	}
}

// Test that relative paths are detected as relative
func TestRelPathIsNotAbs(t *testing.T) {
	test_str := "abc/def"
	if IsAbsolutePath(test_str) == true {
		t.Fail()
	}
}
