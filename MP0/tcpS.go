package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type message struct {
	to      string
	from    string
	data    string
	title   string
	content string
}

func prettyPrint(mess string) string {
	return ""
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

	netData, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// if strings.TrimSpace(string(netData)) != "" {
	fmt.Print("-> ", string(netData))
	c.Write([]byte("Stop"))
	fmt.Print("\nStop signal received, successful!\n\n")
	return
	// }

}