package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

type command struct {
	lamportCounter int
	opcode int
	key string
	val int
}

type LamportClock struct {
	Counter int
	mutex sync.Mutex
}

type Node struct {
	// List of conncted clients
	clientList []net.Conn

	// an instance of the lamport clock for
	// the node
	LamportClock LamportClock

	// Head command
	headCommand command

	// neck command
	neckCommand command

	// data availalbe to the node
	data store
}

type TCPPayload struct {
	// 0 -> new connection into data pool
	// 1 -> sending the latest client list
	// 		to a node in the pool
	// 2 -> broadcast command
	header int
	// DataNode Node
	LamportClockCounter int
	HeadCommand command
	NeckCommand command
	Data store
	ClientList []net.Conn
}


// TODO: Initialize the node
func (n *Node) initNode() {
	// we need to have the client have an empty list
	n.clientList = []net.Conn{}
	// initilising commands to nil
	n.LamportClock.Counter = 0
	n.headCommand = command{lamportCounter: n.LamportClock.Counter}
	n.neckCommand = command{lamportCounter: n.LamportClock.Counter}
	n.data.loadData()
}

// TODO: establish connection with "SOME" node in network
// - ask for the clinet list of the node youre connecting to
// - ask for the data of the node youre connecting to, and use that data
// - Along with data, copy the head and neck values of the node youre connecting to
// - Connect to everyone in clinet list except the node youre connected to
// - If the note is already part of the pool, then when connecting establishing connections
//   with other clinets in the pool, we only update each clients list of clients

