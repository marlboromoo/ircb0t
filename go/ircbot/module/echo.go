// Package modules provides ...
package module

//=============================================================================
// types
//=============================================================================

type Bot interface {
    GetChannels() []string
    Say(channel, msg string)
}

type BotModule func(bot Bot, msg string)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleEcho(bot Bot, msg string) {
	for _, channel := range bot.GetChannels() {
		bot.Say(channel, msg)
	}
}
