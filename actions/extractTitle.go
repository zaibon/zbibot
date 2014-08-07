package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/go.net/html"

	"github.com/Zaibon/ircbot"
)

type TitleExtract struct {
	name string
}

func (u *TitleExtract) Command() []string {
	return []string{
		"PRIVMSG",
	}
}

func (u *TitleExtract) Usage() string {
	return ""
}

func (u *TitleExtract) Do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	do(b, m)
}

func do(b *ircbot.IrcBot, m *ircbot.IrcMsg) {
	m.Args[0] = strings.TrimPrefix(m.Args[0], ":")
	for _, word := range m.Args {
		fmt.Println(word)
		u, err := url.Parse(word)
		if err != nil {
			fmt.Println("err parse url: ", err)
			continue
		}

		go func() {

			resp, err := http.Get(u.String())
			if err != nil {
				fmt.Println("err get url: ", err)
				return
			}

			contentType := resp.Header.Get("Content-Type")

			switch {
			case strings.Contains(contentType, "text/html"):
				doc, err := html.Parse(resp.Body)
				resp.Body.Close()
				if err != nil {
					fmt.Println("err Parse page : ", err)
					return
				}
				title := extractTitle(doc)
				fmt.Println("extract title : ", title)
				if title != "" {
					b.Say(m.Channel, title)
				}
			default:
				fmt.Println("mime not supported")
			}
		}()
	}
}

func extractTitle(n *html.Node) string {
	var curr *html.Node
	var title string
	curr = n
	for curr != nil {
		if curr.Data == "title" {
			if curr.FirstChild != nil {
				title = curr.FirstChild.Data
				break
			}
		}
		if curr.FirstChild != nil {
			curr = curr.FirstChild
		} else {
			curr = curr.NextSibling
		}
	}
	return title
}
