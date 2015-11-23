package main

import (
	"log"
	"os"
)

func initLoggers() {

	infoFile, _ := createLogFile("atmosReader_Info.log")
	errorFile, _ := createLogFile("atmosReader_Error.log")

	Info = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func createLogFile(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	return f, nil
}
