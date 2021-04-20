package functions

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

//Pass the received message from client to the shared messageData array
func ServerReceive(serverIP string, mutex *sync.Mutex, cond *sync.Cond, messageData *[][]string, clientMap *map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	//Create a TCP listener at the given server IP
	l, err := net.Listen("tcp", serverIP)
	Check(err)
	defer l.Close()

	for {
		//Read the message from tcp port
		c, err := l.Accept()
		Check(err)
		netData, err := bufio.NewReader(c).ReadString('#')
		Check(err)
		message := strings.TrimRight(string(netData), "#")

		//Check if an EXIT signal is received, if so shut off server
		if strings.Contains(message, "EXIT") {
			return
		}

		//Generate a go-routine for each client request so that multiple clients can connect concurrently
		wg.Add(1)
		go func(serverIP string, mutex *sync.Mutex, cond *sync.Cond, c net.Conn, clientMap *map[string]string, messageData *[][]string, wg *sync.WaitGroup) {
			defer wg.Done()

			if message != "" {
				//Check special cases
				if strings.Contains(message, "OFFLINE") {
					//A client is going offline, take away the name-IP pair
					name := strings.Split(message, "$$")[0]
					fmt.Printf("Client %s disconnected. Its IP %s is dropped.\n", name, (*clientMap)[name])
					mutex.Lock()
					delete(*clientMap, name)
					mutex.Unlock()
					return
				} else if strings.Contains(message, "INIT") {
					//A new client is just connected, store its name-IP pair
					info := strings.Split(message, "$$")
					fmt.Printf(" Client %s with IP %s just joined the Server...\n", info[1], info[0])
					mutex.Lock()
					(*clientMap)[info[1]] = info[0]
					mutex.Unlock()

					//Reply with an acknolege message
					_, err = c.Write([]byte("SUCCESS\n"))
					Check(err)
					return
				}

				//Pass the parsed message content to the messageData array (main server thread will then send it off to destination client specficied in the message)
				info := strings.Split(message, "$$")
				mutex.Lock()
				*messageData = append(*messageData, info)
				cond.Broadcast()
				mutex.Unlock()
			}
		}(serverIP, mutex, cond, c, clientMap, messageData, wg)
	}
}
