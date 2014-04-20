package actions

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Zaibon/ircbot"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Respond struct{}

func (r *Respond) Command() []string {
	return []string{"PRIVMSG"}
}

func (r *Respond) Usage() string {
	return "respond when someone say my name"
}

func (r *Respond) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	response := []string{
		"oui ?",
		"on parle de moi ?",
		"Je suis pas là",
	}

	s := strings.Join(m.Args, " ")

	if strings.Contains(s, b.Nick) {
		nbr := rand.Intn(len(response))
		line := fmt.Sprintf(":%s", response[nbr])
		b.Out <- &ircbot.IrcMsg{
			Command: "PRIVMSG",
			Channel: m.Channel,
			Args:    []string{line},
		}
	}

}
