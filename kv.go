package main

import (
	"log"
	"os"
)

type store struct {
	// store the data of the file as bytes
	content []byte
}

func (s *store) loadData() {
	datafile, err := os.Open("data/data.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = datafile.Read(s.content)
	if err != nil {
		log.Fatal(err)
	}

}

func (s *store) PersistData() {
	datafile, err := os.Create("data/data.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = datafile.Write(s.content)
}

/*
func snapData(operation string, n *net.Conn) {
	// TODO:
	// ex: LOG "text"

	operations := strings.Split(operation, " ")

	if operations[0] == "LOG" {
		// appendVal(operations[1])
		appendVal(operations[1])
	} else {
		// host name of connection
		hostname := (*n).RemoteAddr().String()
		appendVal(fmt.Sprintf("[%s]: Invalid operation %s", hostname, operations[0]))
	}
}
 */
