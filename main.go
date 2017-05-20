package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

type PetType struct {
	Emoji string
	Name  string
}

func (p PetType) String() string {
	return fmt.Sprintf("%s %s", p.Emoji, p.Name)
}

var Chicken = PetType{Emoji: "🐔", Name: "Chicken"}
var Penguin = PetType{Emoji: "🐧", Name: "Penguin"}
var Dog = PetType{Emoji: "🐶", Name: "Dog"}
var Monkey = PetType{Emoji: "🐵", Name: "Monkey"}
var Fox = PetType{Emoji: "🦊", Name: "Fox"}
var Panda = PetType{Emoji: "🐼", Name: "Panda"}

type Pet struct {
	sync.Mutex
	PlayerID int64
	Name     string
	Type     *PetType
	Health   int
	Happy    int
	Food     int
	Born     time.Time
	Weight   int
	Mood     string
	Alive    bool
	askName  bool
	askType  bool
}

func (p *Pet) Age() time.Duration {
	d := time.Since(p.Born)
	return d - (d % time.Second)
}

func (p *Pet) AgeString() string {
	return pad("Age", p.Age().String())
}

func (p *Pet) HealthString() string {
	return pad("Health", strconv.Itoa(p.Health))
}

func (p *Pet) HappyString() string {
	return pad("Happy", strconv.Itoa(p.Happy))
}

func (p *Pet) FoodString() string {
	return pad("Food", strconv.Itoa(p.Food))
}

func (p *Pet) MoodString() string {
	return pad("Mood", p.Mood)
}

func (p *Pet) WeightString() string {
	return pad("Weight", strconv.Itoa(p.Weight))
}

const padWidth = 20

func pad(first, last string) string {
	return first + strings.Repeat(" ", padWidth-len(first)-len(last)) + last
}

func NewPet(id int64) *Pet {
	return &Pet{
		PlayerID: id,
		Health:   100,
		Happy:    100,
		Food:     100,
		Mood:     "Good",
		Weight:   1,
	}
}

type PetStore struct {
	sync.Mutex
	pets map[int64]*Pet
}

func NewPetStore() *PetStore {
	return &PetStore{pets: make(map[int64]*Pet)}
}

var petStore = NewPetStore()

func (ps *PetStore) Get(id int64) *Pet {
	ps.Lock()
	defer ps.Unlock()
	_, ok := ps.pets[id]
	if !ok {
		ps.pets[id] = NewPet(id)
	}
	return ps.pets[id]
}

func (ps *PetStore) Set(id int64, pet *Pet) {
	ps.Lock()
	ps.pets[id] = pet
	ps.Unlock()
}

func CreatePet(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		pet := petStore.Get(m.ChatID)
		if pet.Name != "" && pet.Type != nil {
			f(m)
			return
		}
		defer petStore.Set(m.ChatID, pet)
		if pet.askType {
			switch m.Text() {
			case Chicken.String():
				pet.Type = &Chicken
			case Penguin.String():
				pet.Type = &Penguin
			case Dog.String():
				pet.Type = &Dog
			case Monkey.String():
				pet.Type = &Monkey
			case Fox.String():
				pet.Type = &Fox
			case Panda.String():
				pet.Type = &Panda
			default:
				m.Replyf("Wrong pet type %s", m.Text())
			}
			pet.askType = false
		}
		if pet.askName {
			pet.Name = m.Text()
			pet.askName = false
			pet.Born = time.Now()
			pet.Alive = true
			rootHandler(m)
		}
		if pet.Type == nil {
			pet.askType = true
			pets := [][]string{
				{Chicken.String(), Penguin.String(), Dog.String()},
				{Monkey.String(), Fox.String(), Panda.String()},
			}
			m.ReplyKeyboard("Choose your pet:", pets, tbot.OneTimeKeyboard)
			return
		}
		if pet.Name == "" {
			pet.askName = true
			m.Reply("Name your pet:")
			return
		}
	}
}

