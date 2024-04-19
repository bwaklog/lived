package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {

	// testStore := store{}
	// testStore.loadData()

	// creating a node instance

	// node := Node{
	// 	clientList: []net.Conn{},
	// }

	node := Node{}
	node.initNode()

	// TCPActiveListener()

	// TODO: Having the tcp listener on one thread
	// and dealing with the clients on another. Hence,
	// writing to a clinet in the client list will

	go startLocalListener(&node)

	time.Sleep(10 * time.Second)

	remoteAddr := flag.String("r", "ERR", "sample help message")
	flag.Parse()

	fmt.Printf("Remote Address: %s\n", *remoteAddr)

	node.EstablishConnection(*remoteAddr)

	for {
		// break from for loop if input is EOF

		fmt.Println("Command: ")
		var usercmd string
		_, err := fmt.Scanln(&usercmd)
		if err != nil {
			log.Println(err)
		}

		if err == io.EOF {
			break
		}

		operation := strings.Split(usercmd, " ")

		// increment lampot clock before performing operation
		node.LamportClock.incrementClock()

		node.neckCommand = node.headCommand

		tempHeadCmd := command{}

		switch operation[0] {
			case "SET":
				tempHeadCmd.opcode = 1
			case "GET":
				tempHeadCmd.opcode = 2
		}

		tempHeadCmd.key = operation[1]
		tempHeadCmd.val, err = strconv.Atoi(operation[2])
		if err != nil {
			log.Println(err)
			// dealing with error for now with replacement with 0
			// needs to be handeled differently
			tempHeadCmd.val = 0
		}
		tempHeadCmd.lamportCounter = node.LamportClock.Counter


		node.headCommand = tempHeadCmd

		node.BroadcastCommand()

	}

}
