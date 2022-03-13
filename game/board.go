package game

import (
	"fmt"
	"os"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/willcliffy/Go-Defenders/utils"
)

var (
	YouWinGraphic = utils.LoadFile("./graphics/youwin.txt")
	YouLoseGraphic = utils.LoadFile("./graphics/youlose.txt")
)

type Board struct {
	log *os.File

	board [][]rune
	BuildingHeight int

	lock *sync.Mutex
}

func NewBoard(logFile *os.File, skyline []int) *Board {
	board := make([][]rune, utils.BoardMaxY+utils.BoardBuffer)
	for i := range board {
		board[i] = make([]rune, utils.BoardMaxX)
	}

	for y := utils.BoardMaxY+utils.BoardBuffer-1; y >= 0; y-- {
		for x := 0; x < utils.BoardMaxX; x++ {
			board[y][x] = ' '
		}
	}

	var l = len(skyline)
	var lasty = skyline[0]
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

	return &Board{
		log: logFile,
		board: board,
		BuildingHeight: champy,
		lock: &sync.Mutex{},
	}
}

func (b *Board) Display(v *gocui.View) {
	for x := utils.BoardMaxY; x > 0; x-- {
		fmt.Fprintln(v, string(b.board[x]))
	}
}

func (b *Board) DisplayDebug() {
	for _, row := range b.board {
		fmt.Printf("%s\n", string(row))
	}
}

func (b *Board) Hit(x, y int) bool {
	if b.board[y][x] != ' ' {
		b.lock.Lock()
		{
			b.board[y][x] = '*'
		}
		b.lock.Unlock()
		return true
	}

	return false
}

// func endgameScreen() {
// 	boardLock.Lock()
	
// 	var lines []string
// 	if def.hp > 0 {
// 		lines = YouWinGraphic
// 	} else {
// 		lines = YouLoseGraphic
// 	}

// 	Iy := len(lines)
// 	var Ix int = -1

// 	for i := Iy-1; i >= 0; i-- {
// 		if Ix == -1 {
// 			Ix = len(lines[i])
// 		}
// 		leftBuffer := (BoardMaxX - Ix) / 2
// 		botBuffer  := (BoardMaxY - Iy - buildingHeight) /  2 + buildingHeight
// 		for c := range lines[i] {
// 			board[botBuffer+i][leftBuffer+c] = lines[Iy-1-i][c]
// 		}
// 	}

// 	boardLock.Unlock()
// }

// func missileHit(m *missile) bool {
// 	switch board[m.Y-1][m.X] {
// 	case '#':
// 		defScored()
// 		boardLock.Lock()
// 		board[m.Y][m.X] = ' '
// 		boardLock.Unlock()
// 		return true
// 	case '_', '|':
// 		if m.Y != 2 {
// 			defHit()
// 			boardLock.Lock()
// 			board[m.Y][m.X]   = ' '
// 			board[m.Y-1][m.X] = ' '
// 			board[m.Y-2][m.X] = '*'
// 			boardLock.Unlock()
// 		}
// 		return true
// 	case '*':
// 		defHit()
// 		boardLock.Lock()
// 		board[m.Y][m.X] = ' '
// 		board[m.Y-1][m.X] = '*'
// 		boardLock.Unlock()
// 		return true
// 	} 
// 	return false
// }

// func updateMissile(m *missile) {
// 	if m.curr % m.spd == 0 {
// 		boardLock.Lock()
// 		m.curr++
// 		board[m.Y][m.X] = ' '
// 		m.Y--
// 		board[m.Y][m.X] = '|'
// 		boardLock.Unlock()
// 	} else {
// 		m.curr++
// 	}
// }

// func updateScoreboard() {
// 	boardLock.Lock()
// 	dStr1 := "DEFENDER: " + def.name

// 	var dStr2 string
// 	if def.score >= 100 { 
// 		dStr2 = "   SCORE: " + strconv.Itoa(def.score)
// 	} else if def.score >= 10 {
// 		dStr2 = "   SCORE:  " + strconv.Itoa(def.score)
// 	} else {
// 		dStr2 = "   SCORE:   " + strconv.Itoa(def.score)
// 	}

// 	var dStr3 string
// 	if def.hp >= 100 {
// 		dStr3 = "   HP:    " + strconv.Itoa(def.hp)
// 	} else if def.hp >= 10 {
// 		dStr3 = "   HP:     " + strconv.Itoa(def.hp)
// 	} else {
// 		dStr3 = "   HP:      " + strconv.Itoa(def.hp)
// 	}
	
// 	aStr1 := "ATTACKER: " + atk.Name
// 	aStr2 := "   AMMO:   " + strconv.Itoa(atk.missilesLeft)
	
// 	scoreboard := []string{dStr1, dStr2, dStr3, aStr1, aStr2}

// 	for i := range scoreboard {
// 		y := BoardMaxY+BoardBuffer-i-1
// 		for x := range scoreboard[i] {
// 			board[y][x] = scoreboard[i][x]
// 		}
// 	}

// 	boardLock.Unlock()
// }

