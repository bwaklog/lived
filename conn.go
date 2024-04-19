package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type command struct {
	opcode int
	key string
	val int
}

type LamportClock struct {
	counter int
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
	// 1 -> sending the latest clinet list
	// 		to a node in the pool
	// 2 -> broadcast command
	header int
	DataNode Node
}


// TODO: Initialize the node
func (n *Node) initNode() {
	// we need to have the client have an empty list
	n.clientList = []net.Conn{}
	n.data.loadData()
	LamportClock.counter = 0
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
		DataNode: Node {
			clientList: []net.Conn{},
		},
	}

	marshNodeData := marshalData(emptyPayload)
	err = conn.Write(marshNodeData)
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

	begPayload := TCPPayload {
		header: 0,
		DataNode: Node{
			clientList: n.clientList,
		},
	}
	marshData, err := json.Marshal(begPayload)

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

func (n *Node) BroadcastCommand (, connection *net.Conn)

func (n *Node) HandlePayload(payload TCPPayload, connection *net.Conn) {
	switch payload.header {
		case 0: // when we are dealing with a new connection
				// into the data pool
			// n.clientList = append(n, payload.DataNode.clientList)

			clinetPayload := TCPPayload {
				DataNode: Node {
					clientList: n.clientList,
					data: n.data,
					LamportClock: n.LamportClock,
				},
			}

			marshalData, err := json.Marshal(clinetPayload)
			if err != nil {
				log.Fatal(err)
			}

			n.writeMarshalData(marshalData, connection)

		case 1:
			// when we are pinging all other nodes in a server after recieving
			// a client list from initial node, we need to send them the updated
			// client list

			clientPayload := TCPPayload {
				DataNode: Node {
					clientList: n.clientList,
				},
			}

			marshalData, err := json.Marshal(clientPayload)
			if err != nil {
				log.Fatal(err)
			}

			n.writeMarshalData(marshalData, connection)

		case 2: // handle broadcast command
			fmt.Print("TODO: Handle broadcast command")

			for _, node := range n.clientList {
				if node != *connection {

				}
			}

	}
}

func (n * Node) HandleServerConnection (connection *net.Conn, buffer string, clientList []net.Conn) {
	if len(clientList) == 0 {
		n.writeData(2, connection)
		return
	}
	var dataRecieved []byte
	_, err := (*connection).Read(dataRecieved)
	if err != nil {
		log.Fatal(err)
	}

	var nodeRecieved Node
	err = json.Unmarshal(dataRecieved, &nodeRecieved)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: As the connection we read from already has a clinet list,
	// we only need to update the "current" node with the new list
	n.clientList = nodeRecieved.clientList

}

func (n *Node)writeMarshalData (bytesMarshalData []byte, connection *net.Conn) {
	_, err = (*connection).Write(bytesMarshalData)
	if err != nil {
		log.Fatal(err);
	}
}

func readBlock(n *net.Conn) (string, error) {
	// fireing a blocking operation
	var buff []byte = make([]byte, 512)
	bytesRead, err := (*n).Read(buff)
	if err != nil {
		return "", err
	}

	return string(buff[:bytesRead]), nil
}

var connectionChan = []net.Conn{}

func appendVal(message string) {
	for _, c := range connectionChan {
		if err := writeToConnection(&c, message); err != nil {
			log.Fatal(err)
		}
	}
}

func writeToConnection(n *net.Conn, message string) error {
	_, err := (*n).Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func marshalData (payload struct, connection *net.Conn) []byte {
	marshaledData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	byteMarshledData := []byte(marshaledData)

	return byteMarshledData
}

func TCPActiveListener() {
	tcp, err := net.Listen("tcp", ":6060")
	if err != nil {
		log.Fatal("Failed to start the tcp listener")
	}
	fmt.Println("Listening to port 6060")

	// create a channel of connections

	for {

		// creating threads, handling connections in separate
		// thread

		go func() {

			conn, err := tcp.Accept()
			if err != nil {
				log.Fatal(err)
			}
			connectionChan = append(connectionChan, conn)
			for _, c := range connectionChan {
				fmt.Println("connection: ", c)
			}
			for {
				// _ = append(connectionChan, conn)
				// append to the channel

				log.Println("appended to channel")

				connResp, err := readBlock(&conn)
				//connRespTrim := strings.Trim(connResp, "\n")

				if err == io.EOF {
					log.Printf("[%s] DISCONNECT {EOF}", conn.RemoteAddr())
				}

				if err != nil {
					conn.Close()
					// remove conn from the connectionChanl where
					// conn == conn

					// remove the connection from the channel
					//
					for i, c := range connectionChan {
						if c == conn {
							connectionChan = append(connectionChan[:i], connectionChan[i+1:]...)
						}
					}

					log.Printf("[%s] DISCONNECT", conn.RemoteAddr())
					break
				}

				// log.Printf("[%s]: %s", conn.LocalAddr().String(), connResp)

				// if err = writeToConnection(conn, connRespTrim); err != nil {
				// 	log.Fatal(err)
				// }
				// iterate througgh connnections and wtite to them

				snapData(connResp, &conn)

			}

		}()

		time.Sleep(5 * time.Second)

	}

}
