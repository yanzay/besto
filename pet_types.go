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
	Chicken = PetType{Emoji: "🐔", Adult: "🐓", Name: "Chicken"}
	Penguin = PetType{Emoji: "🐧", Name: "Penguin"}
	Dog     = PetType{Emoji: "🐶", Adult: "🐕", Name: "Dog"}
	Monkey  = PetType{Emoji: "🐵", Adult: "🐒", Name: "Monkey"}
	Fox     = PetType{Emoji: "🦊", Name: "Fox"}
	Panda   = PetType{Emoji: "🐼", Name: "Panda"}
	Cat     = PetType{Emoji: "🐱", Adult: "🐈", Name: "Cat"}
	Pig     = PetType{Emoji: "🐷", Adult: "🐖", Name: "Pig"}
	Rabbit  = PetType{Emoji: "🐰", Adult: "🐇", Name: "Rabbit"}
	Mouse   = PetType{Emoji: "🐭", Adult: "🐁", Name: "Mouse"}
	Tiger   = PetType{Emoji: "🐯", Adult: "🐅", Name: "Tiger"}
	Lizard  = PetType{Emoji: "🦎", Adult: "🐉", Name: "Lizard"}
)
