// Package modules provides ...
package main

import (
    "fmt"
    "./ircbot"
)

//=============================================================================
// main
//=============================================================================

func main() {
	bot := ircbot.NewBot(
		"10.10.5.32:6667",
		"r0b0t",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bot.Connect()
	defer bot.Disconnect()
	bot.Identify()
	bot.JoinDefault()

        //. process messages
	go bot.Listen()
	go bot.MakeNoise()

	for i := range bot.Channels {
		bot.Say(bot.Channels[i], fmt.Sprintf("Hello %s", bot.Channels[i]))
		bot.Action(bot.Channels[i], "shake his body")
		bot.Action(bot.Channels[i], "唱了一首歌.")
	}

	// run forever
	select {}

}

