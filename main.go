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

// //linked list
// myList := linkedList{}
// n1 := &node{data: 48}
// n2 := &node{data: 15}
// n3 := &node{data: 16}
// n4 := &node{data: 117}
// n5 := &node{data: 59}
// n6 := &node{data: 5}

// myList.prePend(n1)
// myList.prePend(n2)
// myList.prePend(n3)
// myList.prePend(n4)
// myList.prePend(n5)
// myList.prePend(n6)
// myList.printListData()
// fmt.Println(myList.string())
// myList.deleteWithValue(117)
// myList.printListData()
// fmt.Println(myList.string())

// //end of linked list

// hashmap := newHashMap()

// hashmap.add("Karson")
// hashmap.add("Colby")
// hashmap.add("Adam")

// hashmap.displayAll()

// hash := hashmap.hashString("Karson")
// original := hashmap.Get(hash)
// fmt.Println("\nRetrieved using hash: ", original)

// data := []byte("Hello, world!")
// fmt.Println(data)
// //insecure := flag.Bool("insecure", false, "skip certificate verification")
// flag.Parse()

// s := []int{1, 2, 3, -9, 4, 0}
// c := make(chan int)
// go sum(s[:len(s)/2], c) //first half of slice
// go sum(s[len(s)/2:], c) //second half of slice
// x, y := <-c, <-c        //receive from channel c
// fmt.Println(x, y, x+y)

func loops() {
	for i := 0; i < 10; i++ {
		if i < 5 {
			fmt.Println(i)
		}

		if i > 5 {
			fmt.Println("greater than 5")
		}
	}
}

type hashMap struct {
	data map[string]string
}

// newHashMap creates a new hash map instance
func newHashMap() *hashMap {
	return &hashMap{
		data: make(map[string]string),
	}
}

// hash string using sha256
func (h *hashMap) hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// adds string to hashMap
func (h *hashMap) add(s string) {
	hash := h.hashString(s)
	h.data[hash] = s
}

// retrieves a string from hashMap using its hash
func (h *hashMap) Get(hash string) string {
	return h.data[hash]
}

// display all strings and hashes in hashMap
func (h *hashMap) displayAll() {
	for hash, original := range h.data {
		fmt.Println("Original:", original)
		fmt.Println("Hash:", hash)
	}
}

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
