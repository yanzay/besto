package main

import "time"

var moveDuration = 60 * time.Second
var sleepCheckDuration = 5 * time.Second

func mainLoop() {
	tick := time.Tick(moveDuration)
	for range tick {
		for _, pet := range petStore.Alive() {
			if pet.Alive && !pet.Sleep {
				petStore.Update(pet.PlayerID, func(pet *Pet) {
					decreaseFood(pet)
					decreaseHappy(pet)
					decreaseHealth(pet)
					pet.Weight += getWeightDelta(pet)
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
				})
			}
		}
	}
}

const (
	SpeedFood         = 2
	SpeedHappy        = 1
	SpeedHealth       = 2
	SpeedNormalWeight = 1
	SpeedOverWeight   = 2
	NormalWeight      = 42
)

func decreaseFood(pet *Pet) {
	if pet.Food >= SpeedFood {
		pet.Food -= SpeedFood
	} else {
		pet.Food = 0
	}
}

func decreaseHappy(pet *Pet) {
	pet.Happy -= SpeedHappy
	if pet.Food == 0 {
		pet.Happy -= SpeedHappy
	}
	if pet.Happy < 0 {
		pet.Happy = 0
	}
}

func decreaseHealth(pet *Pet) {
	if pet.Happy == 0 || pet.Food == 0 {
		if pet.Health > SpeedHealth {
			pet.Health -= SpeedHealth
		} else {
			pet.Health = 0
			pet.Alive = false
			pet.Died = time.Now()
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
