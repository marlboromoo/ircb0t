// Package modules provides ...
package main

import (
    "fmt"
	"time"
    "./ircbot"
)

//=============================================================================
// main
//=============================================================================

func main() {

	bot1 := ircbot.NewBot(
		"irc.lab:6667",
		"marlboromoo",
		"r0b0t01",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bot2 := ircbot.NewBot(
		"irc.lab:6667",
		"marlboromoo",
		"r0b0t02",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bots := []*ircbot.IRCBot{bot1, bot2}

	base := ircbot.NewBase()
	base.AddBots(bots)
	base.Launch()

	for _, bot := range bots {
		channels := bot.Channels()
		for i := range channels {
			bot.Say(channels[i], fmt.Sprintf("Hello %s", channels[i]))
			bot.Action(channels[i], "shake his body")
			bot.Action(channels[i], "唱了一首歌.")
		}
	}

	// wait to bots exist
	time.Sleep(time.Duration(time.Second * 3))
	//bot1.Disconnect()
	//bot2.Disconnect()
	base.Wait()

}

