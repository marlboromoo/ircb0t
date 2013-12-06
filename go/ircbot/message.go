// Package msg provides ...
package ircbot

import (
	"regexp"
	"strings"
)

//=============================================================================
// type and variable
//=============================================================================

type Msg struct {
	host string
	port string
	raw  string
	tags []string
}

const (
	PRIVMSG = "PRIVMSG"
	PINGMSG = "PING"
	SERVMSG = "SERVER"
	UNKNMSG = "UNKNOWN"
)

var (
	rePINGMSGs = regexp.MustCompile(`^PING :(?P<server>.*)`)
	rePINGMSGu = regexp.MustCompile(
		`:(?P<from>.*) PRIVMSG (?P<to>.*) :\001PING (?P<timestamp>.*)\001`)
	rePRIVMSG = regexp.MustCompile(
		`:(?P<nick>.*)!~(?P<user>.*)@(?P<host>.*) PRIVMSG (?P<to>.*) :(?P<message>.*)`)
	reQuit = regexp.MustCompile(`:(?P<from>.*) PRIVMSG (?P<to>.*) :\.quit`)
	reWho  = regexp.MustCompile(`(?P<nick>.*)!~(?P<user>.*)@(?P<host>.*)`)
)

//=============================================================================
// methods
//=============================================================================

func NewMsg(host string, port string, raw string) *Msg {
	msg := &Msg{
		host: host,
		port: port,
		raw:  raw,
		tags: []string{},
	}
	msg.checkType()
	return msg
}

func (msg *Msg) checkType() {
	raw := msg.Trim()
	if strings.Fields(raw)[0] == PINGMSG {
		msg.tags = append(msg.tags, PINGMSG)
	}
	if strings.Index(raw, PRIVMSG) != -1 {
		if strings.Index(raw, ":\001PING ") != -1 {
			msg.tags = append(msg.tags, PINGMSG)
		}
		msg.tags = append(msg.tags, PRIVMSG)
	}
	if strings.Fields(raw)[0] == msg.host {
		msg.tags = append(msg.tags, SERVMSG)
	}
	if len(msg.tags) == 0 {
		msg.tags = append(msg.tags, UNKNMSG)
	}
}

func (msg *Msg) Tags() []string {
	return msg.tags
}

func (msg *Msg) Raw() string {
	return msg.raw
}

func (msg *Msg) Trim() string {
	if strings.HasPrefix(msg.raw, ":") {
		return strings.TrimPrefix(msg.raw, ":")
	}
	return msg.raw
}

// Parse the raw message and return the map.
func (msg *Msg) Parsemp(r *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	if match := r.FindStringSubmatch(msg.raw); match != nil {
		for i, name := range r.SubexpNames() {
			if i == 0 {
				continue
			}
			result[name] = match[i]
		}
	}
	return result
}

// Parse the raw message and return the slice.
func (msg *Msg) Parsese(r *regexp.Regexp) []string {
	result := []string{}
	mp := msg.Parsemp(r)
	for _, val := range mp {
		result = append(result, val)
	}
	return result
}

func (msg *Msg) ParsePRIVMSG() map[string]string {
	return msg.Parsemp(rePRIVMSG)
}

func hasElement(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func (msg *Msg) IsPRIVMSG() bool {
	return hasElement(msg.Tags(), PRIVMSG)
}

func (msg *Msg) IsSERVMSG() bool {
	return hasElement(msg.Tags(), SERVMSG)
}

func (msg *Msg) IsPINGMSG() bool {
	return hasElement(msg.Tags(), PINGMSG)
}

func (msg *Msg) IsUNKNMSG() bool {
	return hasElement(msg.Tags(), UNKNMSG)
}

func (msg *Msg) GetPRIVMSG() string {
	if msg.IsPRIVMSG() {
		return msg.Parsemp(rePRIVMSG)["message"]
	}
	return ""
}
