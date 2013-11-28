// Package module
package module

//=============================================================================
// types
//=============================================================================

type Module interface {
    Process(bot Bot, msg string)
}

type Bot interface {
    GetChannels() []string
    Say(channel, msg string)
}

type BotModule func(bot Bot, msg string)



