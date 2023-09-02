// client.go
package main

import (
	//"context"
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"

	//"strconv"
	"strings"

	//"net/http"
	"context"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
	//"golang.org/x/text/message"
	//"github.com/quic-go/quic-go/http3"
)

const (
	discoveryTimeout = 20 * time.Second
)

func StartClient() {

	clientConfig := &quic.Config{
		MaxIdleTimeout: time.Hour,
	}

	servers := discoverServers()
	if len(servers) == 0 {
		log.Println("no servers found")
		return
	}
	// For simplicity, let's just connect to the first discovered server.
	// In a real-world scenario, you might want to provide a choice to the user.
	serverAddr := servers[0]

	//set up a quic client
	log.Println("Attempting to connect to server: ", serverAddr)
	ctx := context.Background()
	sess, err := quic.DialAddr(ctx, serverAddr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}, clientConfig)

	if err != nil {
		log.Println("Error connecting to server: ", err)
	}
	log.Println("Successfully connected to server: ", serverAddr)
	log.Println("Attempting to open stream to server")
	stream, err := sess.OpenStreamSync(ctx)
	log.Println("stream Opened")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("stream opened to server")
	log.Println("Session state after opening stream", sess.ConnectionState())
	writeToServer(stream)
	go readFromServer(stream)

}

func discoverServers() []string {
	// Address to listen for broadcasts
	addr := &net.UDPAddr{
		IP:   net.ParseIP(broadcastAddress), // Listen on all available interfaces
		Port: 4243,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to set up UDP listener: %v", err)
	}
	defer conn.Close()

	log.Println("listening for server broadcasts", addr.String())

	var servers []string
	buf := make([]byte, 1024)

	// Set a timeout for reading from the connection
	conn.SetReadDeadline(time.Now().Add(20 * time.Second))

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			log.Printf("Error reading: %v", err)
			continue
		}

		msg := string(buf[:n])
		log.Println("received message: ", msg)
		if strings.HasPrefix(msg, "QUIC SERVER AVAILABLE: ") {
			serverAddr := strings.TrimPrefix(msg, "QUIC SERVER AVAILABLE: ")
			if !contains(servers, serverAddr) {
				servers = append(servers, serverAddr)
				log.Println("Discovered server:", serverAddr)
				log.Println(servers)
			}
		}
	}

	return servers
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func readFromServer(stream quic.Stream) {
	// read response from the server
	for {
		buf := make([]byte, 100)
		n, err := stream.Read(buf)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			return
		}
		fmt.Printf("Received from Server: %s\n", buf[:n])

	}
}

func writeToServer(stream quic.Stream) {
	reader := bufio.NewReader(os.Stdin)

	//automatically send an initial message upon connection
	im := "client connected"
	_, err := stream.Write([]byte(im))
	if err != nil {
		log.Printf("Error writing initial message: %v", err)
		return
	}
	log.Println("initial message sent to server")
	for {
		//send data to server
		fmt.Println("Enter Message: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("EOF Detected. Closing Stream")
				stream.Close()
				return
			}
			log.Println("error reading use input", err)
			continue
		}
		message = strings.TrimSpace(message)

		//if user types 'exit', break loop
		if message == "exit" {
			fmt.Println("Connection closed")
			stream.Close()
			return
		}

		//send message to server
		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message: ", err)
			break
		}
		log.Println("Data written to server")
	}

}

// fmt.Println("Do you want to close the server? (Yes/No)")
// var input string
// fmt.Scan(&input)
// if input == "yes" {
// 	stream.Close()
// 	log.Println("Stream has been closed.")
// }
