// Package filter provides ...
package filter

import (
	"regexp"
)

var (
	ServerPing = regexp.MustCompile(`^PING :(?P<server>.*)`)
	UserPing = regexp.MustCompile(
		`:(?P<who>.*) PRIVMSG (?P<target>.*) :\001PING (?P<timestamp>.*)\001`)
)
