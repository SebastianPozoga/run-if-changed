package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/filesystem/filespace/diskfs"
)

const HashFilePath = "/hash"

type WalkPaths struct {
	Excludes Paths
	Includes Paths
	Exts     Paths
}

type AppServices struct {
	Logs   Logs
	Hasher Hasher
}

type App struct {
	Paths WalkPaths

	cwd      filesystem.Filespace
	cache    filesystem.Filespace
	services AppServices

	commandArgs []string
}

type AppParams struct {
	CWD      string
	CacheDir string

	Paths       WalkPaths
	CommandArgs []string

	Dev bool
}

func NewApp(params AppParams) (app *App, err error) {
	var cache, cwd filesystem.Filespace
	if params.CWD == "" {
		params.CWD = "."
	}
	if params.CacheDir == "" {
		params.CacheDir = "cache"
	}
	params.CacheDir = cleanPath(params.CacheDir)
	params.Paths.Excludes = append(params.Paths.Excludes, params.CacheDir)
	params.Paths.Excludes = cleanPaths(params.Paths.Excludes)
	params.Paths.Includes = cleanPaths(params.Paths.Includes)
	if cwd, err = diskfs.NewFilespace(params.CWD); err != nil {
		return
	}
	if err = cwd.MkdirAll(params.CacheDir, filesystem.DefaultUnixDirMode); err != nil {
		return
	}
	if cache, err = diskfs.NewFilespace(params.CacheDir); err != nil {
		return
	}
	logs := NewLogs(params.Dev)
	hasher := NewHasher(logs, cwd)
	app = &App{
		Paths: params.Paths,

		cwd:   cwd,
		cache: cache,
		services: AppServices{
			Logs:   logs,
			Hasher: hasher,
		},

		commandArgs: params.CommandArgs,
	}
	logs.Dev.Log("logs on")
	logs.Dev.Log("Includes: %v", strings.Join(app.Paths.Includes, ", "))
	logs.Dev.Log("Excludes: %v", strings.Join(app.Paths.Excludes, ", "))
	logs.Dev.Log("Exts: %v", strings.Join(app.Paths.Exts, ", "))
	return
}

func (app *App) Valid() {
	if len(app.Paths.Includes) == 0 {
		panic("include flag is required")
	}
}

func (app *App) Run() (err error) {
	var newHash, hash []byte
	if newHash, err = app.services.Hasher.Hash(app.Paths); err != nil {
		return
	}
	if app.cache.IsExist(HashFilePath) {
		if hash, err = app.cache.ReadFile(HashFilePath); err != nil {
			return
		}
		if bytes.Equal(newHash, hash) {
			app.services.Logs.Log.Log("no modified")
			return
		}
		app.cache.Remove(HashFilePath)
	}
	app.services.Logs.Log.Log("Run script: %s", strings.Join(app.commandArgs, " "))
	if err = app.exec(); err != nil {
		return
	}
	if err = app.cache.WriteFile(HashFilePath, newHash, filesystem.DefaultUnixFileMode); err != nil {
		return
	}
	app.services.Logs.Log.Log("hash updated")
	return
}

func (app *App) exec() (err error) {
	cmd := exec.Command(app.commandArgs[0], app.commandArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
