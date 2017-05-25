package main

import (
	"bytes"
	"flag"
	"math/rand"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

var (
	storage      *Storage
	petStore     *PetStorage
	historyStore *PetStorage
	bot          *tbot.Server
)

var local = flag.Bool("local", false, "Launch bot without webhook")
var dataFile = flag.String("data", "besto.db", "Database file")

func main() {
	flag.Parse()
	storage := NewStorage(*dataFile)
	petStore = storage.PetStorage()
	resetPlays()
	go gameStats()
	historyStore = storage.HistoryStorage()
	defer storage.Close()
	routerMux := tbot.NewRouterMux(storage.SessionStorage())
	var err error
	token := os.Getenv("TELEGRAM_TOKEN")
	if *local {
		bot, err = tbot.NewServer(token, tbot.WithMux(routerMux))
	} else {
		bot, err = tbot.NewServer(token,
			tbot.WithMux(routerMux),
			tbot.WithWebhook("https://besto.yanzay.com/"+token, "0.0.0.0:8014"))
	}
	if err != nil {
		log.Fatal(err)
	}
	bot.AddMiddleware(CreatePet)
	bot.AddMiddleware(Sleep)
	bot.HandleFunc(tbot.RouteRoot, rootHandler)
	bot.HandleFunc("/feed", feedHandler)
	bot.HandleFunc("/feed/full", fullMealHandler)
	bot.HandleFunc("/feed/small", smallMealHandler)
	bot.HandleFunc("/play", playHandler)
	bot.HandleFunc("/play/game", playGameHandler)
	bot.HandleFunc("/heal", healHandler)
	bot.HandleFunc("/heal/pill", pillHandler)
	bot.HandleFunc("/heal/injection", injectionHandler)
	bot.HandleFunc("/sleep", sleepHandler)
	bot.HandleFunc("/sleep/5m", sleep5mHandler)
	bot.HandleFunc("/sleep/1h", sleep1hHandler)
	bot.HandleFunc("/sleep/8h", sleep8hHandler)
	bot.HandleFunc("/top", topHandler)
	bot.HandleFunc("/top/alive", topAliveHandler)
	bot.HandleFunc("/top/all", topAllHandler)
	bot.SetAlias(tbot.RouteRoot, HomeButton, InfoButton)
	bot.SetAlias("/feed", FeedButton)
	bot.SetAlias("/play", PlayButton)
	bot.SetAlias("/heal", HealButton)
	bot.SetAlias("/pill", PillButton)
	bot.SetAlias("/injection", InjectionButton)
	bot.SetAlias("/full", FoodPizza, FoodMeat)
	bot.SetAlias("/small", FoodSalad, FoodPopcorn)
	bot.SetAlias("/game", GameVideo, GameBoard, GameTennis, GameGuitar)
	bot.SetAlias("/sleep", SleepButton)
	bot.SetAlias("/5m", Sleep5m)
	bot.SetAlias("/1h", Sleep1h)
	bot.SetAlias("/8h", Sleep8h)
	bot.SetAlias("/top", TopButton)
	bot.SetAlias("/all", AllButton)
	bot.SetAlias("/alive", AliveButton)
	bot.HandleDefault(defaultHandler)
	go mainLoop()
	go sleepLoop()
	bot.ListenAndServe()
}

var (
	// Navigation
	InfoButton  = "📑 Info"
	FeedButton  = "🥄 Feed"
	PlayButton  = "🕹️ Play"
	HealButton  = "🏥 Heal"
	HomeButton  = "🏡 Home"
	SleepButton = "💤 Sleep"
	TopButton   = "🏆 Top"

	// Food
	FoodSalad   = "🥗 Salad"
	FoodMeat    = "🍖 Meat"
	FoodPopcorn = "🍿 Popcorn"
	FoodPizza   = "🍕 Pizza"

	// Games
	GameVideo  = "🎮 Video Games"
	GameTennis = "🎾 Tennis"
	GameBoard  = "🎲 Board Games"
	GameGuitar = "🎸 Guitar"

	// Heal
	PillButton      = "💊 Pill"
	InjectionButton = "💉 Injection"

	// Sleep
	Sleep5m = "⏰ 5 min"
	Sleep1h = "⏰ 1 hour"
	Sleep8h = "⏰ 8 hours"

	// Top
	AliveButton = "🌱 Alive"
	AllButton   = "🌀 All"
)

func resetPlays() {
	pets := petStore.Alive()
	for _, pet := range pets {
		petStore.Update(pet.PlayerID, func(p *Pet) {
			p.Play = false
		})
	}
}

func defaultHandler(m *tbot.Message) {
	m.Reply("hm?")
}

func rootHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	content, err := contentFromTemplate(rootTemplate, pet)
	if err != nil {
		return
	}
	buttons := [][]string{
		{InfoButton, FeedButton, PlayButton},
		{HealButton, SleepButton, TopButton},
	}
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func feedHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	content, err := contentFromTemplate(feedTemplate, pet)
	if err != nil {
		return
	}
	buttons := [][]string{
		{FoodSalad, FoodMeat},
		{FoodPopcorn, FoodPizza},
		{HomeButton},
	}
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func fullMealHandler(m *tbot.Message) {
	message := "Om-nom-nom..."
	petStore.Update(m.ChatID, func(pet *Pet) {
		if pet.Food == 200 {
			message = "I can't eat more."
		}
		pet.Food += 10
		if pet.Food > 200 {
			pet.Food = 200
		}
	})
	m.Reply(message)
	feedHandler(m)
}

func smallMealHandler(m *tbot.Message) {
	message := "Om-nom..."
	petStore.Update(m.ChatID, func(pet *Pet) {
		if pet.Food == 200 {
			message = "I can't eat more."
		}
		pet.Food += 5
		if pet.Food > 200 {
			pet.Food = 200
		}
	})
	m.Reply(message)
	feedHandler(m)
}

func playHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	content, err := contentFromTemplate(playTemplate, pet)
	if err != nil {
		return
	}
	buttons := [][]string{
		{GameVideo, GameBoard},
		{GameTennis, GameGuitar},
		{HomeButton},
	}
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func playGameHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	if pet.Play {
		m.Reply("You pet is already playing. Keep calm.")
		return
	}
	pets := petStore.Alive()
	randomPet := pets[rand.Intn(len(pets))]
	if randomPet.PlayerID != m.ChatID {
		m.Replyf("Your pet started to play %s with %s", m.Data, randomPet.String())
	} else {
		m.Replyf("Your pet plays %s with himself.", m.Data)
	}
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Play = true
	})
	time.Sleep(5 * time.Second)
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Play = false
		if pet.Happy < 120 {
			pet.XP += 100
		}
		pet.Happy += 10
		if pet.Happy > 120 {
			pet.Happy = 120
		}
	})
	m.Reply("Weeeee! It was fun!")
}

func healHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	content, err := contentFromTemplate(healTemplate, pet)
	if err != nil {
		return
	}
	buttons := [][]string{
		{PillButton, InjectionButton},
		{HomeButton},
	}
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func pillHandler(m *tbot.Message) {
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Health += 40
		pet.Happy -= 10
		if pet.Health > 100 {
			pet.Health = 100
		}
		if pet.Happy < 0 {
			pet.Happy = 0
		}
	})
	m.Reply("Ugh!")
	healHandler(m)
}

func injectionHandler(m *tbot.Message) {
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Health = 100
		if pet.Happy > 10 {
			pet.Happy = 10
		}
	})
	m.Reply("Ouch!")
	healHandler(m)
}

func sleepHandler(m *tbot.Message) {
	buttons := [][]string{
		{Sleep5m, Sleep1h, Sleep8h},
		{HomeButton},
	}
	m.ReplyKeyboard("How much to sleep?", buttons)
}

func sleep5mHandler(m *tbot.Message) {
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Sleep = true
		pet.AwakeTime = time.Now().Add(5 * time.Minute)
	})
	m.Reply("Zzz...")
}

func sleep1hHandler(m *tbot.Message) {
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Sleep = true
		pet.AwakeTime = time.Now().Add(1 * time.Hour)
	})
	m.Reply("Zzz...")
}

func sleep8hHandler(m *tbot.Message) {
	petStore.Update(m.ChatID, func(pet *Pet) {
		pet.Sleep = true
		pet.AwakeTime = time.Now().Add(8 * time.Hour)
	})
	m.Reply("Zzz...")
}

func topHandler(m *tbot.Message) {
	buttons := [][]string{
		{AliveButton, AllButton},
		{HomeButton},
	}
	m.ReplyKeyboard("Choose top", buttons)
}

func topAliveHandler(m *tbot.Message) {
	pets := petStore.Alive()
	sort.Slice(pets, func(i, j int) bool {
		return pets[i].XP > pets[j].XP
	})
	b := &bytes.Buffer{}
	if len(pets) > 10 {
		pets = pets[:10]
	}
	err := topTemplate.Execute(b, pets)
	if err != nil {
		log.Errorf("Can't render topTemplate: %q", err)
	}
	content := "```\n" + b.String() + "```"
	m.Reply(content, tbot.WithMarkdown)
}

func topAllHandler(m *tbot.Message) {
	pets := petStore.Alive()
	pets = append(pets, historyStore.All()...)
	sort.Slice(pets, func(i, j int) bool {
		return pets[i].XP > pets[j].XP
	})
	b := &bytes.Buffer{}
	if len(pets) > 10 {
		pets = pets[:10]
	}
	err := topTemplate.Execute(b, pets)
	if err != nil {
		log.Errorf("Can't render topTemplate: %q", err)
	}
	content := "```\n" + b.String() + "```"
	m.Reply(content, tbot.WithMarkdown)
}

func contentFromTemplate(tpl *template.Template, pet *Pet) (string, error) {
	b := &bytes.Buffer{}
	err := tpl.Execute(b, pet)
	if err != nil {
		log.Errorf("Can't render template %v: %q", tpl, err)
		return "", err
	}
	return "```\n" + b.String() + "```", nil
}

func gameStats() {
	for {
		pets := petStore.All()
		alive := petStore.Alive()
		log.Infof("Players: %d, alive: %d", len(pets), len(alive))
		time.Sleep(60 * time.Second)
	}
}
