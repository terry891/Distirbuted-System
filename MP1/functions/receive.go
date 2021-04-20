package functions

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//Print the received message when received any via TCP
func Unicast_receive(source string, l net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	for {
		//Read the message from tcp channel
		c, err := l.Accept()
		Check(err)
		netData, err := bufio.NewReader(c).ReadString('\n')
		Check(err)
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
