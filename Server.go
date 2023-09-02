package main

import (
	//"context"
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"net"
	//"net/http"

	"github.com/quic-go/quic-go"
	//"github.com/quic-go/quic-go/x509"
	//"github.com/quic-go/quic-go/http3"
)

//type tlsConfig struct {
//key string

// }

type ServerState int

const (
	broadcastAddress              = "127.0.0.1"
	broadcastPort                 = #
	broadcastInterval             = 1 * time.Second
	serverIdentifier              = "QUIC SERVER AVAILABLE:"
	Idle              ServerState = iota
	Broadcasting
	Connected
)

var state ServerState = Idle
var stateMutex sync.Mutex

func StartServer() {
	serverConfig := &quic.Config{
		MaxIdleTimeout: time.Hour,
	}

	// Setup the server with an address and certificate
	listener, err := quic.ListenAddr("127.0.0.1:#", generateTlsConfig(), serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	stopBroadcast := make(chan bool)
	go broadcastPresencePeriodically(stopBroadcast) // Start broadcasting immediately

	var sess quic.Connection
	connectionChannel := make(chan quic.Connection)

	for {
		switch state {
		case Idle:
			// Start broadcasting
			log.Println("Server state : Idle. Transitioning to broadcasting.")
			stateMutex.Lock()
			state = Broadcasting
			stateMutex.Unlock()

		case Broadcasting:
			log.Println("Entering broadcasting state")
			log.Println("Server state : Broadcast. Waiting for client connection.")
			// Accept connection to the server
			stopAccepting := make(chan bool)
			go func() {
				for {
					select {
					case <-stopAccepting:
						log.Println("received stopAccepting signal. Exiting accepting goroutine.")
						return
					default:
						log.Println("Attempting to accept a new client connection.")
						localSess, err := listener.Accept(context.Background())
						if err != nil {
							log.Println("Couldn't connect to client", err)
							continue
						}
						log.Println("accepted a client connection. Sending to connectionChannel")
						connectionChannel <- localSess
					}
				}
			}()
			//wait for either a connection or timeout
			select {
			case sess = <-connectionChannel: // This will assign the value to the sess variable declared at the function level
				log.Println("Received a client session from connectionChannel.")
				log.Println("Client connected: ", sess.RemoteAddr())
				log.Println("session state after connection: ", sess.ConnectionState())
				stateMutex.Lock()
				state = Connected
				stateMutex.Unlock()
				log.Println("Transitioning to Connected state. ")
				stopBroadcast <- true
				stopAccepting <- true
			case <-time.After(60 * time.Second):
				log.Println("No client connected in the last 20 seconds Continuing to broadcast")
				// Do not transition to Connected state here. Just continue the loop.
				continue
			}

			// Stop broadcasting and transition to connected state
			log.Println("Transitioning to connected state.")
			//stopBroadcast <- true
			//state = Connected

		case Connected:
			log.Println("Entered Connected state.")
			log.Println("Server state: Client Connected. Waiting for stream from client")
			log.Printf("Server session state before accepting stream: %v:", sess.ConnectionState())
			// Handle the connected client
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				log.Println("error accepting stream", err)
				log.Println("Transitioning back to broadcasting state.")
				stateMutex.Lock()
				state = Broadcasting
				stateMutex.Unlock()
				return

			}
			log.Println("Stream accepted from client.")
			log.Println("Session state after accepting stream:", sess.ConnectionState())
			handleClient(stream) // Handle client in a separate goroutine

			if stream.Context().Err() == context.Canceled {
				log.Println("Client has closed connection. Transitioning to broadcasting")
				stateMutex.Lock()
				state = Broadcasting
				stateMutex.Unlock()
				go broadcastPresencePeriodically(stopBroadcast)
			}

			// Once done, transition back to Broadcasting
			log.Println("Transitioning back to broadcasting state after handling client.")
			stateMutex.Lock()
			state = Broadcasting
			stateMutex.Unlock()
		}
	}
}

func handleClient(stream quic.Stream) {
	go readFromClient(stream)
	writeToClient(stream)

}

func readFromClient(stream quic.Stream) {
	for {
		buf := make([]byte, 100)
		n, err := stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client has closed the stream.")
			} else {
				log.Printf("Error reading from client: %v", err)
			}
			return
		}
		log.Printf("Received from Client: %s\n", buf[:n])
	}
}

