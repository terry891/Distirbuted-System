package functions

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//Add delay to the message before being sent off
func Delay_send(configFile []string, myID string, myIP string, min int, max int, wg *sync.WaitGroup) {
	defer wg.Done()

	//An overall for loop allowing sending messages to multiple processes
	for {

		//Read message-command input from user
		var message, destinationIP, destinationID string
		for {
			fmt.Print("Please write destination ID and message:\n  ")
			rawCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if rawCommand == "STOP\n" {
				//Terminate the current process if STOP signal received, also send self a STOP signal
				fmt.Printf("\n Stop signal received, terminating current process ID %s\n", myID)
				connection, err := net.Dial("tcp", myIP)
				Check(err)
				fmt.Fprintf(connection, myID+"$$STOP\n")
				return
			}
			command := strings.Split(rawCommand, " ")

			//Prepare for sending by getting IP address and parsing input command
			if command[0] == "send" && len(command) == 3 {
				message = strings.TrimRight(command[2], "\n")
				destinationID = command[1]

				//Find the correct destination IP address from configuration file given process ID
				for _, pair := range configFile[1:] {
					parsedID := strings.Split(pair, " ")
					if parsedID[0] == destinationID {
						destinationIP = parsedID[1]
					}
				}
				//Continue ask for message input if ID is not found
				if destinationIP != "" {
					print("\n")
					break
				} else {
					print("Sorry your ID is not found in the configuration file\n")
					continue
				}
			}
		}

		//Generate the time to delay based on min and max, sending this amount to a new goroutine
		delayMS := min + rand.Intn(max-min)

		//New go routine to execute the delay and message send
		wg.Add(1)
		go unicast_send(destinationIP, destinationID, myID, message, delayMS, wg)
	}
}

//Send message to destination
func unicast_send(ip string, id string, myID string, message string, delayMS int, wg *sync.WaitGroup) {
	defer wg.Done()

	//Delay
	time.Sleep(time.Duration(delayMS) * time.Millisecond)

	//Send message to destination ip
	now := strings.Split(time.Now().String(), " ")[1] //Get the current time in hh:mm:ss:ms
	connection, err := net.Dial("tcp", ip)
	Check(err)
	fmt.Fprintf(connection, myID+"$$"+message+"\n")
	fmt.Printf("Send \"%s\" to process %s, system time is %s, after delayed %d ms\n\n  ", message, id, now, delayMS)
}
