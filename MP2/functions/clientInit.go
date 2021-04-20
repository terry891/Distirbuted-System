package functions

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

//Initialize the connection with Server by sending IP-username pair
func ClientInit(serverIP string, myIP string, username string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		//Send server handshake request by providing IP-name with timeout
		c, err := net.DialTimeout("tcp", serverIP, 1*time.Second)
		Check(err)
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
