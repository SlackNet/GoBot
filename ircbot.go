package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strings"
)

var c net.Conn

func handleCommand(sender, command, channel, message string) {
	if sender == "NickServ!service@rizon.net" && strings.Contains(message, "please choose a different nick.") {
		write("PRIVMSG NickServ :IDENTIFY password")
	}

	if strings.Contains(message, "will originate from") {
		write("JOIN #x86")
	}

	if channel[0] == '#' {
		if message[0] == '`' {
			if strings.HasPrefix(message, "`raw ") && sender == "slacky!~Sl@ck.ware" {
				write(message[5:])
			}
		}
	}
}

func write(s string) {
	if c == nil {
		fmt.Println("Connection is nil.")
		os.Exit(1)
	}

	s += "\n"

	i, _ := c.Write([]byte(s))
	if i < 1 {
		fmt.Println("WARNING: Wrote nothing.")
	}


	fmt.Print(" >> " + s)
}

func main() {
	c, _ = net.Dial("tcp4", "irc.rizon.net:6667")

	for {
		buffer := make([]byte, 10480)

		i, _ := c.Read(buffer)

		if i < 1 {
			time.Sleep(1000000000)
			continue
		}

		lines := strings.Split(string(buffer[0:i]), "\n")
		for _, line := range(lines) {
			if len(line) > 0 {
				fmt.Println(" << " + line)
				words := strings.Split(line, " ")
				if words[0] == "PING" {
					write("PONG " + line[5:])
				}
				if line[0] == ':' {
					if len(words) > 2 {
						sender := words[0][1:]
						command := words[1]
						channel := words[2]

						if strings.Contains(command, "439") {
							write("NICK GoBot")
							write("USER GoBot Go.Bot irc.rizon.net Go Bot")
						}

						// this is lame and I'm being lazy.
						prefix := sender + " " + command + " " + channel + " :"

						if strings.Contains(line, prefix) {
							handleCommand(sender, command, channel, line[strings.Index(line, prefix) + len(prefix):])
						}
					}
				}
			}
		}

		time.Sleep(.5*1000000000)
	}
}
