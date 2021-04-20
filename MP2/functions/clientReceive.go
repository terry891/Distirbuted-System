package functions

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	Red   = "\033[31m"
	Reset = "\033[0m"
)

//Receive messages from the server
func ClientReceive(myIP string, wg *sync.WaitGroup) {
	defer wg.Done()

	//Create a TCP process at my IP address
	l, err := net.Listen("tcp", myIP)
	Check(err)
	defer l.Close()

	for {
		//Read the message from tcp channel
		c, err := l.Accept()
		Check(err)
		netData, err := bufio.NewReader(c).ReadString('#')
		Check(err)
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
			print(Red + "\nMessage Received: " + Reset)
			fmt.Printf(" %s says to %s(you):\n   %s\n\nIP:Port or username to send the message to:\n ", fromName, toName, content)
		}
	}
}
