// Package ircboot
package ircbot

import (
	//"fmt"
	//"runtime"
	"time"
)

//=============================================================================
// type and variable
//=============================================================================
type Scheduler struct {
	highPriority chan Job
	lowPriority  chan Job
}

type Job interface {
	Run()
}

//=============================================================================
// methods
//=============================================================================

func NewScheduler() *Scheduler {
	return &Scheduler{
		highPriority: make(chan Job, 1000),
		lowPriority:  make(chan Job, 1000),
	}
}

func (schder *Scheduler) AddHigh(job Job) {
	schder.highPriority <- job
}

func (schder *Scheduler) AddLow(job Job) {
	schder.lowPriority <- job
}

func (schder *Scheduler) Run() {
	for {
		//fmt.Println("*** scheduler")
		switch {
		case len(schder.highPriority) >= 1:
			job := <-schder.highPriority
			job.Run()
		case len(schder.highPriority) == 0 && len(schder.lowPriority) >= 1:
			job := <-schder.lowPriority
			job.Run()
		}
		time.Sleep(time.Duration(time.Millisecond * 500))
	}
}
