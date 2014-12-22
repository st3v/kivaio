package main

import (
	"fmt"

	"github.com/st3v/kiva-streams/socketio"
)

func main() {
	fmt.Println("Starting Client")

	session, err := socketio.NewSession("streams.kiva.org")
	if err != nil {
		fmt.Printf("Error opening session: %s\n", err.Error())
	}

	channels := []string{
		"loan.purchased",
		"loan.posted",
		"lender.registered",
		"lender.joinedTeam",
	}

	broadcast := make(chan string)

	for _, channelName := range channels {
		channel, err := session.Connect(channelName)
		if err != nil {
			fmt.Printf("Error connecting to channel '%s': %s\n", channelName, err.Error())
			continue
		}

		go func(name string, messages <-chan string) {
			for {
				broadcast <- fmt.Sprintf("%s: %s", name, <-messages)
			}
		}(channelName, channel)
	}

	for {
		fmt.Println(<-broadcast)
	}
}
