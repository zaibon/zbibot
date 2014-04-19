package main

import (
	"flag"
	"fmt"

	"github.com/Zaibon/ircbot"
	"github.com/Zaibon/zbibot/actions"
)

//needed for the flag "channel"
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

	flagNick string

	flagWebEnable bool
	flagWebPort   string
)

func init() {
	flag.StringVar(&flagServer, "server", "irc.freenode.net", "ip adresse of the server you want to connect to")
	flag.StringVar(&flagServer, "s", "irc.freenode.net", "ip adresse of the server you want to connect to")

	flag.StringVar(&flagPort, "port", "6697", "port")
	flag.StringVar(&flagPort, "p", "6697", "port")

	flag.BoolVar(&flagSsl, "ssl", true, "true|false")

	flag.StringVar(&flagNick, "nick", "ZbiBot", "nickname")
	flag.StringVar(&flagNick, "n", "ZbiBot", "nickname")

	flag.BoolVar(&flagWebEnable, "web", false, "enable or not the web interface true|false")
	flag.StringVar(&flagWebPort, "wport", "6697", "port on wich to bind web interface")
}

func main() {
	flag.Var(&flagChannels, "c", "channels")

	flag.Parse()

	//create new bot
	b := ircbot.NewIrcBot()

	//configure bot
	b.Server = flagServer
	b.Port = flagPort
	b.Encrypted = flagSsl
	b.Nick = flagNick
	b.User = b.Nick
	if flag.NFlag() != 0 {
		for i := 0; i < len(flagChannels); i++ {
			b.Channel = append(b.Channel, flagChannels[i])
		}
	}
	b.WebEnable = flagWebEnable
	b.WebPort = flagWebPort

	//set channels
	b.AddInternAction(&actions.Greet{})
	b.AddInternAction(&actions.Respond{})

	//command fire by users
	b.AddUserAction(&actions.Help{})

	//connectin to server, listen and serve
	b.Connect()

	//TODO handle signal system to throw something in b.Exit
	<-b.Exit
	//and then disconenct
	b.Disconnect()
}
