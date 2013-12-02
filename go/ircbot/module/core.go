// Package modules provides ...
package module

import (
	"regexp"
	"fmt"
	"../filter"
)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModulePong(bot Bot, msg string) {
	for _, ping := range []*regexp.Regexp{filter.ServerPing, filter.UserPing} {
		result := bot.ParseMsg(msg, ping)
		if len(result) >= 1 {
			if server, ok := result["server"]; ok {
				bot.Pong(server)
			}
			if who, ok := result["who"]; ok {
				msg:= fmt.Sprintf("\001PING %s\001", result["timestamp"])
				bot.Notice(who, msg)
				fmt.Println(msg)
			}
		}
	}
}
