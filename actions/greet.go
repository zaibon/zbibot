package actions

import (
	"fmt"

	"github.com/Zaibon/ircbot"
)

type Greet struct{}

func (g *Greet) Command() []string {
	return []string{"JOIN"}
}

func (g *Greet) Usage() string {
	return ""
}

func (g *Greet) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if m.Nick == b.Nick {
		b.Joined = true
		return
	}

	s := fmt.Sprintf("%s :Salut %s", m.Channel, m.Nick)
	b.Out <- &ircbot.IrcMsg{
		Command: "PRIVMSG",
		Args:    []string{s},
	}
}
