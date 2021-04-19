//Use go run client.go 127.0.0.1:1234 to run this program

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//A message struct to store the user inputs
type message struct {
	to      string
	from    string
	date    string
	title   string
	content string
}

//Ask user for input message, populating the struct
func readInput(mess *message) {
	fmt.Print("\n\nPlease enter destination email: ")
	fmt.Scanf("%s", &mess.to)
	fmt.Print("Please enter your email: ")
	fmt.Scanf("%s", &mess.from)
	fmt.Print("Please enter the date: ")
	fmt.Scanf("%s", &mess.date)
	fmt.Print("Please enter email title: ")
	fmt.Scanf("%s", &mess.title)
	fmt.Print("Please enter the content: ")
	mess.content, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}

//Prepare the message to send to server, using $$ as the delimiter
func prepareStr(mess *message) string {
	data := mess.to + "$$" + mess.from + "$$" + mess.date + "$$" + mess.title + "$$" + mess.content
	return data
}

//Main function for the client node, ask for input and send it off to server
func main() {
	//Read input from the user, and store the message ready to sent in stringData
	var mess message
	readInput(&mess)
	stringData := prepareStr(&mess)

	//If user does not enter any server IP:Port address, use a default port address; otherwise store the givin IP in "serverIP"
	arguments := os.Args
	serverIP := ""
	if len(arguments) == 1 {
		serverIP = "127.0.0.1:1234"
	} else {
		serverIP = arguments[1]
	}

	//Dial server IP, making a TCP connection to serverIP address
	c, err := net.Dial("tcp", serverIP)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Write the message to server using Fprintf, where c is the net.Connection, representing server's IO
	check := true
	if check {
		fmt.Fprintf(c, stringData+"\n")
		check = false
	}

	//Waiting for a message from server, server will be writing to the c (net.Conn) object
	//If the message from server is Stop, then terminate the program, if not keep waiting
	for {
		message, _ := bufio.NewReader(c).ReadString('\n')
		if message == "Stop" {
			fmt.Print("\nStop signal received, successful!\n\n\n")
			return
		}
	}
}
