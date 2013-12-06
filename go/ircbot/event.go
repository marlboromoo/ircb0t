// Package ircboot
package ircbot

//=============================================================================
// type and variable
//=============================================================================

type MsgEvent struct {
	bot *IRCBot
	raw string
}

//=============================================================================
// methods
//=============================================================================

func NewMsgEvent(bot *IRCBot, raw string) *MsgEvent {
	return &MsgEvent{
		bot: bot,
		raw: raw,
	}
}

func (event *MsgEvent) Run() {
	event.bot.Send(event.raw)
}
