package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		paths   WalkPaths
		cwd     string
		cache   string
		baseArg []string
		exArgs  []string

		app *App
		err error
		ok  bool

		dev bool
	)
	if baseArg, exArgs, ok = findExArg(os.Args); !ok {
		err = fmt.Errorf("external command is required. Add it after '--'. Like 'run-if-changed --include mydir -- command-to-run --commandArg1 v1 --commandArg2 v2' ")
		panic(err)
	}
	flag.Var(&paths.Excludes, "exclude", "directory or file to exclude")
	flag.Var(&paths.Includes, "include", "directory or file to include")
	flag.Var(&paths.Exts, "ext", "file extension to include")
	flag.StringVar(&cwd, "cwd", "", "current working directory")
	flag.StringVar(&cache, "cache", "", "cache directory")
	flag.BoolVar(&dev, "dev", false, "show dev logs")
	if err = flag.CommandLine.Parse(baseArg[1:]); err != nil {
		panic(err)
	}
	if app, err = NewApp(AppParams{
		CWD:         cwd,
		CacheDir:    cache,
		CommandArgs: exArgs,
		Paths:       paths,
		Dev:         dev,
	}); err != nil {
		panic(err)
	}
	app.Valid()
	if err = app.Run(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func findExArg(args []string) (base []string, exts []string, ok bool) {
	for index, arg := range args {
		if arg == "--" {
			return args[:index], args[index+1:], true
		}
	}
	return nil, nil, false
}
