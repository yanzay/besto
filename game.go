package main

import "time"

const (
	SpeedFood         = 2
	SpeedHappy        = 1
	SpeedHealth       = 2
	SpeedNormalWeight = 1
	SpeedOverWeight   = 2
	NormalWeight      = 42
)

var (
	moveDuration       = 60 * time.Second
	sleepCheckDuration = 5 * time.Second
)

func mainLoop() {
	tick := time.Tick(moveDuration)
	for range tick {
		for _, pet := range petStore.Alive() {
			if !pet.Sleep {
				petStore.Update(pet.PlayerID, func(pet *Pet) {
					decreaseFood(pet)
					decreaseHappy(pet)
					decreaseHealth(pet)
					pet.Weight += getWeightDelta(pet)
					if pet.Weight < 2 {
						pet.Weight = 1
					}
				})
			}
		}
	}
}

func sleepLoop() {
	tick := time.Tick(sleepCheckDuration)
	for range tick {
		for _, pet := range petStore.Alive() {
			if pet.Sleep && pet.AwakeTime.Before(time.Now()) {
				petStore.Update(pet.PlayerID, func(p *Pet) {
					p.Sleep = false
					p.Notify("Good morning!")
				})
			}
		}
	}
}

func decreaseFood(pet *Pet) {
	if pet.Food > SpeedFood {
		pet.Food -= SpeedFood
	} else {
		if pet.Food > 0 {
			pet.Notify("Hey! I am hungry!")
		}
		pet.Food = 0
	}
}

func decreaseHappy(pet *Pet) {
	speed := SpeedHappy
	if pet.Food == 0 {
		speed *= 2
	}
	if pet.Happy > speed {
		pet.Happy -= speed
	} else {
		if pet.Happy > 0 {
			pet.Notify("Hey! I am bored!")
		}
		pet.Happy = 0
	}
}

func decreaseHealth(pet *Pet) {
	if pet.Happy == 0 || pet.Food == 0 {
		if pet.Health < 10 {
			pet.Notify("I'm dying! Please help me!")
		}
		if pet.Health > SpeedHealth {
			pet.Health -= SpeedHealth
		} else {
			pet.Die()
		}
	}
}

func getWeightDelta(pet *Pet) int {
	switch {
	case pet.Food > 120:
		return SpeedOverWeight
	case pet.Food > 80 && pet.Weight < NormalWeight:
		return SpeedNormalWeight
	case pet.Food > 80 && pet.Weight > NormalWeight:
		return -SpeedNormalWeight
	case pet.Food > 20 && pet.Weight > NormalWeight:
		return -SpeedNormalWeight
	case pet.Food < 20:
		return -SpeedOverWeight
	}
	return 0
}
