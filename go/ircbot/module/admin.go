// Package modules provides ...
package module

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleQuit(bot Bot, msg Msg) {
	if msg.IsPRIVMSG() {
		pm := msg.GetPRIVMSG()
		//. must recive private message from bot's owner
		if pm.Nick == bot.Owner() && pm.To == bot.Nickname() && pm.Message == ".quit" {
			bot.Quit()
		}
	}
}
