package main

import "fmt"

type Log struct {
	Prefix string
	Active bool
}

func (log Log) Log(format string, args ...any) {
	if !log.Active {
		return
	}
	fmt.Printf(log.Prefix+format+"\n", args...)
}

type Logs struct {
	Dev     Log
	Warning Log
	Error   Log
	Log     Log
}

func NewLogs(dev bool) Logs {
	return Logs{
		Dev:     Log{"     Dev: ", dev},
		Warning: Log{" Warning: ", true},
		Error:   Log{"   Error: ", true},
		Log:     Log{"     Log: ", true},
	}
}

func NewMockLogs() Logs {
	return Logs{
		Dev:     Log{"", false},
		Warning: Log{"", false},
		Error:   Log{"", false},
		Log:     Log{"", false},
	}
}
