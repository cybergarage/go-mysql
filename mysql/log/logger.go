// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type LoggerOutpter func(file string, level LogLevel, msg string) (int, error)

type LogLevel int

type Logger struct {
	File     string
	Level    LogLevel
	outputer LoggerOutpter
}

const (
	logFormat   = "%s %s %s"
	logFilePerm = 0644
	LF          = "\n"
	LevelDebug  = (1 << 6)
	LevelTrace  = (1 << 5)
	LevelInfo   = (1 << 4)
	LevelWarn   = (1 << 3)
	LevelError  = (1 << 2)
	LevelFatal  = (1 << 1)
	LevelAll    = 0

	loggerLevelUnknownString = "UNKNOWN"
	loggerStdout             = "stdout"
)

var sharedLogger *Logger

// SetSharedLogger sets a singleton logger
func SetSharedLogger(logger *Logger) {
	sharedLogger = logger
}

// GetSharedLogger gets a shared singleton logger
func GetSharedLogger() *Logger {
	return sharedLogger
}

var logLevelStrings = map[LogLevel]string{
	LevelDebug: "DEBUG",
	LevelTrace: "TRACE",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

func getLogLevelString(logLevel LogLevel) string {
	logString, hasString := logLevelStrings[logLevel]
	if !hasString {
		return loggerLevelUnknownString
	}
	return logString
}

// SetLevel sets a output log level.
func (logger *Logger) SetLevel(level LogLevel) {
	logger.Level = level
}

// GetLevel gets the current log level.
func (logger *Logger) GetLevel() LogLevel {
	return logger.Level
}

// IsLevel returns true when the specified log level is enable, otherwise false.
func (logger *Logger) IsLevel(logLevel LogLevel) bool {
	if logLevel < logger.Level {
		return false
	}
	return true
}

// NewStdoutLogger creates a stdout logger.
func NewStdoutLogger(level LogLevel) *Logger {
	logger := &Logger{
		File:     loggerStdout,
		Level:    level,
		outputer: outputStdout}
	return logger
}

func outputStdout(file string, level LogLevel, msg string) (int, error) {
	fmt.Println(msg)
	return len(msg), nil
}

// NewFileLogger creates a file based logger.
func NewFileLogger(file string, level LogLevel) *Logger {
	logger := &Logger{
		File:     file,
		Level:    level,
		outputer: outputToFile}
	return logger
}

func outputToFile(file string, level LogLevel, msg string) (int, error) {
	msgBytes := []byte(msg + LF)
	fd, err := os.OpenFile(file, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), logFilePerm)
	if err != nil {
		return 0, err
	}

	writer := bufio.NewWriter(fd)
	writer.Write(msgBytes)
	writer.Flush()

	fd.Close()

	return len(msgBytes), nil
}

func output(outputLevel LogLevel, msgFormat string, msgArgs ...interface{}) int {
	if sharedLogger == nil {
		return 0
	}

	logLevel := sharedLogger.GetLevel()
	if logLevel < outputLevel {
		return 0
	}

	t := time.Now()
	logDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	headerString := fmt.Sprintf("[%s]", getLogLevelString(outputLevel))
	logMsg := fmt.Sprintf(logFormat, logDate, headerString, fmt.Sprintf(msgFormat, msgArgs...))
	logMsgLen := len(logMsg)

	if 0 < logMsgLen {
		logMsgLen, _ = sharedLogger.outputer(sharedLogger.File, logLevel, logMsg)
	}

	return logMsgLen
}

// Debug outputs a debug level message to loggers.
func Debug(format string, args ...interface{}) int {
	return output(LevelDebug, format, args...)
}

// Trace outputs trace level message to loggers.
func Trace(format string, args ...interface{}) int {
	return output(LevelTrace, format, args...)
}

// Info outputs a infomation level message to loggers.
func Info(format string, args ...interface{}) int {
	return output(LevelInfo, format, args...)
}

// Warn outputs a warning level message to loggers.
func Warn(format string, args ...interface{}) int {
	return output(LevelWarn, format, args...)
}

// Error outputs a error level message to loggers.
func Error(format string, args ...interface{}) int {
	return output(LevelError, format, args...)
}

// Fatal outputs a fatal level message to loggers.
func Fatal(format string, args ...interface{}) int {
	return output(LevelFatal, format, args...)
}

// Output outputs the specified level message to loggers.
func Output(outputLevel LogLevel, format string, args ...interface{}) int {
	return output(outputLevel, format, args...)
}
