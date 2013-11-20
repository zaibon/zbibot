package main

import (
	"flag"
	"fmt"
	"github.com/Zaibon/ircbot"
)

type channels []string

func (c *channels) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *channels) Set(value string) error {
	*c = append(*c, value)
	return nil
}

var (
	flagServer   string
	flagPort     string
	flagSsl      bool
	flagChannels channels
)

func init() {
	flag.StringVar(&flagServer, "server", "irc.freenode.net", "ip adresse of the server you want to connect to")
	flag.StringVar(&flagPort, "port", "6697", "port")
	flag.BoolVar(&flagSsl, "ssl", true, "true|false")

}

func main() {
	flag.Var(&flagChannels, "c", "channels")

	flag.Parse()

	b := ircbot.NewIrcBot()
	b.Server = flagServer
	b.Port = flagPort
	b.Encrypted = flagSsl
	b.Nick = "ZbiBotTLS"
	b.User = b.Nick

	if flag.NFlag() != 0 {
		for i := 0; i < len(flagChannels); i++ {
			b.Channel = append(b.Channel, flagChannels[i])
		}
	}

	fmt.Println(b)

	b.Handlers["PING"] = ircbot.Pong
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
