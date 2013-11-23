package zbibot

import (
	"fmt"
	"github.com/Zaibon/ircbot"
	"math/rand"
	"strings"
)

func FckBigx(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if m.Nick == "Bigx" {
		args := fmt.Sprintf("%s !quit", m.Channel)
		b.Out <- &ircbot.IrcMsg{
			Command: "PRIVMSG",
			Args:    []string{args},
		}
	}
}

func Greet(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
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

func Respond(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
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

const (
	gitlabUrl string = "gitlab.gigx.be/api/v3/projects/:id/repository/commits"
)

func GoDock(b *ircbot.IrcBot, m *ircbot.IrcMsg) {

	// if m.Args[0] == ":.doc" {

	// 	client := http.Client{}
	// 	req, err := http.NewRequest("GET", "http://godoc.org/?q=sql", nil)
	// 	if err != nil {
	// 		b.Error <- err
	// 	}
	// 	req.Header.Add("Accept", "text/plain")
	// 	resp, err := client.Do(req)
	// 	// resp, err := http.Get("http://godoc.org/?q=sql")
	// 	if err != nil {
	// 		b.Error <- err
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	m.Command = "PRIVMSG"
	// 	if resp.StatusCode != 200 {
	// 		m.Args = []string{":Pas de correpondance"}
	// 	} else {
	// 		//convert body to string
	// 		// buf := new(bytes.Buffer)
	// 		// buf.ReadFrom(resp.Body)
	// 		// s := buf.String()
	// 		scanner := bufio.NewScanner(resp.Body)
	// 		scanner.Split(bufio.ScanWords)

	// 		var documentation string
	// 		for scanner.Scan() {
	// 			w := scanner.Text()
	// 			last := w
	// 			if w == m.Args[1] {
	// 				documentation += last
	// 				scanner.Split(bufio.ScanBytes)
	// 				by := scanner.Bytes()
	// 				for by != "\n" {
	// 					documentation += by
	// 				}
	// 				break
	// 			}
	// 		}

	// 		m.Args = []string{fmt.Sprintf(":%s", s)}
	// 	}

	// 	b.Out <- m
	// }
}
