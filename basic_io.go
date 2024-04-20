package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func sendDataToNode(ip string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Prepare the URL to send data to
	url := fmt.Sprintf("http://%s:%d", ip, 8080)

	// Prepare the data to send
	data := "Hello from the client!"

	// Send a POST request to the node
	resp, err := http.Post(url, "text/plain", bytes.NewBufferString(data))
	if err != nil {
		log.Printf("Error sending data to %s:%d: %v\n", ip, port, err)
		return
	}
	defer resp.Body.Close()

	// Read the response from the node
	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		log.Printf("Error reading response from %s:%d: %v\n", ip, port, err)
		return
	}

	// Print the response
	log.Printf("Response from %s:%d: %s\n", ip, port, body.String())
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Do you want to be a server (yes/no)?")
	serverInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	serverInput = strings.TrimSpace(serverInput)

	if strings.ToLower(serverInput) == "yes" {
		// Start HTTP server
		fmt.Println("Enter the port number to listen for connections:")
		portInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err)
		}
		portInput = strings.TrimSpace(portInput)

		log.Printf("Starting server and listening for connections on port %s...\n", portInput)
		go func() {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hi from the server!")
			})

			log.Fatal(http.ListenAndServe(":"+portInput, nil))
		}()

		// Keep alive
		select {}
	} else if strings.ToLower(serverInput) == "no" {
		// List of IP addresses and port numbers of other nodes
		nodes := []struct {
			IP   string
			Port int
		}{
			{"tapank.local", 8080},
			{"cider.local", 8080},
			// Add more nodes if needed
		}

		// Use WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup

		// Iterate over nodes and send data concurrently
		for _, node := range nodes {
			wg.Add(1)
			go sendDataToNode(node.IP, node.Port, &wg)
		}

		// Wait for all goroutines to finish
		wg.Wait()
	} else {
		fmt.Println("Invalid input. Please enter 'yes' or 'no'.")
	}
}
