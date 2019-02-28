package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type (
	iLogger interface {
		FATAL(logString ...interface{})
		ERROR(logString ...interface{})
		WARN(logString ...interface{})
		INFO(logString ...interface{})
		DEBUG(logString ...interface{})
		TRACE(logString ...interface{})
	}
	logOutput interface {
		writeLog(level string, callingFile string, lineNumber int, logStr ...interface{})
	}

	logger struct {
		FatalLevel []logOutput
		ErrorLevel []logOutput
		InfoLevel  []logOutput
		WarnLevel  []logOutput
		DebugLevel []logOutput
		TraceLevel []logOutput
	}
	consoleOutput struct{}
)

const callingFunctionLevel = 2

var Logger iLogger

func initLogger() {
	Logger = logger{

		FatalLevel: []logOutput{consoleOutput{}},
		ErrorLevel: []logOutput{consoleOutput{}},
		WarnLevel:  []logOutput{consoleOutput{}},
		InfoLevel:  []logOutput{consoleOutput{}},
		DebugLevel: []logOutput{consoleOutput{}},
		TraceLevel: []logOutput{consoleOutput{}},
	}
}

func (l logger) FATAL(logString ...interface{}) { writeLog(l.FatalLevel, "FATAL", logString) }
func (l logger) ERROR(logString ...interface{}) { writeLog(l.ErrorLevel, "ERROR", logString) }
func (l logger) WARN(logString ...interface{})  { writeLog(l.WarnLevel, "WARN", logString) }
func (l logger) INFO(logString ...interface{})  { writeLog(l.InfoLevel, "INFO", logString) }
func (l logger) DEBUG(logString ...interface{}) { writeLog(l.DebugLevel, "DEBUG", logString) }
func (l logger) TRACE(logString ...interface{}) { writeLog(l.TraceLevel, "TRACE", logString) }

func getCompilePackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		compileDir := strings.Replace(file, "logger.go", "", 1)
		return compileDir
	} else {
		return ""
	}
}

func (c consoleOutput) writeLog(level string, callingFile string, lineNumber int, logStr ...interface{}) {
	var color, logString string
	timeStr := time.Now().Format("15:04:05")
	switch level {
	case "ERROR":
		color = "\033[21;31m"
	case "FATAL":
		color = "\033[21;37;41m"
	case "INFO":
		color = "\033[21;32m"
	case "WARN":
		color = "\033[21;35m"
	case "DEBUG":
		color = "\033[21;34m"
	}

	logString = fmt.Sprintf("%v%s %s %s:%d - %s \033[0m", color, timeStr, level, callingFile, lineNumber, logStr)
	fmt.Println(logString)
}

func writeLog(logOutputs []logOutput, level string, log ...interface{}) {
	projectRootDir := getCompilePackageDir()
	_, file, lineNumber, ok := runtime.Caller(callingFunctionLevel)
	if ok {
		callingFile := strings.Replace(filepath.ToSlash(file), filepath.ToSlash(projectRootDir), "", -1)
		var logMessage interface{}
		if len(log) == 1 { //single value log as first element
			logMessage = log[0]
		} else { //single value log as array
			logMessage = log
		}
		for _, output := range logOutputs {

			output.writeLog(level, callingFile, lineNumber, fmt.Sprintf("%s", logMessage))

		}
	}
}
