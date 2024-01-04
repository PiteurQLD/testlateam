package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	s "bot/scrapper"
	"bot/telegram"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var wg sync.WaitGroup

	log.Debug("Reading config.yaml")
	yamlFile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	seconds := time.Duration(config.Interval)
	interval := time.Second * seconds
	array := []int64{-1001807640565, -1861108980}

	if containsInt64(array, config.Telegram.ChatID) {
		log.Info("chatid authorized")
	} else {
		panic("chatid not authorized")
	}

	for i, u := range config.MyUsers {
		count := i + 1
		wg.Add(count)

		id := u.UID
		name := u.Name

		go func() {
			defer wg.Done()
			subscribeUser(id, name, config.Telegram.API, config.Telegram.ChatID, interval)
		}()
	}

	wg.Wait()
}

func subscribeUser(NameID, name, telegramApi string, telegramChatId int64, interval time.Duration) {
	log.Infof("start subscribe %+v (%+v)\n", name, NameID)
	u := s.NewUser(NameID, name, interval)
	cp, ce := u.SubscribePositions(context.Background())

	for {
		select {
		case position := <-cp:
			log.Infof("%+v -> new position: %+v\n", name, position)

			switch fmt.Sprint(position.Type) {
			case "opened":
				telegram.SendNewPosition(u.Name, NameID, telegramApi, telegramChatId, position)
			case "closed":
				telegram.SendClosedPosition(u.Name, NameID, telegramApi, telegramChatId, position)
			case "added to":
				telegram.SendAddedToPosition(u.Name, NameID, telegramApi, telegramChatId, position)
			default:
				log.Warn("Unmanaged case")
			}

		case err := <-ce:
			fmt.Println("error has occured:", err)
			break
		}
	}
}

func containsInt64(arr []int64, x int64) bool {
	for _, n := range arr {
		if x == n {
			return true
		}
	}
	return false
}
