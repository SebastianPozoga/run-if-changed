package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/memfs"
)

func TestHasherWalk(t *testing.T) {
	var (
		forFS filesystem.Filespace

		err error

		expected = []string{"./dir1/file1", "./test1", "./test2"}
		result   = make([]string, 0)
	)
	if forFS, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	forFS.WriteFile("./test1", []byte{}, filesystem.SafeFilePermissions)
	forFS.WriteFile("./test2", []byte{}, filesystem.SafeFilePermissions)
	forFS.MkdirAll("./dir1", filesystem.SafeDirPermissions)
	forFS.WriteFile("./dir1/file1", []byte{}, filesystem.SafeFilePermissions)

	logs := NewMockLogs()
	hasher := NewHasher(logs, forFS)
	hasher.Walk(WalkParams{
		CB: func(path string, info os.FileInfo) (err error) {
			result = append(result, path)
			return
		},
		Paths: WalkPaths{
			Excludes: nil,
			Includes: Paths{"."},
			Exts:     nil,
		},
	})
	if reflect.DeepEqual(expected, result) {
		return
	}
	t.Errorf("Expected:\n %+v\n  and take:\n %+v", expected, result)
}

func TestHasherWalkExcludes(t *testing.T) {
	var (
		forFS filesystem.Filespace

		err error

		expected = []string{"./test1", "./test2"}
		result   = make([]string, 0)
	)
	if forFS, err = memfs.NewFilespace(); err != nil {
		t.Error(err)
		return
	}
	forFS.WriteFile("./test1", []byte{}, filesystem.SafeFilePermissions)
	forFS.WriteFile("./test2", []byte{}, filesystem.SafeFilePermissions)
	forFS.MkdirAll("./dir1", filesystem.SafeDirPermissions)
	forFS.WriteFile("./dir1/file1", []byte{}, filesystem.SafeFilePermissions)

	logs := NewMockLogs()
	hasher := NewHasher(logs, forFS)
	hasher.Walk(WalkParams{
		CB: func(path string, info os.FileInfo) (err error) {
			result = append(result, path)
			return
		},
		Paths: WalkPaths{
			Excludes: Paths{"./dir1"},
			Includes: Paths{"."},
			Exts:     nil,
		},
	})
	if reflect.DeepEqual(expected, result) {
		return
	}
	t.Errorf("Expected:\n %+v\n  and take:\n %+v", expected, result)
}
