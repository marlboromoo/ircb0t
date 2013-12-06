// Package ircboot
package ircbot

import (
	"sync"
)

//=============================================================================
// type and variable
//=============================================================================

type Controller struct {
	wg   *sync.WaitGroup
	bots []*IRCBot
}

//=============================================================================
// methods
//=============================================================================

func NewController() *Controller {
	var wg sync.WaitGroup
	p := Controller{
		wg:   &wg,
		bots: []*IRCBot{},
	}
	return &p
}

func (ctrler *Controller) LinkBot(bot *IRCBot) {
	ctrler.bots = append(ctrler.bots, bot)
	ctrler.wg.Add(1)
	bot.wg = ctrler.wg
}

func (ctrler *Controller) LinkBots(bots []*IRCBot) {
	for _, bot := range bots {
		ctrler.LinkBot(bot)
	}
}

func (ctrler *Controller) Wait() {
	ctrler.wg.Wait()
}

func (ctrler *Controller) Launch() {
	for _, bot := range ctrler.bots {
		go bot.Launch()
	}
}

func (ctrler *Controller) BotsAreConnected() bool {
	results := []bool{}
	for _, bot := range ctrler.bots {
		results = append(results, bot.IsConnected())
	}
	for _, r := range results {
		if !r {
			return false
		}
	}
	return true
}
