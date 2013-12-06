// Package modules provides ...
package main

import (
	"./ircbot"
	"fmt"
	"runtime"
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

	//. link bots with controller
	bot1.PipeOn() //. enable pipe to recive the IRC messages
	bot2.PipeOn()
	bots := []*ircbot.IRCBot{bot1, bot2}
	ctrler := ircbot.NewController()
	ctrler.LinkBots(bots)

	//. launch the bots
	ctrler.Launch()

	//. ensure the robot are connected
	for !ctrler.BotsAreConnected() {
		time.Sleep(time.Duration(time.Millisecond * 500))
	}

	//. do other stuffs
	for _, bot := range bots {
		channels := bot.Channels()
		for i := range channels {
			bot.Say(channels[i], fmt.Sprintf("Hello %s", channels[i]))
			bot.Action(channels[i], "向大家問好")
		}
	}

	//. using pipe to interactive with IRC messages
	pipe := bot1.GetPipe()
	go func() {
		for msg := range pipe {
			if msg.IsPRIVMSG() {
				//fmt.Println(msg.Raw())
				fmt.Sprintln(msg.Raw())
			}
			runtime.Gosched()
		}
	}()

	//. just for debug
	go func() {
		for {
			for _, bot := range bots {
				bot.Debug()
			}
			time.Sleep(time.Duration(time.Second * 1))
		}
	}()

	// wait to bots exist
	ctrler.Wait()
}
