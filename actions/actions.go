package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Zaibon/ircbot"
	"math/rand"
	"os"
	"strings"
	"time"
)

func KickAlex(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if m.Nick == "banzounet" {
		time.Sleep(10 * time.Second)
		b.Out <- &ircbot.IrcMsg{
			Command: "KICK",
			Channel: m.Channel,
			Args:    []string{m.Nick},
		}
	}
}

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

func InteractiveCommands(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	if !strings.HasPrefix(m.Args[0], ":.") {
		return
	}

	//read file that contains link
	f, err := os.Open("infoUrl.json")
	if err != nil {
		return
	}
	defer f.Close()
	//unmarshall into map
	links := make(map[string][]string)
	dec := json.NewDecoder(f)
	dec.Decode(&links)

	//parse irc command
	command := strings.TrimPrefix(m.Args[0], ":.")
	switch command {
	case "link":
		if len(m.Args) < 2 {
			//no link spécified
			//display all available
			helpLinks := ""
			for k, _ := range links {
				helpLinks += k + " "
			}
			b.Say(m.Channel, helpLinks)
			return
		}

		for _, v := range links[m.Args[1]] {
			b.Say(m.Channel, v)
		}
	case "ticker":
		Ticker(b, m)

	case "last":
		LastBlock(b, m)

	case "status":
		Status(b, m)

	case "u", "user":
		User(b, m)
	case "help":
		b.Say(m.Channel, ".link .ticker .last <nbr> .status .u <user> .user <user> .stats")
	case "stats":
		OverallStats(b, m)
	}
}

// func Info(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
// 	if !strings.HasPrefix(m.Args[0], ":.") {
// 		return
// 	}

// 	command := strings.TrimPrefix(m.Args[0], ":.")
// 	switch command {
// 	case "link":
// 		if len(m.Args) < 2 {
// 			return
// 		}
// 		switch m.Args[1] {
// 		case "stratum":
// 			b.Say(m.Channel, "http://mining.bitcoin.cz/stratum-mining")
// 		case "bter":
// 			b.Say(m.Channel, "https://bter.com/")
// 		}

// 	case "pool":
// 		b.Say(m.Channel, ":pool    : http://laminerie.eu")
// 		b.Say(m.Channel, ":stratum : stratum+tcp://laminerieu.eu:3333")
// 		b.Say(m.Channel, ":getwork : http://laminerie.eu:3335")
// 	}
// }

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
