package main

import (
	"log"
	"time"
)

var (
	nodesHosts      []string
	nodesMacaroons  []string
	nodesTLSCerts   []string
	pollingInterval time.Duration // Parseable by time.ParseDuration
	ErrorLog        *log.Logger
	InfoLog         *log.Logger
	DebugLog        *log.Logger
	WarnLog         *log.Logger
	//TODO LogLevel
)

type App struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DebugLog *log.Logger
	WarnLog  *log.Logger
}

func NewApp(errorLog *log.Logger, InfoLog *log.Logger, DebugLog *log.Logger, WarnLog *log.Logger) *App {
	return &App{
		ErrorLog: errorLog,
		InfoLog:  InfoLog,
		DebugLog: DebugLog,
		WarnLog:  WarnLog,
	}
}

