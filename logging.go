package main

import (
	"log"
	"os"
	"path"
)

func initLoggers() {

	infoFile, _ := createLogFile("atmosreader_Info.log")
	errorFile, _ := createLogFile("atmosreader_Error.log")

	Info = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func createLogFile(fileName string) (*os.File, error) {
	f, err := os.OpenFile(path.Join(storagePath, fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	return f, nil
}
