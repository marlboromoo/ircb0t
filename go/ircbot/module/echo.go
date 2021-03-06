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
		pm := msg.GetPRIVMSG()
		if strings.Fields(pm.Message)[0] == ".echo" {
			bot.Say(pm.To, strings.TrimLeft(pm.Message, ".echo  "))
		}
	}
}
