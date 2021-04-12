package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type message struct {
	to      string
	from    string
	data    string
	title   string
	content string
}

func prettyPrint(mess string) {
	data := strings.Split(mess, "$$")
	fmt.Printf("\nEmail recieved!\n")
	fmt.Printf("  Titled %s\n  From %s to %s with data %s\n  %s\n", data[3], data[1], data[0], data[2], data[4])
}

func main() {
	arguments := os.Args
	PORT := ":"

	if len(arguments) == 1 {
		PORT += "1234"
	} else {
		PORT += arguments[1]
	}

	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.TrimSpace(string(netData)) != "" {
			prettyPrint(string(netData))
			c.Write([]byte("Stop")) //Write to client that the message is received
			fmt.Print("\nStop signal received, successful!\n\n")
			return
		}

	}
}