func main() {
	routerMux := tbot.NewRouterMux(tbot.NewSessionStorage())
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"), tbot.WithMux(routerMux))
	if err != nil {
		log.Fatal(err)
	}
	bot.AddMiddleware(CreatePet)
	bot.HandleFunc(tbot.RouteRoot, rootHandler)
	bot.HandleFunc("/feed", feedHandler)
	bot.HandleFunc("/feed/full", fullMealHandler)
	bot.HandleFunc("/feed/small", smallMealHandler)
	bot.HandleFunc("/play", playHandler)
	bot.HandleFunc("/play/game", playGameHandler)
	bot.HandleFunc("/heal", healHandler)
	bot.HandleFunc("/heal/pill", pillHandler)
	bot.HandleFunc("/heal/injection", injectionHandler)
	bot.SetAlias(tbot.RouteRoot, "Home", HomeButton)
	bot.SetAlias(tbot.RouteBack, "Back", BackButton)
	bot.SetAlias(tbot.RouteRefresh, "Info", InfoButton)
	bot.SetAlias("/feed", FeedButton)
	bot.SetAlias("/play", PlayButton)
	bot.SetAlias("/heal", HealButton)
	bot.SetAlias("/full", FoodPizza, FoodMeat)
	bot.SetAlias("/small", FoodSalad, FoodPopcorn)
	bot.SetAlias("/game", GameVideo, GameBoard, GameTennis, GameGuitar)
	bot.HandleDefault(defaultHandler)
	go gameLoop()
	bot.ListenAndServe()
}

var (
	InfoButton = "📑 Info"
	FeedButton = "🥄 Feed"
	PlayButton = "🕹️ Play"
	HealButton = "🏥 Heal"
)

var moveDuration = 2 * time.Second

func gameLoop() {
	tick := time.Tick(moveDuration)
	for range tick {
		petStore.Lock()
		for _, pet := range petStore.pets {
			petLoop(pet)
		}
		petStore.Unlock()
	}
}

const (
	SpeedFood         = 2
	SpeedHappy        = 1
	SpeedHealth       = 4
	SpeedNormalWeight = 1
	SpeedOverWeight   = 2
	NormalWeight      = 42
)

func petLoop(pet *Pet) {
	pet.Lock()
	if pet.Alive {
		if pet.Food >= SpeedFood {
			pet.Food -= SpeedFood
		} else {
			pet.Food = 0
		}
		if pet.Happy >= SpeedHappy {
			pet.Happy -= SpeedHappy
		} else {
			pet.Happy = 0
		}
		if pet.Happy == 0 || pet.Food == 0 {
			pet.Happy = 0
			if pet.Health >= SpeedHealth {
				pet.Health -= SpeedHealth
			} else {
				pet.Health = 0
				pet.Alive = false
			}
		}
		pet.Weight += getWeightDelta(pet)
		pet.Mood = getMood(pet)
	}
	pet.Unlock()
}

func getWeightDelta(pet *Pet) int {
	delta := 0
	if pet.Food > 100 {
		delta += SpeedOverWeight
	}
	if pet.Food < 100 && pet.Food > 80 && pet.Weight < NormalWeight {
		delta += SpeedNormalWeight
	}
	if pet.Food < 100 && pet.Weight > NormalWeight {
		delta -= SpeedNormalWeight
	}
	return delta
}

func getMood(pet *Pet) string {
	if !pet.Alive {
		return "Dead"
	}
	if pet.Food < 80 {
		return "Hungry"
	}
	if pet.Happy == 0 {
		return "Stress"
	}
	if pet.Happy < 80 {
		return "Sorrow"
	}
	return "Good"
}

func defaultHandler(m *tbot.Message) {
	m.Reply("hm?")
}

