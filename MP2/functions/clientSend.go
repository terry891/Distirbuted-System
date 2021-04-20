package functions

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

//Send message to server
func ClientSend(serverIP string, myIP string, username string, wg *sync.WaitGroup) {
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
				Check(err)
				fmt.Fprintf(connectionSelf, "self$$EXIT#")

				//Tell server that this client is no longer active
				connectionServer, err := net.Dial("tcp", serverIP)
				Check(err)
				fmt.Fprintf(connectionServer, username+"$$OFFLINE#")
				return
			}

			//Read content of the message from user
			fmt.Printf("Your message to %s plese:\n ", toIPorName)
			message, _ = bufio.NewReader(os.Stdin).ReadString('\n')

			//Send message to server with destination username
			connection, err := net.Dial("tcp", serverIP)
			Check(err)
			fmt.Fprintf(connection, username+"$$"+toIPorName+"$$"+message+"#")
		}
	}
}
