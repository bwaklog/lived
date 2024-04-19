package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

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
