// Package modules provides ...
package module

import (
	"fmt"
	"regexp"
	"strings"
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
	Writef(format string, args ...interface{})
	Log(format string, v ...interface{})
	Disconnect()
}

type Msg interface {
	Raw() string
	Tags() []string
	Trim() string
	Parsemp(r *regexp.Regexp) map[string]string
	Parsese(r *regexp.Regexp) []string
	ParsePRIVMSG() map[string]string
	IsPRIVMSG() bool
	IsSERVRMSG() bool
	IsPINGMSG() bool
	IsUNKNMSG() bool
	GetPRIVMSG() string
}

type BotModule func(bot Bot, msg Msg)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleDebugMSG(bot Bot, msg Msg) {
	//fmt.Println(msg.Tags())
	//fmt.Println(msg.Raw())
}

func ModulePong(bot Bot, msg Msg) {
	if msg.IsPINGMSG() {
		if msg.IsPRIVMSG() {
			//. from user
			mp := msg.ParsePRIVMSG()
			timestamp := strings.Fields(strings.Trim(mp["message"], "\001"))[1]
			msg := fmt.Sprintf("\001PING %s\001", timestamp)
			bot.Notice(mp["nick"], msg)
		} else {
			//. from server
			bot.Pong(strings.Split(msg.Raw(), ":")[1])
		}
	}
}

func ModuleQuit(bot Bot, msg Msg) {
	if msg.IsPRIVMSG() {
		mp := msg.ParsePRIVMSG()
		//. must recive private message from bot's owner
		if mp["nick"] == bot.Owner() && mp["to"] == bot.Nickname() && mp["message"] == ".quit" {
			bot.Disconnect()
		}
	}
}
