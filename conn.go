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
	var connectionChan = []net.Conn{}

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

					log.Printf("[%s] DISCONNECT", conn.RemoteAddr())
					break
				}

				log.Printf("[%s]: %s", conn.LocalAddr().String(), connResp)

				// if err = writeToConnection(conn, connRespTrim); err != nil {
				// 	log.Fatal(err)
				// }
				// iterate througgh connnections and wtite to them

				for _, c := range connectionChan {
					fmt.Println("Writing to conncetion: ", c)
					if err = writeToConnection(&c, connResp); err != nil {
						log.Fatal(err)
					}
				}
			}

		}()

		time.Sleep(5 * time.Second)

	}

}
