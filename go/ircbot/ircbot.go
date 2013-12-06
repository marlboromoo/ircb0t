// Package modules provides ...
package ircbot

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

import extmod "./module"

//=============================================================================
// type and variable
//=============================================================================

type IRCBot struct {
	//. IRC server
	address string "server:port"
	//. Client infos
	owner      string
	nickname   string
	username   string
	hostname   string
	servername string
	realname   string
	channels   []string
	//. MISC
	conn           net.Conn
	keepConnection bool
	pipeMessage    bool
	pipeBuffer     *list.List
	pipe           Botpipe
	reader         *bufio.Reader
	writer         *textproto.Writer
	//noises         chan string
	noises    *list.List
	modules   BotModules
	logger    *log.Logger
	wg        *sync.WaitGroup
	scheduler *Scheduler
}

type BotModules map[string]reflect.Value
type Botpipe chan *Msg

//=============================================================================
// methods
//=============================================================================

func NewBot(address, owner, nickname, username, realname string,
	channels []string) *IRCBot {

	bot := IRCBot{
		owner:          owner,
		address:        address,
		nickname:       nickname,
		username:       username,
		hostname:       "hostname",
		servername:     "servername",
		realname:       realname,
		channels:       channels,
		noises:         list.New(),
		modules:        make(BotModules),
		conn:           nil,
		keepConnection: false,
		pipeMessage:    false,
		pipeBuffer:     list.New(),
		pipe:           make(Botpipe),
		scheduler:      NewScheduler(),
	}

	//. logger
	f, _ := os.OpenFile(
		fmt.Sprintf("/tmp/%s.log", bot.nickname),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	bot.logger = log.New(f, bot.nickname+" ", log.Ldate|log.Lmicroseconds)

	//. bot pointer
	return &bot
}

func (bot *IRCBot) Log(format string, v ...interface{}) {
	bot.logger.Printf(format, v...)
}

func (bot *IRCBot) RegisterModule(modname string, mod reflect.Value) {
	modt := extmod.Types["BotModule"]
	// check module == module.BotModule
	if modt.ConvertibleTo(mod.Type()) {
		bot.modules[modname] = mod
		bot.Log("++ inject module: %v(%v)\n", modname, mod.Type())
	} else {
		bot.Log("-- skip module: %v(%v)\n", modname, mod.Type())
	}
}

func (bot *IRCBot) RegisterModules(mods BotModules) {
	bot.Log("?? Register moduels ...\n")
	for modname, mod := range mods {
		bot.RegisterModule(modname, mod)
	}
	bot.Log("** Register %v modules, %v modules fail to load.\n",
		len(bot.modules), len(extmod.Functions)-len(bot.modules))
}

func (bot *IRCBot) Connect() bool {
	conn, err := net.DialTimeout("tcp", bot.address, time.Duration(time.Second*10))
	if err != nil {
		bot.Log("!! Fail to connect to IRC server!\n")
		return false
	}
	bot.Log(">> Connect to IRC server.\n")
	bot.conn = conn
	bot.reader = bufio.NewReader(bot.conn)
	bot.writer = textproto.NewWriter(bufio.NewWriter(bot.conn))
	return true
}

func (bot *IRCBot) MustConnect() bool {
	bot.keepConnection = true
	for !bot.Connect() {
		bot.Log(">> Retrying connect to IRC server ...")
		time.Sleep(time.Duration(time.Second * 1))
	}
	return true
}

func (bot *IRCBot) Disconnect() {
	if bot.conn != nil {
		bot.conn.Close()
		bot.conn = nil
	}
	if bot.wg != nil {
		bot.wg.Done()
	}
	bot.Log(">> Disconnect from IRC server !!\n")
}

func (bot *IRCBot) IsConnected() bool {
	if bot.conn != nil && bot.reader != nil && bot.writer != nil {
		return true
	}
	return false
}

func (bot *IRCBot) Send(format string, args ...interface{}) {
	if bot.conn != nil {
		bot.writer.PrintfLine(format, args...)
	}
}

func (bot *IRCBot) Identify() {
	bot.Send("USER %s %s %s :%s",
		bot.nickname, bot.hostname, bot.servername, bot.realname)
	bot.Send("NICK %s", bot.nickname)
}

func (bot *IRCBot) JoinDefault() {
	for i := range bot.channels {
		bot.Send("JOIN %s", bot.channels[i])
	}
}

func (bot *IRCBot) Pong(server string) {
	//bot.Send(fmt.Sprintf("PONG %s", server))
	//bot.Log(">> PONG !")
	event := NewMsgEvent(bot, fmt.Sprintf("PONG %s", server))
	bot.scheduler.AddHigh(event)
}

func (bot *IRCBot) Reply(target, message string) {
	bot.Send("PRIVMSG %s :%s", target, message)
}

func (bot *IRCBot) Notice(target, message string) {
	//bot.Send("NOTICE %s :%s", target, message)
	event := NewMsgEvent(bot, fmt.Sprintf("NOTICE %s :%s", target, message))
	bot.scheduler.AddHigh(event)
}

func (bot *IRCBot) Say(channel, message string) {
	//msg := fmt.Sprintf("PRIVMSG %s :%s", channel, message)
	//bot.noises.PushFront(msg)
	//bot.Send(msg)
	event := NewMsgEvent(bot, fmt.Sprintf("PRIVMSG %s :%s", channel, message))
	bot.scheduler.AddLow(event)
}

// see: http://www.irchelp.org/irchelp/rfc/ctcpspec.html
func (bot *IRCBot) Action(channel, message string) {
	//msg := fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", channel, message)
	//bot.noises.PushFront(msg)
	//bot.Send(msg)
	event := NewMsgEvent(bot, fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", channel, message))
	bot.scheduler.AddLow(event)
}

func (bot *IRCBot) readLine() (string, error) {
	var line_ []byte
	for true {
		line, isPrefix, err := bot.reader.ReadLine()
		if err != nil {
			return "", err
		}
		line_ = append(line_, line...)
		if !isPrefix {
			break
		}
	}
	return string(line_), nil
}

func (bot *IRCBot) toMsg(msg string) *Msg {
	return NewMsg(
		strings.Split(bot.address, ":")[0],
		strings.Split(bot.address, ":")[1],
		msg)
}

func (bot *IRCBot) listen() {
	for bot.conn != nil {
		msg, err := bot.readLine()
		if err != nil {
			bot.Log("!! Lost connection !!")
			if bot.keepConnection {
				if bot.MustConnect() {
					bot.Link()
				}
			} else {
				bot.Disconnect()
			}
		}
		if err == nil && len(msg) >= 1 {
			bot.Log("<< %s\n", msg)
			go bot.process(bot.toMsg(msg))
			if bot.pipeMessage {
				bot.pipeBuffer.PushBack(bot.toMsg(msg))
			}
		}
	}
}

// see: http://golang.org/pkg/reflect/#Value.Call
func (bot *IRCBot) process(msg *Msg) {
	for _, mod := range bot.modules {
		botv := reflect.ValueOf(bot)
		msgv := reflect.ValueOf(msg)
		go mod.Call([]reflect.Value{botv, msgv})
	}
}

func (bot *IRCBot) makeNoise() {
	for bot.conn != nil {
		//bot.Log("** makeNoise")
		for bot.noises.Len() >= 1 {
			if e := bot.noises.Front(); e != nil {
				msg := bot.noises.Remove(e)
				bot.Send(msg.(string))
				time.Sleep(time.Duration(time.Millisecond * 500))
			}
		}
		time.Sleep(time.Duration(time.Millisecond * 500))
	}
}

func (bot IRCBot) Owner() string {
	return bot.owner
}

func (bot *IRCBot) Nickname() string {
	return bot.nickname
}

func (bot *IRCBot) Channels() []string {
	return bot.channels
}

func (bot *IRCBot) makePipe() {
	for bot.conn != nil {
		//bot.Log("** makePipe")
		for bot.pipeBuffer.Len() >= 1 {
			if e := bot.pipeBuffer.Front(); e != nil {
				msg := bot.pipeBuffer.Remove(e).(*Msg)
				bot.pipe <- msg
			}

		}
		time.Sleep(time.Duration(time.Millisecond * 500))
	}
}

func (bot *IRCBot) GetPipe() Botpipe {
	return bot.pipe
}

func (bot *IRCBot) PipeOn() {
	bot.pipeMessage = true
}

func (bot *IRCBot) PipeOff() {
	bot.pipeMessage = false
}

func (bot *IRCBot) Link() {
	bot.Identify()
	bot.JoinDefault()
}

func (bot *IRCBot) Launch() {
	bot.RegisterModules(extmod.Functions)
	if bot.MustConnect() {
		go bot.listen()
		//go bot.makeNoise()
		go bot.makePipe()
		go bot.scheduler.Run()
		bot.Link()
	}
}

func (bot *IRCBot) Debug() {
	fmt.Printf("%v Bot: %v, NumGoroutine: %v, pBuffer: %v, noises: %v, eHigh: %v, eLow: %v\n",
		time.Now().Format(time.RFC3339),
		bot.nickname,
		runtime.NumGoroutine(),
		bot.pipeBuffer.Len(),
		bot.noises.Len(),
		len(bot.scheduler.highPriority),
		len(bot.scheduler.lowPriority),
	)

}
