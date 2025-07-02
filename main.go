package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	_ "modernc.org/sqlite"
)

func main() {
	container, err := sqlstore.New("sqlite", "file:session.db?_foreign_keys=on", nil)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	client := whatsmeow.NewClient(deviceStore, nil)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			for _, msg := range v.Messages {
				if msg.Message.GetConversation() != "" {
					fmt.Println("Received message:", msg.Message.GetConversation())
					// Send reply
					_, err := client.SendMessage(context.Background(), msg.Info.Chat, &waProto.Message{
						Conversation: proto.String("Hello, this is a simple WA bot!"),
					})
					if err != nil {
						fmt.Println("Failed to send message:", err)
					}
				}
			}
		}
	})

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	// Wait for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect()
}
