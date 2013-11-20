package main

import (
	"github.com/Zaibon/ircbot"
	"log"
	// "time"
)

func main() {
	log.SetPrefix("irc> ")

	b := ircbot.NewIrcBot()
	b.Server = "irc.freenode.net"
	b.Port = "6667"
	b.Nick = "ZbiBot"
	b.User = b.Nick

	b.Channel = append(b.Channel, "#testgigx")

	b.Handlers["PONG"] = ircbot.Pong
	b.Handlers["JOIN"] = ircbot.Join
	b.Handlers["PRIVMSG"] = ircbot.Respond

	b.Connect()

	b.Listen()
	b.HandleActionIn()
	b.HandleActionOut()

	b.Join()

	//TODO handle signal system to throw something in b.Exit
	<-b.Exit
	//and then disconenct
	b.Disconnect()
}
