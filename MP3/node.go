package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Panic if error is thrown
func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

//Add delay to the message before being sent off
func delay_send(configFile []string, myID string, myIP string, min int, max int, wg *sync.WaitGroup) {
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
				check(err)
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
	check(err)
	fmt.Fprintf(connection, myID+"$$"+message+"\n")
	fmt.Printf("Send \"%s\" to process %s, system time is %s, after delayed %d ms\n\n  ", message, id, now, delayMS)
}

//Print the received message when received any via TCP
func unicast_receive(source string, l net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	for {
		//Read the message from tcp channel
		c, err := l.Accept()
		check(err)
		netData, err := bufio.NewReader(c).ReadString('\n')
		check(err)
		message := strings.TrimSpace(string(netData))

		if message != "" {
			//Output the message content
			now := strings.Split(time.Now().String(), " ")[1]
			sourceID := strings.Split(message, "$$")[0]
			message = strings.Split(message, "$$")[1]

			//Close the connection if received STOP signal
			if message == "STOP" {
				if sourceID == source { //STOP signal initiated by self process, simply cloase unicast_receive routine
					fmt.Printf(" Stop signaled received from self! Closing self process ID: %s\n\n", sourceID)
					return
				} else { //STOP signal from another service, close down unicast_receive and unicast_write routine
					fmt.Printf(" Stop signal received from ID: %s; Initiate self termination ID: %s\n", sourceID, source)
					//FIXME: need to write to stdin to terminate the send go routine

					// var b bytes.Buffer
					// b.Write([]byte("STOP\n"))
					// os.Stdin = &b

					//os.Stdin.Write([]byte("STOP"))

					fmt.Fprintf(os.Stdin, "STOP\n")
					//bufio.NewWriter(os.Stdin).Write([]byte("STOP\n"))

					return
				}
			} else {
				fmt.Printf("Received \"%s\" from process %s, system time is %s\n\n  ", message, sourceID, now)
			}
		}
	}
}

//Main function
func main() {
	//Set the current node IP address
	arguments := os.Args
	if len(arguments) == 1 {
		print("Please provide host:port.\n")
		return
	}
	myIP := arguments[1]
	var myID string

	//Read the configuration file and test if the corresponding ID for the entered IP is in the configuration file
	file, err := ioutil.ReadFile("config.txt")
	check(err)
	configFile := strings.Split(string(file), "\n")
	for _, pair := range configFile[1:] {
		parsedID := strings.Split(pair, " ")
		if parsedID[1] == myIP {
			myID = parsedID[0]
		}
	}
	if myID == "" {
		print("The IP:Port you entered cannot be found in the configuration file\n")
		return
	}
	print("\n")

	//Get the min and max delay time as type integer
	lineOne := strings.Split(configFile[0], " ")
	min, _ := strconv.Atoi(lineOne[0])
	max, _ := strconv.Atoi(lineOne[1])

	//Create a TCP process at the given port
	l, err := net.Listen("tcp", myIP)
	check(err)
	defer l.Close()

	//Generate one thread for sending message and one thread for receiving message
	var wg sync.WaitGroup
	wg.Add(1)
	go delay_send(configFile, myID, myIP, min, max, &wg)
	wg.Add(1)
	go unicast_receive(myID, l, &wg)
	wg.Wait()

}
