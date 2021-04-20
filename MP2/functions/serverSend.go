package functions

import (
	"fmt"
	"net"
	"sync"
)

//Server forwards message to designated client by its IP address
func ServerSend(sIP string, toIP string, fromName string, toName string, content string, wg *sync.WaitGroup) {
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
	Check(err)
	fmt.Fprint(connection, message)
}
