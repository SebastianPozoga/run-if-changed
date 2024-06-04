package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil"
)

type WalkParams struct {
	CB    func(path string, info os.FileInfo) (err error)
	Paths WalkPaths
}

type Hasher struct {
	cwd  filesystem.Filespace
	logs Logs
}

func NewHasher(logs Logs, cwd filesystem.Filespace) Hasher {
	return Hasher{
		cwd:  cwd,
		logs: logs,
	}
}

func (hasher *Hasher) Hash(paths WalkPaths) (result []byte, err error) {
	hash := sha256.New()
	if err = hasher.Walk(WalkParams{
		CB: func(path string, node os.FileInfo) (err error) {
			var (
				bytes []byte
			)
			hasher.logs.Dev.Log("hash %s", path)
			if _, err = hash.Write([]byte(path)); err != nil {
				return
			}
			if bytes, err = node.ModTime().MarshalBinary(); err != nil {
				return
			}
			if _, err = hash.Write(bytes); err != nil {
				return
			}
			bytes = make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, uint64(node.Size()))
			_, err = hash.Write(bytes)
			return
		},
		Paths: paths,
	}); err != nil {
		return
	}
	result = hash.Sum(nil)
	return
}

func (hasher *Hasher) Walk(params WalkParams) (err error) {
	sort.Strings(params.Paths.Includes)
	for _, node := range params.Paths.Includes {
		if !hasher.cwd.IsExist(node) {
			return fmt.Errorf("node '%s does not exist", node)
		}
		if err = hasher.walk(params, node); err != nil {
			return
		}
	}
	return
}

func (hasher *Hasher) walk(params WalkParams, includePath string) (err error) {
	var (
		infos    []fs.FileInfo
		basePath string
	)
	if varutil.IsArrContainStr(params.Paths.Excludes, includePath) {
		hasher.logs.Dev.Log("exclude %s", includePath)
		return
	}
	if hasher.cwd.IsFile(includePath) {
		var info os.FileInfo
		if info, err = hasher.cwd.Lstat(includePath); err != nil {
			return
		}
		return params.CB(includePath, info)
	}
	basePath = includePath + "/"
	if infos, err = hasher.cwd.ReadDir(includePath); err != nil {
		return
	}
	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})
	for _, node := range infos {
		if node.Name() == "." || node.Name() == ".." {
			continue
		}
		nodePath := basePath + node.Name()
		if varutil.IsArrContainStr(params.Paths.Excludes, nodePath) {
			continue
		}
		if node.IsDir() {
			hasher.logs.Dev.Log("dir %s", includePath)
			if err = hasher.walk(params, nodePath); err != nil {
				return
			}
			continue
		}
		if err = params.CB(nodePath, node); err != nil {
			return
		}
	}
	return
}
