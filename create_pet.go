package main

import (
	"time"

	"github.com/yanzay/tbot"
)

func CreatePet(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		pet := petStore.Get(m.ChatID)
		if !pet.Alive {
			petStore.Set(m.ChatID, NewPet(m.ChatID))
			buttons := [][]string{{"Create"}}
			content, err := contentFromTemplate(rootTemplate, pet)
			if err != nil {
				return
			}
			m.Reply(content, tbot.WithMarkdown)
			m.ReplyKeyboard("Your pet is dead. Create new one?", buttons)
			return
		} else {
			if pet.Name != "" && pet.Emoji != "" {
				f(m)
				return
			}
		}
		defer petStore.Set(m.ChatID, pet)
		if pet.AskType {
			switch m.Text() {
			case Chicken.String():
				pet.Emoji = Chicken.Emoji
			case Penguin.String():
				pet.Emoji = Penguin.Emoji
			case Dog.String():
				pet.Emoji = Dog.Emoji
			case Monkey.String():
				pet.Emoji = Monkey.Emoji
			case Fox.String():
				pet.Emoji = Fox.Emoji
			case Panda.String():
				pet.Emoji = Panda.Emoji
			default:
				m.Replyf("Wrong pet type %s", m.Text())
			}
			pet.AskType = false
		}
		if pet.AskName {
			pet.Name = m.Text()
			pet.AskName = false
			pet.Born = time.Now()
			pet.Alive = true
			petStore.Set(pet.PlayerID, pet)
			rootHandler(m)
		}
		if pet.Emoji == "" {
			pet.AskType = true
			pets := [][]string{
				{Chicken.String(), Penguin.String(), Dog.String()},
				{Monkey.String(), Fox.String(), Panda.String()},
			}
			m.ReplyKeyboard("Choose your pet:", pets, tbot.OneTimeKeyboard)
			return
		}
		if pet.Name == "" {
			pet.AskName = true
			m.Reply("Name your pet:")
			return
		}
	}
}
