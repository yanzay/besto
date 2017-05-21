package main

import (
	"time"

	"github.com/yanzay/tbot"
)

func Sleep(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		pet := petStore.Get(m.ChatID)
		if pet.Sleep {
			if time.Until(pet.AwaikTime) > time.Minute {
				m.Replyf("Your pet is sleeping. Time to wake up: %s", roundDuration(time.Until(pet.AwaikTime)))
			} else {
				m.Reply("Your pet will wake up soon.")
			}
			return
		}
		f(m)
	}
}
