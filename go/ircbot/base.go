// Package ircboot
package ircbot

import (
	"sync"
)

//=============================================================================
// type and variable
//=============================================================================

type Base struct {
	wg   *sync.WaitGroup
	bots []*IRCBot
}

//=============================================================================
// methods
//=============================================================================

func NewBase() *Base {
	var wg sync.WaitGroup
	p := Base{
		wg:   &wg,
		bots: []*IRCBot{},
	}
	return &p
}

func (base *Base) AddBot(bot *IRCBot) {
	base.bots = append(base.bots, bot)
	base.wg.Add(1)
	bot.wg = base.wg
}

func (base *Base) AddBots(bots []*IRCBot) {
	for _, bot := range bots {
		base.AddBot(bot)
	}
}

func (base *Base) Wait() {
	base.wg.Wait()
}

func (base *Base) Launch() {
	for _, bot := range base.bots {
		go bot.Launch()
	}
}

func (base *Base) Capture() {
	for _, bot := range base.bots {
		go bot.Capture()
	}
}
