// Package modules provides ...
package ircbot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"reflect"
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
	conn    net.Conn
	reader  *bufio.Reader
	writer  *textproto.Writer
	noises  chan string
	modules BotModules
	logger  *log.Logger
	wg      *sync.WaitGroup
}

type BotModules map[string]reflect.Value

//=============================================================================
// methods
//=============================================================================

func NewBot(address, owner, nickname, username, realname string,
	channels []string) *IRCBot {

	bot := IRCBot{
		owner:      owner,
		address:    address,
		nickname:   nickname,
		username:   username,
		hostname:   "hostname",
		servername: "servername",
		realname:   realname,
		channels:   channels,
		noises:     make(chan string, 1000),
		modules:    make(BotModules),
	}

	//. logger
	f, _ := os.OpenFile(
		fmt.Sprintf("/tmp/%s.log", bot.nickname),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	bot.logger = log.New(f, bot.nickname, log.Ldate|log.Lmicroseconds)

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

func (bot *IRCBot) Connect() {
	conn, err := net.Dial("tcp", bot.address)
	if err != nil {
		bot.Log("!! Fail to connect to IRC server!\n")
	}
	bot.Log(">> Connect to IRC server.\n")
	bot.conn = conn
	bot.reader = bufio.NewReader(bot.conn)
	bot.writer = textproto.NewWriter(bufio.NewWriter(bot.conn))
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

func (bot *IRCBot) ReadLine() string {
	var line_ []byte
	for true {
		line, isPrefix, _ := bot.reader.ReadLine()
		line_ = append(line_, line...)
		if !isPrefix {
			break
		}
	}
	return string(line_)
}

func (bot *IRCBot) Listen() {
	for bot.conn != nil {
		msg := bot.ReadLine()
		bot.Log("<< %s\n", msg)
		bot.Process(msg)
	}
}

// see: http://golang.org/pkg/reflect/#Value.Call
func (bot *IRCBot) Process(msg string) {
	for _, mod := range bot.modules {
		botv := reflect.ValueOf(bot)
		msg := NewMsg(
			strings.Split(bot.address, ":")[0],
			strings.Split(bot.address, ":")[1],
			msg)
		msgv := reflect.ValueOf(msg)
		mod.Call([]reflect.Value{botv, msgv})
	}
}

func (bot *IRCBot) MakeNoise() {
	for {
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

func (bot *IRCBot) Launch() {
	bot.RegisterModules(extmod.Functions)
	bot.Connect()

	//. process messages
	go bot.Listen()

	//. say hello
	bot.Identify()
	bot.JoinDefault()

	//. say something
	go bot.MakeNoise()
}
