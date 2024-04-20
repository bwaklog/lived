package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type command struct {
	lamportCounter int
	opcode         int
	key            string
	val            int
}

type LamportClock struct {
	Counter int
	mutex   sync.Mutex
}

func (lc *LamportClock) incrementClock() {
	lc.mutex.Lock()
	lc.Counter++
	lc.mutex.Unlock()
}

type Node struct {
	// List of conncted clients
	clientList []net.Addr
	clientObjs []net.Conn

	// an instance of the lamport clock for
	// the node
	LamportClock LamportClock

	// Head command
	headCommand command

	// neck command
	neckCommand command

	// data availalbe to the node
	data store

	// store the tcp remote addres of the nodes system
	remoteAddr string
}

type TCPPayload struct {
	// 0 -> new connection into data pool
	// 1 -> sending the latest client list
	// 		to a node in the pool
	// 2 -> broadcast command
	Header int
	LamportClockCounter int
	HeadCommand         command
	NeckCommand         command
	Data                store
	ClientList          []net.Addr
}

// TODO: Initialize the node
func (n *Node) initNode() {
	// we need to have the client have an empty list
	n.clientList = []net.Addr{}
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

func (n *Node) EstabshConnection(connectionAddr string) {
	// develop a connection with the node

	conn, err := net.Dial("tcp", connectionAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to node: ", conn.RemoteAddr().String())

	// adding connection to the clinet list
	temp := conn.RemoteAddr()

	n.clientList = append(n.clientList, temp)

	emptyPayload := TCPPayload{
		Header:     0,
		ClientList: n.clientList,
	}

	marshalData, err := json.Marshal(emptyPayload)
	if err != nil {
		log.Fatal(err)
	}

	var buff TCPPayload
	err = json.Unmarshal(marshalData, &buff)

	fmt.Println(buff)

	marshalDataBytes := []byte(marshalData)

	fmt.Println("Sending an empty payload")

	_, err = conn.Write(marshalDataBytes)
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

	var nodeRecieved TCPPayload
	fmt.Println("Retrieving the data from node about connection list")

	err = json.Unmarshal(dataRecieved, &nodeRecieved)
	if err != nil {
		log.Fatal(err)
	}

	// copy unmarsheld data to local n data

	n.clientList = nodeRecieved.ClientList
	n.headCommand = nodeRecieved.HeadCommand
	n.neckCommand = nodeRecieved.NeckCommand
	n.data = nodeRecieved.Data

	// marshal the node data

	shareLatestClients := TCPPayload{
		Header:     1,
		ClientList: n.clientList,
	}

	marshData, err := json.Marshal(shareLatestClients)

	for _, node := range n.clientList {
		// establish connection with all nodes in client client list
		clinetCon, err := net.Dial("tcp", node.String())
		if err != nil {
			log.Fatal(err)
		}
		n.clientObjs = append(n.clientObjs, clinetCon)
		clinetCon.Write(marshData)
	}
}

// TODO: Broadcast commands

func (n *Node) BroadcastCommand() {
	// n.neckCommand = n.headCommand
	// n.headCommand = headCommand

	payload := TCPPayload{
		Header:              2,
		HeadCommand:         n.headCommand,
		NeckCommand:         n.neckCommand,
		LamportClockCounter: n.LamportClock.Counter,
	}
	bytesMarshalData := marshalData(payload)

	for _, conn := range n.clientObjs {
		if conn.RemoteAddr().String() == n.remoteAddr {
			n.writeMarshalData(bytesMarshalData, &conn)
		}
	}
}

func (n *Node) HandlePayload(payload TCPPayload, connection *net.Conn) {
	switch payload.Header {
	case 0: // when we are dealing with a new connection
		// into the data pool
		// n.clientList = append(n, payload.DataNode.clientList)

		tempClientList := payload.ClientList
		tempClientList = append(tempClientList, n.clientList...)

		clientPayload := TCPPayload{
			ClientList:          tempClientList,
			Data:                n.data,
			LamportClockCounter: n.LamportClock.Counter,
		}

		fmt.Println(clientPayload)

		marshalData, err := json.Marshal(clientPayload)
		if err != nil {
			log.Fatal(err)
		}

		marshalDataBytes := []byte(marshalData)

		n.writeMarshalData(marshalDataBytes, connection)

		n.clientList = append(n.clientList, (*connection).RemoteAddr())

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
	// case 3: // this deals with the data payload
	// that we use to confirm that the connection
	// is successfully established

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
func (n *Node) writeMarshalData(bytesMarshalData []byte, connection *net.Conn) {
	_, err := (*connection).Write(bytesMarshalData)
	if err != nil {
		log.Fatal(err)
	}
}

// DEBUG FUNCTION
func writeToConnection(n *net.Conn, message string) error {
	_, err := (*n).Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func marshalData(payload TCPPayload) []byte {
	marshaledData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	byteMarshledData := []byte(marshaledData)

	return byteMarshledData
}

func readConn(n *net.Conn) ([]byte, error) {
	var byte []byte = []byte{}
	_, err := (*n).Read(byte)
	if err != nil {
		return nil, err
	}
	return byte, nil
}

func startLocalListener(node *Node, addrFlag string, wg *sync.WaitGroup) {
	tcp, err := net.Listen("tcp", ":"+addrFlag)
	if err != nil {
		log.Fatal(err)
	}
	node.remoteAddr = tcp.Addr().String()
	defer tcp.Close()

	fmt.Println("Starting a TCP listener at ", tcp.Addr().String())

	for {
		go func() {
			// fmt.Println("listening for connec")
			conn, err := tcp.Accept()

			// if len of node.ClinetLIst is 0, then we finish the weight grp

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("[%s] CONNECTED\n", conn.RemoteAddr().String())

			for {
				connResp, err := readConn(&conn)
				if err == io.EOF {
					log.Printf("[%s] CLOSED - EOF", conn.RemoteAddr().String())
				}

				if err != nil {
					conn.Close()
					break
				}
				var payload TCPPayload
				// fmt.Println(string(connResp))
				json.Unmarshal(connResp, &payload)

				if len(node.clientList) == 0 {
					node.HandlePayload(payload, &conn)
					wg.Done()
					break
				}

				// log.Printf("%v", payload)
			}

		}()
		time.Sleep(1 * time.Second)
	}
}
