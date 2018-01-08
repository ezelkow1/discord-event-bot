package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	//"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//Configuration for bot
type Configuration struct {
	Token            string
	EventURL         string
	BroadcastChannel string
}

var (
	config     = Configuration{}
	configfile string
	embedColor = 0x00ff00
)

func init() {
	flag.StringVar(&configfile, "c", "", "Configuration file location")
	flag.Parse()

	if configfile == "" {
		fmt.Println("No config file entered")
		os.Exit(1)
	}

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		fmt.Println("Configfile does not exist, you should make one")
		os.Exit(2)
	}

	fileh, _ := os.Open(configfile)
	decoder := json.NewDecoder(fileh)
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(3)
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for message events
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "keys go in my piehole")
	SendEmbed(s, config.BroadcastChannel, "", "I iz here", "Eventbot has arrived, servicing all your scheduling needs")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only allow messages in either DM or broadcast channel
	dmchan, _ := s.UserChannelCreate(m.Author.ID)
	if (m.ChannelID != config.BroadcastChannel) && (m.ChannelID != dmchan.ID) {
		return
	}

	if strings.HasPrefix(m.Content, "!help") == true {
		PrintHelp(s, m)
	}

	if strings.HasPrefix(m.Content, "!schedule") == true {
		printSchedule(s, m)
	}
}

//PrintHelp will print out the help dialog
func PrintHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	var buffer bytes.Buffer
	buffer.WriteString("!schedule - prints out the upcoming schedule\n")
	buffer.WriteString("!next - prints the next upcoming game\n")
	SendEmbed(s, m.ChannelID, "", "Available Commands", buffer.String())
}

func printSchedule(s *discordgo.Session, m *discordgo.MessageCreate) {
	resp, _ := http.Get(config.EventURL)
	//bytes, _ := ioutil.ReadAll(resp.Body)

	root, err := html.Parse(resp.Body)
	if err != nil {
		// handle error
	}
	// Search for the title

	/* 	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent, "class") == "upcoming"
		}
		return false
	} */
	// grab all articles and print them

	allResults := scrape.FindAllNested(root, scrape.ByClass("upcoming_event"))
	var gametimes, gamenamesS []string
	//articles := scrape.FindAll(root, matcher)
	for i, article := range allResults {
		times := scrape.FindAll(article, scrape.ByClass("upcoming_event_timestamp"))
		for k, thisart := range times {
			gametimes = append(gametimes, scrape.Text(thisart))
			fmt.Printf("%d %2d %s\n", i, k, scrape.Text(thisart))
		}
	}

	for i, article := range allResults {
		gamenames := scrape.FindAll(article, scrape.ByTag(atom.A))
		for k, thisart := range gamenames {
			gamenamesS = append(gamenamesS, scrape.Text(thisart))
			fmt.Printf("%d %2d %s\n", i, k, scrape.Text(thisart))
		}
	}

	fmt.Println(gametimes, gamenamesS)

	//SendEmbed(s, config.BroadcastChannel, "", "Current Schedule", game)
	resp.Body.Close()
}
