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
	date    string
	title   string
	content string
}

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

func prepareStr(mess *message) string {
	data := mess.to + "$$" + mess.from + "$$" + mess.date + "$$" + mess.title + "$$" + mess.content
	return data
}

func main() {
	var mess message
	readInput(&mess)
	stringData := prepareStr(&mess)

	arguments := os.Args
	CONNECT := ""

	if len(arguments) == 1 {
		CONNECT = "127.0.0.1:1234"

	} else {
		CONNECT = arguments[1]
	}

	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	check := true
	if check {
		fmt.Fprintf(c, stringData+"\n")
		check = false
	}

	message, _ := bufio.NewReader(c).ReadString('\n')
	if message == "Stop" {
		fmt.Print("\nStop signal received, successful!\n\n\n")
		return
	}
}
