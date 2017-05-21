package main

import "fmt"

type PetType struct {
	Emoji string
	Name  string
}

func (p PetType) String() string {
	return fmt.Sprintf("%s %s", p.Emoji, p.Name)
}

var (
	Chicken = PetType{Emoji: "🐔", Name: "Chicken"}
	Penguin = PetType{Emoji: "🐧", Name: "Penguin"}
	Dog     = PetType{Emoji: "🐶", Name: "Dog"}
	Monkey  = PetType{Emoji: "🐵", Name: "Monkey"}
	Fox     = PetType{Emoji: "🦊", Name: "Fox"}
	Panda   = PetType{Emoji: "🐼", Name: "Panda"}
)