func (n *Node) EstablishConnection(connectionAddr string) {
	// develop a connection with the node

	conn, err := net.Dial("tcp", connectionAddr)
	if err != nil {
		log.Fatal(err)
	}

	// adding connection to the clinet list
	n.clientList = append(n.clientList, conn)


	emptyPayload := TCPPayload {
		header: 0,
		ClientList: n.clientList,
	}

	// marshNodeData := marshalData(emptyPayload)
	marshalData, err := json.Marshal(emptyPayload)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(marshalData)
	if err != nil {
		log.Fatal(err)
	}

	// we are giving an empty connection to the connected node
	// in the cluster

	// the connected node will evaluate this payload

	// have a read blocking call to read the data from node
	// the data recienved is a byte format of marshaled json
	// data for the strut Node
	var dataRecieved []byte
	_, err = conn.Read(dataRecieved)
	if err != nil {
		log.Fatal(err)
	}

	var nodeRecieved Node
	err = json.Unmarshal(dataRecieved, &nodeRecieved)
	if err != nil {
		log.Fatal(err)
	}

	// copy unmarsheld data to local n data

	n.clientList = nodeRecieved.clientList
	n.headCommand = nodeRecieved.headCommand
	n.neckCommand = nodeRecieved.neckCommand
	n.data = nodeRecieved.data

	// marshal the node data

	shareLatestClients := TCPPayload {
		header: 1,
		ClientList: n.clientList,
	}

	marshData, err := json.Marshal(shareLatestClients)

	for _, node := range n.clientList {
		// establish connection with all nodes in client client list
		clinetCon, err := net.Dial("tcp", node.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		clinetCon.Write(marshData)
	}
}


// TODO: Broadcast commands

func (n *Node) BroadcastCommand (headCommand command, neckCommand command, connection *net.Conn) {
	n.headCommand = headCommand
	n.neckCommand = neckCommand

	payload := TCPPayload {
		header: 2,
		HeadCommand: n.headCommand,
		NeckCommand: n.neckCommand,
		LamportClockCounter: n.LamportClock.Counter,
	}
	bytesMarshalData := marshalData(payload)

	for _ = range n.clientList {
		n.writeMarshalData(bytesMarshalData, connection)
	}
}

func (n *Node) HandlePayload(payload TCPPayload, connection *net.Conn) {
	switch payload.header {
		case 0: // when we are dealing with a new connection
				// into the data pool
				// n.clientList = append(n, payload.DataNode.clientList)

			clientPayload := TCPPayload {
				ClientList: n.clientList,
				Data: n.data,
				LamportClockCounter: n.LamportClock.Counter,
			}

			marshalData, err := json.Marshal(clientPayload)
			if err != nil {
				log.Fatal(err)
			}

			n.writeMarshalData(marshalData, connection)

		case 1:
			// the client has sent the updated clinet list over the conn
			// in the payload. We update the current node client list with
			// this sent over data
			n.clientList = payload.ClientList

		case 2: // handle broadcast command
			n.LamportClock.mutex.Lock()
			if payload.LamportClockCounter > n.LamportClock.Counter {
				n.LamportClock.Counter = payload.LamportClockCounter
			}
			n.LamportClock.Counter++
			n.LamportClock.mutex.Unlock()

			// if payload.HeadCommand


	}
}

func (lc *LamportClock) compareClocks(revievedClock int) {
	lc.mutex.Lock()
	if revievedClock > lc.Counter {
		lc.Counter = revievedClock
	}
	lc.Counter++
	lc.mutex.Unlock()
}


// quick function to write marshaled data to a conn
func (n *Node)writeMarshalData (bytesMarshalData []byte, connection *net.Conn) {
	_, err := (*connection).Write(bytesMarshalData)
	if err != nil {
		log.Fatal(err);
	}
}

// func readBlock(n *net.Conn) (string, error) {
// 	// fireing a blocking operation
// 	var buff []byte = make([]byte, 512)
// 	bytesRead, err := (*n).Read(buff)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(buff[:bytesRead]), nil
// }

// var connectionChan = []net.Conn{}


// DEBUG FUNCTION
func writeToConnection(n *net.Conn, message string) error {
	_, err := (*n).Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func marshalData (payload TCPPayload) []byte {
	marshaledData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	byteMarshledData := []byte(marshaledData)

	return byteMarshledData
}

// func TCActiveListener (connection *net.Conn) {
// 	tcp, err := net.Listen("tcp", ":6060")
// 	if (err != nil) {
// 		log.Fatal("TCP Error")
// 	}

// 	for {

// 		}
// 	}
// }

// func TCPActiveListener() {
// 	tcp, err := net.Listen("tcp", ":6060")
// 	if err != nil {
// 		log.Fatal("Failed to start the tcp listener")
// 	}
// 	fmt.Println("Listening to port 6060")

// 	// create a channel of connections

// 	for {

// 		// creating threads, handling connections in separate
// 		// thread

// 		go func() {

// 			conn, err := tcp.Accept()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			connectionChan = append(connectionChan, conn)
// 			for _, c := range connectionChan {
// 				fmt.Println("connection: ", c)
// 			}
// 			for {
// 				// _ = append(connectionChan, conn)
// 				// append to the channel

// 				log.Println("appended to channel")

// 				connResp, err := readBlock(&conn)
// 				//connRespTrim := strings.Trim(connResp, "\n")

// 				if err == io.EOF {
// 					log.Printf("[%s] DISCONNECT {EOF}", conn.RemoteAddr())
// 				}

// 				if err != nil {
// 					conn.Close()
// 					// remove conn from the connectionChanl where
// 					// conn == conn

// 					// remove the connection from the channel
// 					//
// 					for i, c := range connectionChan {
// 						if c == conn {
// 							connectionChan = append(connectionChan[:i], connectionChan[i+1:]...)
// 						}
// 					}

// 					log.Printf("[%s] DISCONNECT", conn.RemoteAddr())
// 					break
// 				}

// 				// log.Printf("[%s]: %s", conn.LocalAddr().String(), connResp)

// 				// if err = writeToConnection(conn, connRespTrim); err != nil {
// 				// 	log.Fatal(err)
// 				// }
// 				// iterate througgh connnections and wtite to them


// 				snapData(connResp, &conn)

// 			}

// 		}()

// 		time.Sleep(5 * time.Second)

// 	}

// }
