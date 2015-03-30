/*
 * filenames contains utility functions for dealing with
 * filenames used in the trigger_loadtest program.
 */
package main

import (
	"bytes"
	"math/rand"
	"strings"
)

const (
	random_fname_size = 32
	alphabet          = "abcdefghijklmnopqrstuvwxyz"
)

// Check if a string represents an absolute path
func IsAbsolutePath(dirname string) bool {
	return strings.HasPrefix(dirname, "/")
}

// canonicalizeWorkingDir ensures that a directory name has
// exactly one trailing slash.
func CanonicalizeDirName(wd string) string {
	if strings.HasSuffix(wd, "/") == false {
		wd += "/"
	}
	return wd
}

// Returns a random string of length random_fname_size
func RandomFileName(directory string, rng *rand.Rand) string {
	var pos int32
	var result bytes.Buffer
	result.WriteString(directory)
	for i := 0; i < random_fname_size; i++ {
		pos = rng.Int31n(26)
		result.WriteByte(alphabet[pos])
	}
	result.WriteString(".MAT")
	return result.String()
}
