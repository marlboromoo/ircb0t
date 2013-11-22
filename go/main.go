package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"time"
)

//=============================================================================
// type
//=============================================================================

type IRCbot struct {
	//. IRC server
	address string //server:port
	//. Client infos
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
	modules []BotModule
}

type BotModule func(bot *IRCbot, msg string)

//=============================================================================
// methods
//=============================================================================

func NewBot(address, nickname, username, realname string,
	channels []string) *IRCbot {
	return &IRCbot{
		address:    address,
		nickname:   nickname,
		username:   username,
		hostname:   "hostname",
		servername: "servername",
		realname:   realname,
		channels:   channels,
		noises:     make(chan string, 1000),
		modules:    []BotModule{ModuleFoo},
	}
}

func (bot *IRCbot) Connect() {
	conn, err := net.Dial("tcp", bot.address)
	if err != nil {
		fmt.Printf("Fail to connect to IRC server!\n")
	}
	fmt.Printf("Connect to IRC server.\n")
	bot.conn = conn
	bot.reader = bufio.NewReader(bot.conn)
	bot.writer = textproto.NewWriter(bufio.NewWriter(bot.conn))
}

func (bot *IRCbot) Disconnect() {
	if bot.conn != nil {
		bot.conn.Close()
	}
}

func (bot *IRCbot) Writef(format string, args ...interface{}) {
	bot.writer.PrintfLine(format, args...)
}

func (bot *IRCbot) Identify() {
	bot.Writef("USER %s %s %s :%s",
		bot.nickname, bot.hostname, bot.servername, bot.realname)
	bot.Writef("NICK %s", bot.nickname)
}

func (bot *IRCbot) JoinDefault() {
	for i := range bot.channels {
		bot.Writef("JOIN %s", bot.channels[i])
	}
}

func (bot *IRCbot) Say(channel, message string) {
	bot.noises <- fmt.Sprintf("PRIVMSG %s :%s", channel, message)
	//bot.Writef("PRIVMSG %s :%s", channel, message)
}

func (bot *IRCbot) Action(channel, message string) {
	bot.noises <- fmt.Sprintf("PRIVMSG %s :ACTION %s", channel, message)
	//bot.Writef("PRIVMSG %s :ACTION %s", channel, message)
}

func (bot *IRCbot) ReadLine() string {
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

func (bot *IRCbot) Listen() {
	for {
		msg := bot.ReadLine()
		fmt.Printf("%s\n", msg)
		bot.Process(msg)
	}
}

func (bot *IRCbot) Process(msg string) {
	for _, module := range bot.modules {
		module(bot, msg)
	}
}

func (bot *IRCbot) MakeNoise() {
	for {
		msg := <-bot.noises
		time.Sleep(time.Duration(time.Second * 3))
		bot.Writef(msg)
	}
}

//=============================================================================
// funtions (module for IRCbot)
//=============================================================================

func ModuleFoo(bot *IRCbot, msg string) {
	for i := range bot.channels {
		bot.Say(bot.channels[i], msg)
		break
	}
}

//=============================================================================
// main
//=============================================================================

func main() {
	bot := NewBot(
		"10.10.5.32:6667",
		"r0b0t",
		"robot",
		"robot",
		[]string{"#foo", "#bar"},
	)

	bot.Connect()
	defer bot.Disconnect()
	bot.Identify()
	bot.JoinDefault()

        //. process messages
	go bot.Listen()
	go bot.MakeNoise()

	for i := range bot.channels {
		bot.Say(bot.channels[i], fmt.Sprintf("Hello %s", bot.channels[i]))
		bot.Action(bot.channels[i], "shake his body")
	}

	// run forever
	select {}

}
