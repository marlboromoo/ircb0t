// Package module
package module

import (
	"regexp"
)

//=============================================================================
// types
//=============================================================================

type Module interface {
    Process(bot Bot, msg string)
}

type Bot interface {
    GetChannels() []string
    Say(channel, msg string)
	Reply(target, msg string)
	Notice(target, msg string)
	Pong(server string)
	Writef(format string, args ... interface{})
	ParseMsg(msg string, r *regexp.Regexp) map[string]string
}

type BotModule func(bot Bot, msg string)

