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
	pipe           BotPipe
	reader         *bufio.Reader
	writer         *textproto.Writer
	noises         chan string
	modules        BotModules
	logger         *log.Logger
	wg             *sync.WaitGroup
}

type BotModules map[string]reflect.Value
type BotPipe chan *Msg

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
		noises:         make(chan string, 1000),
		modules:        make(BotModules),
		conn:           nil,
		keepConnection: false,
		pipeMessage:    false,
		pipeBuffer:     list.New(),
		pipe:           make(BotPipe),
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

func (bot *IRCBot) Writef(format string, args ...interface{}) {
	bot.writer.PrintfLine(format, args...)
}

func (bot *IRCBot) Identify() {
	bot.Writef("USER %s %s %s :%s",
		bot.nickname, bot.hostname, bot.servername, bot.realname)
	bot.Writef("NICK %s", bot.nickname)
}

func (bot *IRCBot) JoinDefault() {
	for i := range bot.channels {
		bot.Writef("JOIN %s", bot.channels[i])
	}
}

func (bot *IRCBot) Pong(server string) {
	bot.Writef(fmt.Sprintf("PONG %s", server))
	bot.Log(">> PONG !")
}

func (bot *IRCBot) Reply(target, message string) {
	bot.Writef("PRIVMSG %s :%s", target, message)
}

func (bot *IRCBot) Notice(target, message string) {
	bot.Writef("NOTICE %s :%s", target, message)
}

func (bot *IRCBot) Say(channel, message string) {
	bot.noises <- fmt.Sprintf("PRIVMSG %s :%s", channel, message)
	//bot.Writef("PRIVMSG %s :%s", channel, message)
}

// see: http://www.irchelp.org/irchelp/rfc/ctcpspec.html
func (bot *IRCBot) Action(channel, message string) {
	bot.noises <- fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", channel, message)
	//bot.Writef("PRIVMSG %s :ACTION %s", channel, message)
}

func (bot *IRCBot) ReadLine() (string, error) {
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

func (bot *IRCBot) Listen() {
	for bot.conn != nil {
		msg, err := bot.ReadLine()
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
			if bot.pipeMessage {
				bot.pipeBuffer.PushBack(bot.toMsg(msg))
			}
			bot.Process(bot.toMsg(msg))
		}
		runtime.Gosched()
	}
}

// see: http://golang.org/pkg/reflect/#Value.Call
func (bot *IRCBot) Process(msg *Msg) {
	for _, mod := range bot.modules {
		botv := reflect.ValueOf(bot)
		msgv := reflect.ValueOf(msg)
		mod.Call([]reflect.Value{botv, msgv})
	}
}

func (bot *IRCBot) MakeNoise() {
	for bot.conn != nil {
		msg := <-bot.noises
		time.Sleep(time.Duration(time.Second * 3))
		bot.Writef(msg)
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

func (bot *IRCBot) Pipe() {
	for bot.conn != nil {
		if e := bot.pipeBuffer.Front(); e != nil {
			msg := bot.pipeBuffer.Remove(e).(*Msg)
			bot.pipe <- msg
		}
		runtime.Gosched()
	}
}

func (bot *IRCBot) GetPipe() BotPipe {
	return bot.pipe
}

func (bot *IRCBot) Link() {
	bot.Identify()
	bot.JoinDefault()
}

func (bot *IRCBot) Launch() {
	bot.RegisterModules(extmod.Functions)
	if bot.MustConnect() {
		go bot.Listen()
		go bot.MakeNoise()
		bot.Link()
	}
}

// create  a channel for user to recive the IRC messages
func (bot *IRCBot) Capture() {
	// init pipe
	bot.pipeMessage = true

	// other stuffs
	bot.RegisterModules(extmod.Functions)
	if bot.MustConnect() {
		go bot.Listen()
		go bot.MakeNoise()
		go bot.Pipe()
		bot.Link()
	}
}

func (bot *IRCBot) Debug() {
	fmt.Printf("%v Bot: %v, NumGoroutine: %v, pipeBuffer: %v, noises: %v\n",
		time.Now().Format(time.RFC3339),
		bot.nickname,
		runtime.NumGoroutine(),
		bot.pipeBuffer.Len(),
		len(bot.noises),
	)
}
