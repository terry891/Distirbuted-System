package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	util "./functions"
)

var (
	wg sync.WaitGroup
)

//Main function
func main() {

	var myIP string
	arguments := os.Args

	//Set the server IP address from user input
	if len(arguments) == 1 {
		print("Please indicate if you're creating a server or client\n")
		return
	}

	//Check if the user wants to create a clinet or server
	option := arguments[1]
	if option == "server" {

		//Check if thre are enough arguments to create a server node
		if len(arguments) != 3 {
			print("Please provide Server IP:Port\n")
			return
		}

		//Initiate the server's information: ip
		myIP = arguments[2]
		wg.Add(1)
		go server(myIP)

	} else if option == "client" {

		//Check if thre are enough arguments to create a client node
		if len(arguments) != 5 {
			print("Please provide Server IP:Port, your IP:Port, and your name\n")
			return
		}

		//Initiate the client's information: ip, username and serverIP
		hostIP := arguments[2]
		myIP = arguments[3]
		username := arguments[4]
		wg.Add(1)
		go client(hostIP, myIP, username)

	} else {
		print("Please indicate if you're creating a server or client\n")
		return
	}

	//Wait for all thread the finish and then exit
	wg.Wait()
	print("Exited cleanly!! :)\n\n")
}

//Client function
func client(hostIP string, myIP string, username string) {
	defer wg.Done()

	//Create the send and receive routines for this client, wait until EXIT signals both to terminate
	wg.Add(1)
	go util.ClientInit(hostIP, myIP, username, &wg)
	wg.Add(1)
	go util.ClientSend(hostIP, myIP, username, &wg)
	wg.Add(1)
	go util.ClientReceive(myIP, &wg)
}

func server(myIP string) {
	defer wg.Done()

	//Set a map serving as DNS lookup and a shared data array to store and pass message to client threads
	var clientMap = make(map[string]string)
	var messageData [][]string
	var mutex = &sync.Mutex{}
	cond := sync.NewCond(mutex)

	//Generate threads for the functionality of Server
	wg.Add(1)
	go util.ServerReceive(myIP, mutex, cond, &messageData, &clientMap, &wg)

	//Read user inputs from the terminal, waiting for an EXIT input
	wg.Add(1)
	go func(messageData *[][]string, mutex *sync.Mutex) {
		defer wg.Done()

		for {
			//Pass only EXIT signal to the shared messageData array when EXIT is typped
			exitReady, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if strings.Contains(exitReady, "EXIT") {
				toexit := []string{"EXIT", " ", " "}

				//Add appread the STOP signal to messageData and wake up the waiting main server thread
				mutex.Lock()
				*messageData = append(*messageData, toexit)
				cond.Broadcast()
				mutex.Unlock()

				print("\nEXIT signal received\n")
				return
			}
		}
	}(&messageData, mutex)

	//Generate a send thread as soon as there is message in messageData array
	for {
		//parsedMessage is an array of info of a specific message received by serverReceive
		var parsedMessage []string

		//Wake up only if there is some message are pushed to messageData
		mutex.Lock()
		for len(messageData) == 0 {
			cond.Wait()
		}
		parsedMessage, messageData = messageData[0], messageData[1:]
		mutex.Unlock()

		//Get parsed message info from the shared messageData array, where serverReceive threads poplate
		fromName := parsedMessage[0]
		toNameOrIP := parsedMessage[1]
		message := parsedMessage[2]
		toIP := clientMap[toNameOrIP]
		var toName string

		//Baordcast EXIT to all client if EXIT signal received on Server
		if fromName == "EXIT" {
			for cName, cIP := range clientMap {
				wg.Add(1)
				go util.ServerSend(myIP, cIP, "", "", "EXIT", &wg)
				fmt.Printf(" Signal client %s to terminate...\n", cName)
			}
			wg.Add(1)
			go util.ServerSend(myIP, myIP, "", "", "EXIT", &wg)
			print(" Signal Server to close down.\n")
			break
		}

		//Return error message if the destination client is not found
		if toIP == "" {
			ipExist := false

			//Check if toNameOrIP is actually an IP or invalid
			for cName, cIP := range clientMap {
				if cIP == toNameOrIP {
					toName = cName
					toIP = toNameOrIP //The client entered IP, not destination username
					ipExist = true
					break
				}
			}

			//If the destination user is not found, return a NOT-FOUND signal
			if !ipExist {
				name := clientMap[fromName]
				wg.Add(1)
				go util.ServerSend("", name, toNameOrIP, "", "NOT-FOUND", &wg)
				fmt.Printf("Client %s sends a message, but no one receives it...\n", name)
				continue
			}
		} else {
			toName = toNameOrIP
		}

		//Forward the message to the appropriate client address
		wg.Add(1)
		go util.ServerSend(myIP, toIP, fromName, toName, message, &wg)
		fmt.Printf("Client %s sends a message to Client %s\n", fromName, toName)
	}
}
