package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asad9711/log-responder/utils"
	"github.com/gorilla/websocket"
)

// endpoint handler to read n lines from a file and return data in response
func read(w http.ResponseWriter, req *http.Request) {

	reqPayload, err := utils.ReadReqPayload(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	log.Println("req payload - ", reqPayload)

	// upgrade incoming http conn to websocket
	// then, start reading from the provided file
	fileName := reqPayload.FileName
	numberOfLinesToRead, _ := strconv.Atoi(reqPayload.NumberOfLines)
	fileLines, err := utils.ReadNLinesFromFile(fileName, numberOfLinesToRead)
	if err != nil {
		log.Println("error while reading file ", fileName, err.Error())
	} else {
		log.Println("successfully read from file - ", fileLines)
	}

	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fileLines))

}

// endpoint handler to watch/follow n lines from a file and stream data in websocket connection
func watch(w http.ResponseWriter, req *http.Request) {

	reqPayload, err := utils.ReadReqPayload(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	log.Println("req payload - ", reqPayload)

	// upgrade incoming http conn to websocket
	// then, start reading from the provided file
	fileName := reqPayload.FileName
	numberOfLinesToRead, _ := strconv.Atoi(reqPayload.NumberOfLines)

	// upgrade to websocket
	var upgrader = websocket.Upgrader{}
	clientWSConnection, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		errMsg := fmt.Sprintf("could not upgrade to WS - %s", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	defer clientWSConnection.Close()
	log.Println("successfully upgraded to WS")

	// read from file continuously
	for {
		fileLines, err := utils.ReadNLinesFromFile(fileName, numberOfLinesToRead)
		if err != nil {
			log.Println("error while reading file ", fileName, err.Error())
		} else {
			log.Println("successfully read from file - ", fileLines)
		}

		// TODO: reading from last of file is tricky. Take ideas from the linux utility tail
	}
}

func StartServer() {
	http.HandleFunc("/read", read)
	http.HandleFunc("/watch", watch)
	log.Println("starting http server")

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Println("error while launching server - ", err.Error())
	} else {
		log.Println("launched server successfully at port 8090")
	}

}

func main() {
	StartServer()
}
