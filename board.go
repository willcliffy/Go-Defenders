package main

import (
	"sync"
	"strconv"

	"github.com/jroimartin/gocui"
)

const (
	MaxX = 111
	MaxY =  30
	BoardBuffer  = 5
	PaddleLength = 6
)

var (
	boardLock      sync.Mutex
	board          [MaxY+BoardBuffer][MaxX] byte
	buildingHeight int
)

func boardSetup(skyline []int) {
	boardLock.Lock()

    for y := MaxY+BoardBuffer-1; y >= 0; y-- {
        for x := 0; x < MaxX; x++ {
            board[y][x] = ' '
        }
	}

	var l      = len(skyline)
	var lasty  = skyline[0]
	board[lasty][0] = '_'
	var champy = -1
	for x := 1; x < l; x++ {
		y := skyline[x]
		board[y][x] = '_'
		if y > lasty {
			for yi := 2; yi < y; yi++ {
				board[yi][x-1] = '|'
			}
		} else if y < lasty {
			for yi := lasty; yi > 2; yi-- {
				board[yi-1][x] = '|'
			}
		}

		if lasty > champy {
			champy = lasty
		}
		lasty = y
	}
	buildingHeight = champy
	boardLock.Unlock()
}

func endgameScreen() {
	boardLock.Lock()
	
	var lines []string
	if def.hp > 0 {
		lines = loadFile(YouWin)
	} else {
		lines = loadFile(YouLose) 
	}

	Iy := len(lines)
	var Ix int = -1

	for i := Iy-1; i >= 0; i-- {
		if Ix == -1 {
			Ix = len(lines[i])
		}
		leftBuffer := (MaxX - Ix) / 2
		botBuffer  := (MaxY - Iy - buildingHeight) /  2 + buildingHeight
		for c := range lines[i] {
			board[botBuffer+i][leftBuffer+c] = lines[Iy-1-i][c]
		}
	}

	boardLock.Unlock()
}

func boardLoadDefender(d defender) {
	for x := d.X; x < d.X + PaddleLength + 1; x++ {
		board[d.Y][x] = '#'
	}
}

func defLeft(g *gocui.Gui, v *gocui.View) error {
	boardLock.Lock()
	if def.X > 0 {
		def.X--
		board[def.Y][def.X+PaddleLength+1] = ' '
		board[def.Y][def.X] = '#'
	}
	boardLock.Unlock()
	return nil
}

func defRight(g *gocui.Gui, v *gocui.View) error {
	boardLock.Lock()
	if def.X+PaddleLength+2 < MaxX {
		def.X++
		board[def.Y][def.X-1] = ' '
		board[def.Y][def.X+PaddleLength] = '#'
	}
	boardLock.Unlock()
	return nil
}

func missileHit(m *missile) bool {
	switch board[m.Y-1][m.X] {
	case '#':
		defScored()
		boardLock.Lock()
		board[m.Y][m.X] = ' '
		boardLock.Unlock()
		return true
	case '_', '|':
		if m.Y != 2 {
			defHit()
			boardLock.Lock()
			board[m.Y][m.X]   = ' '
			board[m.Y-1][m.X] = ' '
			board[m.Y-2][m.X] = '*'
			boardLock.Unlock()
		}
		return true
	case '*':
		defHit()
		boardLock.Lock()
		board[m.Y][m.X] = ' '
		board[m.Y-1][m.X] = '*'
		boardLock.Unlock()
		return true
	} 
	return false
}

func updateMissile(m *missile) {
	if m.curr % m.speed == 0 {
		boardLock.Lock()
		m.curr++
		board[m.Y][m.X] = ' '
		m.Y--
		board[m.Y][m.X] = '|'
		boardLock.Unlock()
	} else {
		m.curr++
	}
}

func updateScoreboard() {
	boardLock.Lock()
	dStr1 := "DEFENDER: " + def.name

	var dStr2 string
	if def.score >= 100 { 
		dStr2 = "   SCORE: " + strconv.Itoa(def.score)
	} else if def.score >= 10 {
		dStr2 = "   SCORE:  " + strconv.Itoa(def.score)
	} else {
		dStr2 = "   SCORE:   " + strconv.Itoa(def.score)
	}

	var dStr3 string
	if def.hp >= 100 {
		dStr3 = "   HP:    " + strconv.Itoa(def.hp)
	} else if def.hp >= 10 {
		dStr3 = "   HP:     " + strconv.Itoa(def.hp)
	} else {
		dStr3 = "   HP:      " + strconv.Itoa(def.hp)
	}
	
	aStr1 := "ATTACKER: " + atkName
	aStr2 := "   AMMO:   " + strconv.Itoa(missilesLeft)
	
	scoreboard := []string{dStr1, dStr2, dStr3, aStr1, aStr2}

	for i := range scoreboard {
		y := MaxY+BoardBuffer-i-1
		for x := range scoreboard[i] {
			board[y][x] = scoreboard[i][x]
		}
	}

	boardLock.Unlock()
}

func gameOver() bool {
	return defenseOver() || attackOver() 
}
