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
	Send(format string, args ...interface{})
	Log(format string, v ...interface{})
	Quit()
}

type Msg interface {
	Raw() string
	Tags() []string
	Trim() string
	Parsemp(r *regexp.Regexp) map[string]string
	Parsese(r *regexp.Regexp) []string
	//ParsePRIVMSG() map[string]string
	IsPRIVMSG() bool
	IsSERVMSG() bool
	IsPINGMSG() bool
	IsUNKNMSG() bool
	GetPRIVMSG() PRIVMSG
}

type BotModule func(bot Bot, msg Msg)

type PRIVMSG struct {
	Nick    string
	User    string
	Host    string
	To      string
	Message string
}

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModulePong(bot Bot, msg Msg) {
	if msg.IsPINGMSG() {
		if msg.IsPRIVMSG() {
			//. from user
			pm := msg.GetPRIVMSG()
			timestamp := strings.Fields(strings.Trim(pm.Message, "\001"))[1]
			msg := fmt.Sprintf("\001PING %s\001", timestamp)
			bot.Notice(pm.Nick, msg)
		} else {
			//. from server
			bot.Pong(strings.Split(msg.Raw(), ":")[1])
		}
	}
}
