/*
 * ircbot.go
 *
 * Copyright (C) 2011-2012,  Alex Petrovich <alex@libslack.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
