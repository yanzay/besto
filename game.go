package main

import "time"

var moveDuration = 60 * time.Second

func gameLoop() {
	tick := time.Tick(moveDuration)
	for range tick {
		for _, pet := range petStore.Alive() {
			petLoop(pet)
		}
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
	petStore.Update(pet.PlayerID, func(pet *Pet) {
		if pet.Alive {
			if pet.AwaikTime.Before(time.Now()) {
				pet.Sleep = false
			}
			if pet.Sleep {
				return
			}
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
				if pet.Health > SpeedHealth {
					pet.Health -= SpeedHealth
				} else {
					pet.Health = 0
					pet.Alive = false
					pet.Died = time.Now()
				}
			}
			pet.Weight += getWeightDelta(pet)
			pet.Mood = getMood(pet)
		}
	})
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
	if pet.Food < 20 {
		delta -= SpeedOverWeight
	}
	return delta
}

func getMood(pet *Pet) string {
	if !pet.Alive {
		return "Dead"
	}
	if pet.Food < 20 {
		return "Hungry"
	}
	if pet.Happy == 0 {
		return "Stress"
	}
	if pet.Happy < 20 {
		return "Sorrow"
	}
	if pet.Happy >= 100 {
		return "Great!"
	}
	return "Good"
}
