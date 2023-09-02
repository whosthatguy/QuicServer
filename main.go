package main

import (
	"crypto/sha256"
	"encoding/hex"

	//"log"
	"os/exec"
	//"runtime"
	//"syscall"

	//"flag"

	"fmt"
	"os"
	//"github.com/whosthatguy/Quic/client"
	//"github.com/whosthatguy/Quic/server"
	//"github.com/quic-go/quic-go"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run . <server|client>")
		return
	}
	switch os.Args[1] {
	case "server":
		StartServer()
		return
	case "client":
		StartClient()
		return
	default:
		fmt.Println("Invalid Choice: use 'server' or 'client'. ")
	}

}
func openNewTerminal(command string) {
	c := exec.Command("cmd.exe", "c", "start", "cmd.exe", "/k", command)
	err := c.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n ", err)
	}
}

// // Step 1: Create participants and generate their keys
// participants := make([]*Participant, 0)
// names := []string{"Alice", "Bob", "Charlie"}

// for _, name := range names {
// 	privKey, pubKey, err := GenerateKeys()
// 	if err != nil {
// 		fmt.Println("Error generating keys:", err)
// 		return
// 	}
// 	participant := &Participant{
// 		Name:       name,
// 		PrivateKey: privKey,
// 		PublicKey:  *pubKey,
// 	}
// 	participants = append(participants, participant)
// }

// // Step 2: Each participant signs a message
// message := []byte("This is a multisig contract!")
// signatures := make([][]byte, 0)

// for _, participant := range participants {
// 	signature, err := signMessage(participant.PrivateKey, message)
// 	if err != nil {
// 		fmt.Println("Error signing message:", err)
// 		return
// 	}
// 	signatures = append(signatures, signature)
// }

// // Step 3: Build the Merkle tree from the signatures
// merkleRoot, tree := buildMerkleTree(signatures)

// // Display the Merkle root
// fmt.Println("Merkle Root:", merkleRoot)

// // For demonstration purposes, let's generate a Merkle proof for Alice's signature and then verify it
// aliceSignature := signatures[0]
// proof := generateMerkleProof(aliceSignature, tree)

// // Verify the Merkle proof
// isValid := verifyMerkleProof(proof, aliceSignature, merkleRoot)
// if isValid {
// 	fmt.Println("Merkle proof for Alice's signature is valid!")
// } else {
// 	fmt.Println("Merkle proof for Alice's signature is invalid!")
// }



//automat the process of creating server and client
//if len(os.Args) > 1{
// 	switch os.Args[1] {
// 	case "handleServerStream":
// 		handleServerStream()
// 		return
// 	case "handleClientStream":
// 		handleClientStream()
// 		return
// 	case "server":
// 		StartServer()
// 		return
// 	case "client":
// 		StartClient()
// 	}
// }

// for {
// 	fmt.Println("Select an option:")
// 	fmt.Println("1: Start server")
// 	fmt.Println("2: Start client")
// 	fmt.Println("3: Start both server and client")
// 	fmt.Println("4: Exit")

// 	var choice int
// 	_, err := fmt.Scanln(&choice)
// 	if err != nil {
// 		fmt.Println("Error reading choice", err)
// 		continue
// 	}

// 	switch choice {
// 	case 1:
// 		fmt.Println("Starting server in new Window: ")
// 		StartServer()
// 	case 2:
// 		fmt.Println("Starting client in new window: ")
// 		StartClient()
// 	case 3:
// 		fmt.Println("starting Server and Client in a new window : ")
// 		StartServer()
// 		StartClient()
// 	case 4:
// 		fmt.Println("Exiting Program: ")
// 		os.Exit(0)
// 	default:
// 		fmt.Println("Invalid choice")
// 	}
// }