var rootTemplate = template.Must(template.New("root").Parse(
	`{{ .Type.Emoji }} {{ .Name }} {{ if not .Alive }}☠️{{ end }}
{{ .AgeString }}
{{ .WeightString }}
{{ .MoodString }}
{{ .HealthString }}
{{ .FoodString }}
{{ .HappyString }}
`))

var feedTemplate = template.Must(template.New("feed").Parse(
	`{{ .FoodString }}
What do you prefer?
`))

var playTemplate = template.Must(template.New("feed").Parse(
	`{{ .HappyString }}
Let's play!
`))

var healTemplate = template.Must(template.New("heal").Parse(
	`{{ .HealthString }}
Heal me...
`))

func rootHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	b := &bytes.Buffer{}
	err := rootTemplate.Execute(b, pet)
	if err != nil {
		log.Errorf("Can't render rootTemplate: %q", err)
		return
	}
	buttons := [][]string{
		{InfoButton, FeedButton, PlayButton},
		{HealButton},
	}
	content := "```\n" + b.String() + "```"
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

var (
	FoodSalad   = "🥗 Salad"
	FoodMeat    = "🍖 Meat"
	FoodPopcorn = "🍿 Popcorn"
	FoodPizza   = "🍕 Pizza"
)

func feedHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	b := &bytes.Buffer{}
	err := feedTemplate.Execute(b, pet)
	if err != nil {
		log.Errorf("Can't render feedTemplate: %q", err)
		return
	}
	buttons := [][]string{
		{FoodSalad, FoodMeat},
		{FoodPopcorn, FoodPizza},
		{BackButton, HomeButton},
	}
	content := "```\n" + b.String() + "```"
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func fullMealHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	pet.Lock()
	pet.Food += 10
	pet.Unlock()
	m.Reply("Om-nom-nom...")
	feedHandler(m)
}

func smallMealHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	pet.Lock()
	pet.Food += 5
	pet.Unlock()
	m.Reply("Om-nom...")
	feedHandler(m)
}

var (
	GameVideo  = "🎮 Video Games"
	GameTennis = "🎾 Tennis"
	GameBoard  = "🎲 Board Games"
	GameGuitar = "🎸 Guitar"
)

func playHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	b := &bytes.Buffer{}
	err := playTemplate.Execute(b, pet)
	if err != nil {
		log.Errorf("Can't render play template: %q", err)
		return
	}
	buttons := [][]string{
		{GameVideo, GameBoard},
		{GameTennis, GameGuitar},
		{BackButton, HomeButton},
	}
	content := "```\n" + b.String() + "```"
	m.ReplyKeyboard(content, buttons, tbot.WithMarkdown)
}

func playGameHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	pet.Lock()
	pet.Happy += 10
	if pet.Happy > 100 {
		pet.Happy = 100
	}
	pet.Unlock()
	m.Reply("Weeeee! It's fun!")
	playHandler(m)
}

func healHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	b := &bytes.Buffer{}
	err := healTemplate.Execute(b, pet)
	if err != nil {
		log.Errorf("Can't render heal template: %q", err)
		return
	}
	buttons := [][]string{
		{"💊 Pill", "💉 Injection"},
		{BackButton, HomeButton},
	}
	m.ReplyKeyboard("Heal me...", buttons)
}

func pillHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	pet.Lock()
	defer pet.Unlock()
	pet.Health += 40
	pet.Happy -= 10
	if pet.Health > 100 {
		pet.Health = 100
	}
	if pet.Happy < 0 {
		pet.Happy = 0
	}
	m.Reply("Ugh!")
	healHandler(m)
}

func injectionHandler(m *tbot.Message) {
	pet := petStore.Get(m.ChatID)
	pet.Lock()
	defer pet.Unlock()
	pet.Health = 100
	if pet.Happy > 10 {
		pet.Happy = 10
	}
	m.Reply("Ouch!")
	healHandler(m)
}

var BackButton = "⬅️  Back"
var HomeButton = "🏡 Home"
