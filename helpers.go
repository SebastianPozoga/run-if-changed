package main

import (
	"path/filepath"
	"strings"
)

func cleanPaths(paths []string) (cleaned []string) {
	for i := range paths {
		paths[i] = cleanPath(paths[i])
	}
	return paths
}

func cleanPath(path string) (cleanPath string) {
	path = filepath.Clean(path)
	return strings.TrimLeft(path, "./")
}
