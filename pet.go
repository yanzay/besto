package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const padWidth = 18
const FirstLevelXP = 1000

type Pet struct {
	PlayerID  int64
	Name      string
	Emoji     string
	Health    int
	Happy     int
	Food      int
	Born      time.Time
	Died      time.Time
	AwakeTime time.Time
	Weight    int
	Mood      string
	Alive     bool
	Sleep     bool
	AskName   bool
	AskType   bool
	XP        int64
}

func NewPet(id int64) *Pet {
	return &Pet{
		PlayerID: id,
		Health:   100,
		Happy:    80,
		Food:     80,
		Mood:     "Good",
		Weight:   1,
		Alive:    true,
	}
}

func (p *Pet) SetMood() {
	switch {
	case !p.Alive:
		p.Mood = "Dead"
	case p.Health < 50:
		p.Mood = "Sick"
	case p.Food < 20:
		p.Mood = "Hungry"
	case p.Happy < 5:
		p.Mood = "Stress"
	case p.Happy < 50:
		p.Mood = "Sorrow"
	case p.Happy >= 100:
		p.Mood = "Great"
	default:
		p.Mood = "Good"
	}
}

func (p *Pet) Age() time.Duration {
	var d time.Duration
	if p.Alive {
		d = time.Since(p.Born)
	} else {
		d = p.Died.Sub(p.Born)
	}
	return roundDuration(d)
}

func (p *Pet) Die() {
	p.Health = 0
	p.Alive = false
	p.Died = time.Now()
	go historyStore.Create(p)
}

func (p *Pet) XPForLevel(level int) int64 {
	return int64(FirstLevelXP * math.Pow(2.1, float64(level-1)))
}

func (p *Pet) Level() int {
	if p.XP < FirstLevelXP {
		return 0
	}
	return int(math.Log(float64(p.XP)/FirstLevelXP)/math.Log(2.1)) + 1
}

func (p *Pet) XPForNextLevel() int64 {
	return int64(float64(p.XPForLevel(p.Level())) * 1.1)
}

func (p *Pet) XPFromCurrentLevel() int64 {
	return p.XP - p.XPForLevel(p.Level())
}

func (p *Pet) XPString() string {
	return fmt.Sprintf("%d/%d", p.XPFromCurrentLevel(), p.XPForNextLevel())
}

func roundDuration(d time.Duration) time.Duration {
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
	return pad("Weight", fmt.Sprintf("%dg", p.Weight))
}

func (p *Pet) TopString() string {
	name := fmt.Sprintf("%s%s", p.Emoji, p.Name)
	deadStr := ""
	if !p.Alive {
		deadStr = "💀"
	}
	return pad(name, p.Age().String()) + deadStr
}

func pad(first, last string) string {
	repeatCount := padWidth - utf8.RuneCountInString(first) - utf8.RuneCountInString(last)
	if repeatCount < 0 {
		repeatCount = 1
	}
	return first + strings.Repeat(" ", repeatCount) + last
}