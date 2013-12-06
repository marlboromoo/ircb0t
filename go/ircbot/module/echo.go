// Package modules provides ...
package module

import (
	"strings"
)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleEcho(bot Bot, msg Msg) {
	if msg.IsPRIVMSG() {
		r := msg.ParsePRIVMSG()
		if strings.Fields(r["message"])[0] == ".echo" {
			bot.Say(r["to"], strings.TrimLeft(r["message"], ".echo  "))
		}
	}
}
