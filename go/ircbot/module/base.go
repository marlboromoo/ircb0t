// Package module
package module

//=============================================================================
// types
//=============================================================================

type BotModule struct {
    IsModule bool
    Name string
}

type Module interface {
    Process(bot Bot, msg string)
}

type Bot interface {
    GetChannels() []string
    Say(channel, msg string)
}

//type BotModule func(bot Bot, msg string)



