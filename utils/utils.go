package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/asad9711/log-responder/models"
)

// ReadReqPayload read and unmarshal req payload into object
func ReadReqPayload(r *http.Request) (reqPayload models.ReqPayload, err error) {
	var reqBody []byte
	reqBody, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error while reading payload from request - ", err.Error())
		return
	}

	// Unmarshal
	err = json.Unmarshal(reqBody, &reqPayload)
	if err != nil {
		log.Println("error while unmarshalling req payload - ", err.Error())
		return
	}
	return
}

// checkIfFileExists to check if file exists
func checkIfFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}
	return true
}

// ReadNLinesFromFile read n lines from the given file
func ReadNLinesFromFile(fileName string, n int) (fileLines string, err error) {
	lineCounter := 0
	// check if file exists
	if checkIfFileExists(fileName) {
		readFile, _ := os.Open(fileName)

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			if lineCounter > n {
				break
			}
			fileLines += fileScanner.Text()
			fileLines += "\n"
			// fileLines = append(fileLines, []byte(fileScanner.Text()))
			lineCounter++
		}

		readFile.Close()

		log.Println("read data -", fileLines)
	} else {
		err = errors.New("file doesn't exist")
	}
	return
}
