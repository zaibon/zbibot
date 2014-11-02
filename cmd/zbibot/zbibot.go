package main

import (
	"fmt"

	"github.com/zaibon/zbibot"

	"github.com/jessevdk/go-flags"
	"github.com/zaibon/ircbot"
	"github.com/zaibon/ircbot/actions"
)

var opts struct {
	Server   string   `short:"s" long:"server" description:"ip adresse of the server you want to connect to" default:"irc.freenode.net"`
	Port     uint     `short:"p" long:"port" description:"port to connect to" default:"6667"`
	Channels []string `short:"c" long:"channels" description:"channels the bot has to joined" required:"true"`
	SSL      int      `long:"ssl" description:"enable ssl on not" default:"false"`
	Nick     string   `short:"n" long:"nick" description:"nickname" default:"Zbibot"`
	Password string   `long:"pass" long:"password" description:"password"`
	DBPath   string   `long:"db" long:"database" description:"path to the sqlite database file" default:"irc.db"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	//create new bot
	b := ircbot.NewIrcBot(opts.Nick, opts.Nick, opts.Password, opts.Server, opts.Port, opts.Channels, opts.DBPath)

	b.AddInternAction(&actions.Greet{})
	b.AddInternAction(actions.NewTitleExtract())
	b.AddInternAction(actions.NewLogger(b))
	b.AddInternAction(actions.NewURLLog(b))

	//command fire by users
	b.AddUserAction(&actions.Help{})
	b.AddUserAction(zbibot.NewLetMeKnow(b))
	b.AddUserAction(actions.NewURL(b))

	//connectin to server, listen and serve
	b.Connect()

	// //TODO handle signal system to throw something in b.Exit
	<-b.Exit
	//and then disconenct
	b.Disconnect()
}
