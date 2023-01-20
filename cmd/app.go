package cmd

import (
	"log"
)

var (
	nodesHosts      []string
	nodesMacaroons  []string
	pollingInterval int64 //in seconds
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

const name = "Moneytor"

func NewApp(errorLog *log.Logger, InfoLog *log.Logger, DebugLog *log.Logger, WarnLog *log.Logger) *App {
	return &App{
		ErrorLog: errorLog,
		InfoLog:  InfoLog,
		DebugLog: DebugLog,
		WarnLog:  WarnLog,
	}
}

func handleErr(err error, message string) {
	if err != nil {
		ErrorLog.Fatalf("%s: %v", message, err)
	}
}

func handleInfo(message string) {

	InfoLog.Printf("%s", message)

}

func handleDebug(message string) {

	DebugLog.Printf("%s", message)

}
