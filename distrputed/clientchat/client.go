package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer client.Close()

	// Join the chat
	var joinReply struct{ ID int }
	err = client.Call("ChatServer.Join", struct{}{}, &joinReply)
	if err != nil {
		fmt.Println("Join error:", err)
		return
	}
	id := joinReply.ID
	fmt.Printf("You joined as User %d\n", id)

	// Goroutine to listen for broadcasts
	go func() {
		for {
			var msgs []string
			err := client.Call("ChatServer.Stream", id, &msgs)
			if err == nil {
				for _, m := range msgs {
					fmt.Println(m)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Main loop: read input and send
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "exit" {
			fmt.Println("Exiting...")
			break
		}
		client.Call("ChatServer.SendMessage", struct {
			ID      int
			Message string
		}{ID: id, Message: text}, &struct{}{})
	}
}
