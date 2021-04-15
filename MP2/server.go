package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

//Panic if an error is thrown
func check(e error) {
	if e != nil {
		fmt.Println("\n\nWe got an error \n%s", e)
		panic(e)
	}
}

//Server forwards message to designated client by its IP address
func serverSend(sIP string, toIP string, fromName string, toName string, content string, wg *sync.WaitGroup) {
	defer wg.Done()

	//Construct message, 2 special cases: exit and destination client IP not found
	var message string
	if content == "EXIT" {
		message = sIP + "$$EXIT#"
	} else if content == "NOT-FOUND" {
		message = fromName + "$$NOT-FOUND#"
	} else {
		message = fromName + "$$" + toName + "$$" + content + "#"
	}

	//Send message to destination IP
	connection, err := net.Dial("tcp", toIP)
	check(err)
	fmt.Fprintf(connection, message)
}

//Pass the received message from client to data channel
func serverReceive(serverIP string, mutex *sync.Mutex, data chan<- []string, clientMap *map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	//Create a TCP listener at the given server IP
	l, err := net.Listen("tcp", serverIP)
	check(err)
	defer l.Close()

	//Make an exit channel, so that the individual subroutine can signal exit
	for {
		//Read the message from tcp channel
		c, err := l.Accept()
		check(err)
		netData, err := bufio.NewReader(c).ReadString('#')
		check(err)
		message := strings.TrimRight(string(netData), "#")

		//Check if an EXIT signal is received, if so shut off server
		if strings.Contains(message, "EXIT") {
			return
		}

		//Generate a go routine for each client request so that multiple clients can connect concurrently
		wg.Add(1)
		go func(serverIP string, mutex *sync.Mutex, c net.Conn, clientMap *map[string]string, data chan<- []string, wg *sync.WaitGroup) {
			defer wg.Done()

			if message != "" {
				//Check special cases
				if strings.Contains(message, "OFFLINE") {
					//A client is going offline, take away the name-IP pair
					name := strings.Split(message, "$$")[0]
					fmt.Printf("Client %s disconnected. Its IP %s is dropped.\n", name, (*clientMap)[name])
					mutex.Lock()
					delete(*clientMap, name)
					mutex.Unlock()
					return
				} else if strings.Contains(message, "INIT") {
					//A new client is just connected, store its name-IP pair
					info := strings.Split(message, "$$")
					fmt.Printf(" Client %s with IP %s just joined the Server...\n", info[1], info[0])
					mutex.Lock()
					(*clientMap)[info[1]] = info[0]
					mutex.Unlock()

					//Reply with an acknolege message
					_, err = c.Write([]byte("SUCCESS\n"))
					check(err)
					return
				}

				//Pass the parsed message content to the data channel (main will take care of this)
				info := strings.Split(message, "$$")
				data <- info
			}
		}(serverIP, mutex, c, clientMap, data, wg)
	}
}

//Main function
func main() {
	//Set the server IP address from user input
	arguments := os.Args
	if len(arguments) == 1 {
		print("Please provide host:port.\n")
		return
	}
	myIP := arguments[1]

	//Set a map serving as DNS lookup and a data channel to pass message to client threads
	var clientMap = make(map[string]string)
	var mutex = &sync.Mutex{}
	data := make(chan []string, 5)

	//Generate threads for the functionality of Server
	var wg sync.WaitGroup
	wg.Add(1)
	go serverReceive(myIP, mutex, data, &clientMap, &wg)

	//Read user inputs from the terminal, waiting for an EXIT input
	wg.Add(1)
	go func(data chan<- []string) {
		defer wg.Done()
		//Only pass in an EXIT signal to channel when EXIT is typped
		for {
			exitReady, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if strings.Contains(exitReady, "EXIT") {
				toexit := []string{"EXIT", " ", " "}
				data <- toexit
				print("\nEXIT signal received\n")
				return
			}
		}
	}(data)

	//Dealing with return values from channels and pass on to function serverSend
	exitNow := false //Help to break the following for loop
	for {
		//Keep getting values from channel, skip if there isn't any value ready
		select {
		case dataDump, _ := <-data:
			// if (len(dataDump)) == 0 {
			// 	continue
			// }

			//Get parsed message information from the data channel
			fromName := dataDump[0]
			toNameOrIP := dataDump[1]
			message := dataDump[2]
			toIP := clientMap[toNameOrIP]
			var toName string

			//Baordcast EXIT to all client if EXIT signal received on Server
			if fromName == "EXIT" {
				for cName, cIP := range clientMap {
					wg.Add(1)
					go serverSend(myIP, cIP, "", "", "EXIT", &wg)
					fmt.Printf(" Signal client %s to terminate...\n", cName)
				}
				wg.Add(1)
				go serverSend(myIP, myIP, "", "", "EXIT", &wg)
				print(" Signal Server to close down.\n")
				close(data)
				exitNow = true //Break out of for loop
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
					go serverSend("", name, toNameOrIP, "", "NOT-FOUND", &wg)
					fmt.Printf("Client %s sends a message, but no one receives it...\n", name)
					continue
				}
			} else {
				toName = toNameOrIP
			}

			//Forward the message to the appropriate client address
			wg.Add(1)
			go serverSend(myIP, toIP, fromName, toName, message, &wg)
			fmt.Printf("Client %s sends a message to Client %s\n", fromName, toName)
		default:
			print("")
		}

		if exitNow {
			break
		}
	}

	wg.Wait()
	print("Exited cleanly!! :)\n\n")
}
