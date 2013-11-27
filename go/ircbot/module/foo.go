// Package modules provides ...
package module

type bot interface {
    GetChannels() []string
    Say(channel, msg string)
}

type BotModule func(bot bot, msg string)

//=============================================================================
// funtions (module for IRCBot)
//=============================================================================

func ModuleFoo(bot bot, msg string) {
	for _, channel := range bot.GetChannels() {
		bot.Say(channel, msg)
		break
	}
}
