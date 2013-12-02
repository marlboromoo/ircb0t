// Package modules provides ...
package module

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleIllegal1(bot Bot, msg string, mystuff string) {
	bot, msg, mystuff = nil, "", ""
}

func ModuleIllegal2(bot Bot) string {
	bot = nil
	return "foo."
}
