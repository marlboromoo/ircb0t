// Package modules provides ...
package module

import (
	"../filter"
	"fmt"
	"regexp"
	//"strings"
)

//=============================================================================
// types
//=============================================================================

type Module interface {
    Process(bot Bot, msg string)
}

type Bot interface {
    Owner() string
	Nickname() string
	Channels() []string
    Say(channel, msg string)
	Reply(target, msg string)
	Notice(target, msg string)
	Pong(server string)
	Writef(format string, args ... interface{})
	ParseMsg(msg string, r *regexp.Regexp) map[string]string
	ParseWho(who string) map[string]string
	Log(format string, v ...interface{})
	Disconnect()
}

type BotModule func(bot Bot, msg string)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

//func ModulePong(bot Bot, msg string) {
//	for _, ping := range []*regexp.Regexp{filter.ServerPing, filter.UserPing} {
//		result := bot.ParseMsg(msg, ping)
//		if len(result) >= 1 {
//			if server, ok := result["server"]; ok {
//				bot.Pong(server)
//			}
//			if who, ok := result["who"]; ok {
//				msg := fmt.Sprintf("\001PING %s\001", result["timestamp"])
//				bot.Notice(who, msg)
//			}
//		}
//	}
//}

func ModuleQuit(bot Bot, msg string) {
	result := bot.ParseMsg(msg, filter.Quit)
	
	if len(result) >=1 {
		if who, ok := result["who"]; ok {
			fmt.Printf("%v\n", bot.ParseWho(who)["nick"])
			who = bot.ParseWho(who)["nick"]

			if result["target"] == bot.Nickname() && who == bot.Owner() {
				bot.Disconnect()
			}
		}
	}
}