func writeToClient(stream quic.Stream) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter Message to Client: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("EOF Detected. Closing Stream")
				stream.Close()
				return
			}
			log.Println("error reading server input", err)
			continue
		}
		message = strings.TrimSpace(message)

		// If the server types 'exit', break the loop
		if message == "exit" {
			fmt.Println("Connection closed by server command.")
			stream.Close()
			return
		}

		// Send message to client
		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message to client: ", err)
			break
		}
		log.Println("Data written to client")
	}
}

func broadcastPresencePeriodically(stopBroadcast chan bool) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(broadcastAddress),
		Port: broadcastPort,
	})
	if err != nil {
		log.Printf("failed to setup UDP broadcast: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(broadcastInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			msg := fmt.Sprintf("QUIC SERVER AVAILABLE: %s:%d", broadcastAddress, 4242)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("Failed to send broadcast message: %v", err)
			} else {
				log.Printf("Broadcasted server presence on %s:%d %s", broadcastAddress, broadcastPort, serverIdentifier)
			}

		case <-stopBroadcast:
			log.Println("Stopping server broadcast")
			return
		}
	}
}

func handleLlIncomingStream(listener quic.Listener) {
	ctx := context.Background()
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			log.Println("Problem accepting session: ", err)
			continue
		}
		stream, err := sess.AcceptStream(ctx)
		if err != nil {
			log.Println("Error accepting stream: ", err)
			continue // instead of terminating the entire function, we'll log the error and continue accepting streams
		}

		go handleServerStream(sess, stream)
	}

}

func handleLlConnection(conn quic.Connection) {
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("failed to accept stream %v", err)
			return
		}
		go handleServerStream(conn, stream)
	}
}

func handleServerStream(sess quic.Connection, stream quic.Stream) {
	log.Println("Handling new stream from:", sess.RemoteAddr())
	//read data from stream
	buf := make([]byte, 100)
	reader := bufio.NewReader(os.Stdin)

	for {
		n, err := stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client Closed the connection")
				return
			}
			log.Printf("Could not read stream %v", err)
			return
		}
		fmt.Printf("received data: %s\n", buf[:n])

		// Write a response back to the client
		fmt.Println("Server: Type a message to send to client: ")

		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Could not write to stream %v", err)
			return
		}
		message = strings.TrimSpace(message)

		//if user types 'exit', break loop
		if message == "exit" {
			fmt.Println("Connection closed")
			stream.Close()
			return
		}
		//send message to client
		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message: ", err)
			break
		}
		log.Println("Data written to client")
	}
}

func getLlCertificates() tls.Certificate {
	//load your certificates and keys
	cert, err := tls.LoadX509KeyPair("TLSKeys/cert.pem", "TLSKeys/key_without_passphrase.pem")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loading cert from:", "TLSKeys/cert.pem")
	fmt.Println("Loading key from:", "TLSKeys/key_without_passphrase.pem")
	return cert
}

func timeOfMessage() {
	fmt.Println(time.Now())
}

func generateTlsConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("TLSKeys/cert.pem", "TLSKeys/key_without_passphrase.pem")
	if err != nil {
		log.Fatal("Failed to load key pair!", err)
	}

	//create a certificate pool
	rootCas := x509.NewCertPool()

	//read in the cert file
	certs, err := ioutil.ReadFile("TLSKeys/cert.pem")
	if err != nil {
		log.Printf("Failed to append %q to RootCas: %v", "TLSKeys/cert.pem", err)
	}

	//append cert to the system pool
	if ok := rootCas.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs APPENDED, using system certs only")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCas,
		NextProtos:   []string{"quic-echo-example"},
		MinVersion:   tls.VersionTLS13,
	}

}

// func handleSession(session quic.Connection) {
//     for {
//         stream, err := session.AcceptStream(context.Background())
//         if err != nil {
//             log.Printf("Failed to accept stream: %v", err)
//             return
//         }
//         go handleStream(stream)
//     }
// }

// func handleStream(stream quic.Stream) {
//     // Handle the stream, e.g., read/write data
//     // For demonstration, we'll just read data and print it
//     buf := make([]byte, 4096)
//     n, err := stream.Read(buf)
//     if err != nil {
//         log.Printf("Failed to read from stream: %v", err)
//         return
//     }
//     fmt.Printf("Received: %s\n", buf[:n])
// }
