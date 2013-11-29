// Package modules provides ...
package module

import (
	"regexp"
	"fmt"
)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModulePong(bot Bot, msg string) {
	var serverPing = regexp.MustCompile(`^PING :(?P<server>.*)`)
	var userPing = regexp.MustCompile(
		`:(?P<who>.*) PRIVMSG (?P<target>.*) :\001PING (?P<timestamp>.*)\001`)
	for _, ping := range []*regexp.Regexp{serverPing, userPing} {
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
