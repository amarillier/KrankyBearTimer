package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

// see error logging:
// https://rollbar.com/blog/golang-error-logging-guide/
func logInit() {
	// typically written to Resources/...
	file, err := os.OpenFile("TaniumTimer0.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func lineCounter(r io.Reader) (int, error) {
	// count lines in a file, used for log rotation and possible other uses
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func logRotate() {
	if checkFileExists("TaniumTimer2.txt") {
		f := os.Remove("TaniumTimer2.txt")
		if f != nil {
			ErrorLog.Println("Error attempting to remove TaniumTimer2.txt")
		}
	}
	if checkFileExists("TaniumTimer1.txt") {
		os.Rename("TaniumTimer1.txt", "TaniumTimer2.txt")
	}
	if checkFileExists("TaniumTimer0.txt") {
		os.Rename("TaniumTimer0.txt", "TaniumTimer1.txt")
	}
}

func listMatchingFiles(directory, pattern string) ([]string, error) {
	var matchingFiles []string

	// Read the directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Loop through the files and match the pattern
	for _, file := range files {
		if matched, err := filepath.Match(pattern, file.Name()); err != nil {
			return nil, err
		} else if matched {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}
	return matchingFiles, nil
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
