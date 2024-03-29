package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// WLog writes log to info.log file and sends to user
func WLog(msg string, ClientID uint) {
	if Conf.Debug {
		log.Println(msg)
	} else if strings.HasPrefix(msg, "Error:") {
		log.Println(msg)
	}
	UpdateUserMessage(msg, ClientID)
}

// Prints everything if Conf.Debug == true
func DebugLog(msg string) {
	if Conf.Debug {
		log.Println(msg)
	}
}

func PrintStruct(i interface{}, structName string) {
	if Conf.ShowStructs {
		j, err := json.Marshal(i)
		if err != nil {
			log.Println("PrintStruct")
			log.Panicln(err)
		}
		DebugLog(structName + ":")
		DebugLog(string(j))
	}
}

// OpenLogFile opens log file and/or creates it
func OpenLogFile(filepath string) error {
	var (
		err     error
		logFile *os.File
	)

	// Create log file if not exist
	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		_, err = os.Create(filepath)
		if err != nil {
			log.Fatalf("error creating file: %v", err)
			return err
		}
	}

	// Open log file
	logFile, err = os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return err
	}

	// Set log writer to log file insted of std
	log.SetOutput(logFile)

	return nil
}
