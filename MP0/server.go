//Use go run server.go 127.0.0.1:1234 to run this program

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

//Upon receving message from the client, print the email in a nice format
func prettyPrint(mess string) {
	data := strings.Split(mess, "$$")
	fmt.Printf("\nEmail recieved!\n")
	fmt.Printf("  Titled %s\n  From %s to %s with data %s\n  %s\n", data[3], data[1], data[0], data[2], data[4])
}

//main() listen on the given port for message, and send off a terminate signal
func main() {
	//If user does not enter any server IP:Port address, use a default port address; otherwise store the given IP in "myIP"
	arguments := os.Args
	myIP := ":"
	if len(arguments) == 1 {
		myIP += "1234"
	} else {
		myIP += arguments[1]
	}

	//Server listen on port myIP for any TCP connection, close the listening port after terminating the program
	l, err := net.Listen("tcp", myIP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	//l.Acccept waits for the next call and returns a generic Conn, where message can be extracted
	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Server has already received a message from l.Accept(), now use Reader to read the string content, and stored in var message
	message, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	//Print the message from the client nicely and tell client to close down
	if strings.TrimSpace(string(message)) != "" {
		prettyPrint(string(message))
		c.Write([]byte("Stop")) //Write to client that the message is received and it should shut down
		fmt.Print("\nStop signal received, successful!\n\n")
		return
	}
}
