package main

import (
	s "bot/scrapper"
	t "bot/telegram"
)

type Config struct {
	MyUsers  []s.User   `yaml:"traders"`
	Telegram t.Telegram `yaml:"telegram"`
	Interval int        `yaml:"interval"`
}
