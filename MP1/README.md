# MP1

## Established peer-to-perr TCP connection with simulated delay among many nodes according to a configuration file

### 1. Inputs and Setup

* Configuration file: first line has the min and max value of simulated delay time in milliseconds. The lines followed contian the IP:Port address paired with ID serving as a dictionary. Each node can look up this configuration file to find the IP address by ID. 

* Node: upon receiving a message from user with a destination ID, the node.go look up the ID's IP address in the configuration text file. If the ID is not found, print a warning message. If found, generate a thread to simulate the delay using time.Sleep() for a random amount of time between min and max delay time, and then send the message via TCP.

  

### 2. Runing & Testing

```tex
Change the configuration file to the desire min and max delay time as well as the IP:Port address for possible nodes
```

``` sh
go run node.go 127.0.0.1:1234
go run node.go 127.0.0.1:1235
go run node.go 127.0.0.1:1236

```

Message Format:

```tex
send ID message
   send 2 Hello!  (example)
STOP
	 (stop signal to terminate the current node)
```



### 3. Underlying Structure

**Note**: each numbered block is a functionaniltiy available to node

**Note**: <u>each colored rectangle block represents a thread/go-routine</u>

1. node can receive and send message concurrently. Each out going message is sent to a *delay* go-routine to simulate delay using time.Sleep()
2. node can send *multiple* and receive concurrently. The amount of time to delay is selected randomly
3. node prints a warning if destination node ID cannot be found in the config.txt file. Node does not generate any other thread (no delay thread)
4. node receives user input 'STOP' and terminates the receive thread that is constantly listening
5. node will panic if it tries to send message to an already closed node

![](./imgs/diagram.png)



### 4. Codes in Action

Node 1:

![](./imgs/1.png)



Node 2:

![](./imgs/2.png)



Code in action:

![Run](./imgs/run.gif)



 ### 4. Credits

Using the Create a TCP and UDP Client and Server using Go article from Linode: [Source](https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/)





