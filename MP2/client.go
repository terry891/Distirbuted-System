package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//Panic if an error is thrown
func check(e error) {
	if e != nil {
		fmt.Println("\n\nWe got an error \n%s", e)
		panic(e)
	}
}

//Initialize the connection with Server by sending IP-username pair
func initClient(serverIP string, myIP string, username string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		//Send server handshake request by providing IP-name with timeout
		c, err := net.DialTimeout("tcp", serverIP, 1*time.Second)
		check(err)
		fmt.Fprintf(c, myIP+"$$"+username+"$$INIT#")

		//Checking aknowledgement from the Server
		message, _ := bufio.NewReader(c).ReadString('\n')
		if message == "SUCCESS\n" {
			break
		} else {
			print("Connection now correctly established! Retrying... \n")
		}
	}
}

//Send message to server
func send(serverIP string, myIP string, username string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		//Read message-command input from user
		var message, toIPorName string
		for {
			print("\nIP:Port or username to send the message to:\n ")
			fmt.Scanf("%s", &toIPorName)

			//Exit the program
			if toIPorName == "EXIT" {
				println("Client side connection closes...\n")

				//Tell the client's recieve thread to close down
				connectionSelf, err := net.Dial("tcp", myIP)
				check(err)
				fmt.Fprintf(connectionSelf, "self$$EXIT#")

				//Tell server that this client is no longer active
				connectionServer, err := net.Dial("tcp", serverIP)
				check(err)
				fmt.Fprintf(connectionServer, username+"$$OFFLINE#")
				return
			}

			//Read content of the message from user
			fmt.Printf("Your message to %s plese:\n ", toIPorName)
			message, _ = bufio.NewReader(os.Stdin).ReadString('\n')

			//Send message to server with destination username
			connection, err := net.Dial("tcp", serverIP)
			check(err)
			fmt.Fprintf(connection, username+"$$"+toIPorName+"$$"+message+"#")
		}
	}
}

//Receive messages from the server
func receive(myIP string, wg *sync.WaitGroup) {
	defer wg.Done()

	//Create a TCP process at my IP address
	l, err := net.Listen("tcp", myIP)
	check(err)
	defer l.Close()

	for {
		//Read the message from tcp channel
		c, err := l.Accept()
		check(err)
		netData, err := bufio.NewReader(c).ReadString('#')
		check(err)
		message := strings.TrimSpace(string(netData))

		if message != "" {
			//Check if a EXIT signal is received
			if strings.Contains(message, "EXIT") {
				info := strings.Split(message, "$$")[0]
				if info == "self" {
					//Client signals EXIT
					return
				} else {
					//Server signals EXIT
					fmt.Printf("Server IP %s is down, client process terminated\n\n", info)
					fmt.Fprintf(os.Stdin, "EXIT\n")
					return
				}
			} else if strings.Contains(message, "NOT-FOUND") {
				//Print a message if the destination client IP address is not found
				notFound := strings.Split(message, "$$")[0]
				fmt.Printf("Sorry, your detination client %s is not found.\n\nIP:Port or username to send the message to:\n ", notFound)
				continue
			} else if strings.Contains(message, "INIT") {
				continue
			}

			//Print out the message received from the server
			parsed := strings.Split(message, "$$")
			fromName := parsed[0]
			toName := parsed[1]
			content := strings.TrimRight(parsed[2], "#")
			fmt.Printf("\n %s says to %s(you):\n   %s\n\nIP:Port or username to send the message to:\n ", fromName, toName, content)
		}
	}
}

//Main client function
func main() {
	//Initiate the client's information: ip, username and serverIP
	arguments := os.Args
	if len(arguments) != 4 {
		print("Please provide Server IP:Port, your IP:Port, and your name\n")
		return
	}

	hostIP := arguments[1]
	myIP := arguments[2]
	username := arguments[3]

	//Create the send and receive routines for this client, wait until EXIT signals both to terminate
	var wg sync.WaitGroup
	wg.Add(1)
	go initClient(hostIP, myIP, username, &wg)
	wg.Add(1)
	go send(hostIP, myIP, username, &wg)
	wg.Add(1)
	go receive(myIP, &wg)
	wg.Wait()
}
