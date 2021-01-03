package main

import "sync"

type defender struct {
	name  string
	X     int
	Y     int
	hp    int
	score int
}

const (
	defaultHP = 100
)

var (
	def     defender
	defLock sync.Mutex
)

func defenderSetup(name string) {
	defLock.Lock()
	def.name  = name
	def.X     = (MaxX - (PaddleLength / 2)) / 2
	def.Y     = buildingHeight + 3
	def.hp    = defaultHP
	def.score = 0
	defLock.Unlock()
	boardLoadDefender(def)
}

func defScored() {
	defLock.Lock()
	def.score++
	defLock.Unlock()
}

func defHit() {
	defLock.Lock()
	def.hp -= 3
	defLock.Unlock()
}

func defenseOver() bool {
	return def.hp <= 0
}
