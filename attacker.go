package main

import (
	"math/rand"
	"sync"
	"time"
)

type missile struct {
	X     int
	Y     int
	speed int
	curr  int
}

var (
	atkName        string
	attacking      bool = false
	missilesLeft   int = -1
	activeMissiles int = 0
	attackerWG     sync.WaitGroup
	atkLock        sync.Mutex
)

func attackerSetup(name string, missiles int) {
	atkLock.Lock()
	atkName = name
	if missiles == 0 {
		missilesLeft = defaultHP / 3 + 10
	} else {
		missilesLeft = missiles 
	}
	rand.Seed(time.Now().UnixNano())
	atkLock.Unlock()
}

func attackerRun(wg *sync.WaitGroup) {
	defer wg.Done()

	min := 0 + 2
	max := MaxX - 2
	
	for missilesLeft > 0 {
		column := rand.Intn(max - min) + min
		spd    := rand.Intn(5) + 1
		newMissile := missile{X: column, Y: MaxY - 1, speed: spd, curr: 0}
		
		attackerWG.Add(1)
		go missileRun(newMissile, &attackerWG)

		missilesLeft--
		activeMissiles++

		n := time.Duration(rand.Intn(2) + 1)
		time.Sleep(n * time.Second)
	}
}

func missileRun(m missile, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select{
		case <- done:
			return
		case <- time.After(Tick):
			if missileHit(&m) {
				atkLock.Lock()
				activeMissiles--
				atkLock.Unlock()
				return
			}
			updateMissile(&m)
		}
	}
}

func attackOver() bool {
	return missilesLeft == 0 && activeMissiles == 0
}