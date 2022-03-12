#!/usr/bin/expect

spawn ssh admin@172.172.172.172
expect "*: " {send "P@ssw0rd\r"}
#telnet to CSW
expect "*> " {send "172.24.32.10\r"}
expect "*: " {send "adminnet\r"}
expect "*: " {send "adminnet#01\r"}
expect "*>" {send "enable\r"}
expect "*: " {send "adminnet#01\r"}
#telnet to VPN
expect "*#" {send "telnet 172.24.33.5\r"}
expect "*: " {send "adminnet\r"}
expect "*: " {send "adminnet#01\r"}
expect "*>" {send "enable\r"}
expect "*: " {send "adminnet#01\r"}
expect "*#" {send "terminal length 0\r"}
expect "*#" {send "show run | i user\r"}
expect "*#" {send -- "exit\r"}
