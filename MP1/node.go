package main

import (
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	util "./functions"
)

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
	util.Check(err)
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
	util.Check(err)
	defer l.Close()

	//Generate one thread for sending message and one thread for receiving message
	var wg sync.WaitGroup
	wg.Add(1)
	go util.Delay_send(configFile, myID, myIP, min, max, &wg)
	wg.Add(1)
	go util.Unicast_receive(myID, l, &wg)
	wg.Wait()

}
