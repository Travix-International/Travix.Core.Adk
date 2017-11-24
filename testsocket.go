package main

import (
	"fmt"
	"os"
	// "time"

	"golang.org/x/net/websocket"
)

const address string = "127.0.0.1:5000"

type Message struct {
	RequestID int
	Command   string
}

func main() {
	fmt.Println("Starting client")

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/socket.io/?EIO=3&transport=websocket", address), "", fmt.Sprintf("http://%s/", address))

	if err != nil {
		fmt.Print(err)
		fmt.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}
	incomingMessages := make(chan string)
	go readClientMessages(ws, incomingMessages)

	// i := 0
	for {
		select {
		/*
		   case <-time.After(time.Duration(2e9)):
		       i++
		       response := new(Message)
		       response.RequestID = i
		       response.Command = "Eject the hot dog."
		       err = websocket.JSON.Send(ws, response)
		       if err != nil {
		           fmt.Printf("Send failed: %s\n", err.Error())
		           os.Exit(1)
		       }
		*/
		case message := <-incomingMessages:
			fmt.Println(`Message Received:`, message)

		}
	}

}

func readClientMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		// err := websocket.JSON.Receive(ws, &message)
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}
