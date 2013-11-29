// Package modules provides ...
package ircbot

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"reflect"
	"time"
)

import extmod "./module"

//=============================================================================
// type
//=============================================================================

type IRCBot struct {
	//. IRC server
	address string "server:port"
	//. Client infos
	nickname   string
	username   string
	hostname   string
	servername string
	realname   string
	Channels   []string
	//. MISC
	conn    net.Conn
	reader  *bufio.Reader
	writer  *textproto.Writer
	noises  chan string
	modules BotModules
}

type BotModules map[string]reflect.Value

//=============================================================================
// methods
//=============================================================================

func NewBot(address, nickname, username, realname string,
	Channels []string) *IRCBot {

	fmt.Printf("Initial IRCBot ...\n")
	bot := &IRCBot{
		address:    address,
		nickname:   nickname,
		username:   username,
		hostname:   "hostname",
		servername: "servername",
		realname:   realname,
		Channels:   Channels,
		noises:     make(chan string, 1000),
		modules:    make(BotModules),
	}
	bot.RegisterModules(extmod.Functions)
	return bot
}

func (bot *IRCBot) RegisterModule(modname string, mod reflect.Value) {
	modt := extmod.Types["BotModule"]
	// check module == module.BotModule
	if modt.ConvertibleTo(mod.Type()) {
		bot.modules["modname"] = mod
		fmt.Printf("++ inject module: %v(%v)\n", modname, mod.Type())
	} else {
		fmt.Printf("-- skip module: %v(%v)\n", modname, mod.Type())
	}
}

func (bot *IRCBot) RegisterModules(mods BotModules) {
	fmt.Printf("Register moduels ...\n")
	for modname, mod := range mods {
		bot.RegisterModule(modname, mod)
	}
}

func (bot *IRCBot) Connect() {
	conn, err := net.Dial("tcp", bot.address)
	if err != nil {
		fmt.Printf("Fail to connect to IRC server!\n")
	}
	fmt.Printf("Connect to IRC server.\n")
	bot.conn = conn
	bot.reader = bufio.NewReader(bot.conn)
	bot.writer = textproto.NewWriter(bufio.NewWriter(bot.conn))
}

func (bot *IRCBot) Disconnect() {
	if bot.conn != nil {
		bot.conn.Close()
	}
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
	for i := range bot.Channels {
		bot.Writef("JOIN %s", bot.Channels[i])
	}
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
	for {
		msg := bot.ReadLine()
		fmt.Printf("%s\n", msg)
		bot.Process(msg)
	}
}

// see: http://golang.org/pkg/reflect/#Value.Call
func (bot *IRCBot) Process(msg string) {
	for _, mod := range bot.modules {
		botv := reflect.ValueOf(bot)
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

func (bot *IRCBot) GetChannels() []string {
	return bot.Channels
}
