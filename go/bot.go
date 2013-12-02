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
		"r0b0t01",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bot2 := ircbot.NewBot(
		"irc.lab:6667",
		"r0b0t02",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bots := []*ircbot.IRCBot{bot1, bot2}

	base := ircbot.NewBase()
	base.AddBots(bots)
	base.Launch()

	//defer bot1.Disconnect()
	//defer bot2.Disconnect()

	for _, bot := range bots {
		for i := range bot.Channels {
			bot.Say(bot.Channels[i], fmt.Sprintf("Hello %s", bot.Channels[i]))
			bot.Action(bot.Channels[i], "shake his body")
			bot.Action(bot.Channels[i], "唱了一首歌.")
		}
	}

	// wait to bots exist
	time.Sleep(time.Duration(time.Second * 3))
	bot1.Disconnect()
	bot2.Disconnect()
	base.Wait()

}

