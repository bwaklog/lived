package main

import (
	"fmt"

	"sync"
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

	var serverFlag bool

	var server string
	fmt.Printf("Server(y/n): ")
	fmt.Scanln(&server)
	switch server {
	case "y":
		serverFlag = true
	case "n":
		serverFlag = false

	}

	var addrPort string
	fmt.Print("Enter port to expose: ")
	fmt.Scanln(&addrPort)

	wg := sync.WaitGroup{}

	go startLocalListener(&node, addrPort, &wg)
	time.Sleep(1 * time.Second)

	if serverFlag == true {
		wg.Add(1)
		wg.Wait()
		fmt.Println("Established connection to a node")
	}
	if serverFlag == false {
		var remoteAddr string
		fmt.Print("Enter remote address: ")
		fmt.Scanln(&remoteAddr)

		node.EstabshConnection(remoteAddr)

	}

	for {
		// fmt.Print("Command: ")
		// var usercmd string
		// _, err := fmt.Scanln(&usercmd)
		// if err != nil {
		// 	log.Println(err)
		// }

		// if err == io.EOF {
		// 	break
		// }

		// operation := strings.Split(usercmd, " ")

		// // increment lampot clock before performing operation
		// node.LamportClock.incrementClock()

		// node.neckCommand = node.headCommand

		// tempHeadCmd := command{}

		// switch operation[0] {
		// case "SET":
		// 	tempHeadCmd.opcode = 1
		// case "GET":
		// 	tempHeadCmd.opcode = 2
		// }

		// tempHeadCmd.key = operation[1]
		// tempHeadCmd.val, err = strconv.Atoi(operation[2])
		// if err != nil {
		// 	log.Println(err)
		// 	// dealing with error for now with replacement with 0
		// 	// needs to be handeled differently
		// 	tempHeadCmd.val = 0
		// }
		// tempHeadCmd.lamportCounter = node.LamportClock.Counter

		// node.headCommand = tempHeadCmd

		// node.BroadcastCommand()
	}

	// time.Sleep(1 * time.Second)

	// serverFlag := flag.Bool("s", true, "sample help message")
	// remoteAddr := flag.String("r", "ERR", "sample help message")

	// if !*serverFlag {
	// 	// go startServer(&node)
	// 	flag.Parse()

	// 	fmt.Printf("Remote Address: %s\n", *remoteAddr)

	// 	node.EstablishConnection(*remoteAddr)
	// }

	/*
		for {
			// break from for loop if input is EOF

			fmt.Print("Command: ")
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
	*/

}
