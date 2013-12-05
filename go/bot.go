// Package modules provides ...
package main

import (
	"./ircbot"
	"fmt"
	"time"
)

//=============================================================================
// main
//=============================================================================

func main() {

	//. define bots
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

	//. put bots in base
	bots := []*ircbot.IRCBot{bot1, bot2}
	base := ircbot.NewBase()
	base.AddBots(bots)

	//. launch the bots
	//base.Launch() //. do not have pipe to interactive with IRC messages
	base.Capture()

	//. do other stuffs
	for _, bot := range bots {
		channels := bot.Channels()
		for i := range channels {
			bot.Say(channels[i], fmt.Sprintf("Hello %s", channels[i]))
			bot.Action(channels[i], "shake his body")
			bot.Action(channels[i], "唱了一首歌.")
			bot.Debug()
		}
	}

	//. using pipe to interactive with IRC messages
	go func() {
		pipe := bot1.GetPipe()
		for msg := range pipe {
			if msg.IsPRIVMSG() {
				//fmt.Println(msg.Raw())
				fmt.Sprintln(msg.Raw())
			}
		}
	}()

	//. debug
	go func() {
		for {
			for _, bot := range bots {
				bot.Debug()
			}
			time.Sleep(time.Duration(time.Second * 1))
		}
	}()

	// wait to bots exist
	base.Wait()
}
