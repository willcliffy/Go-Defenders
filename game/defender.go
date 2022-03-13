package game

import (
	"os"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/willcliffy/Go-Defenders/utils"
)

type Defender struct {
	log   *os.File

	Name  string
	X     int
	Y     int
	HP    int
	Score int

	lock *sync.Mutex
}

const (
	defaultHP = 100
)

func NewDefender(logFile *os.File, name string, buildingHeight int) *Defender {
	return &Defender{
		log:   logFile,
		Name:  name,

		X:     (utils.BoardMaxX - (utils.PaddleLength / 2)) / 2,
		Y:     buildingHeight + 3,

		HP:    defaultHP,
		Score: 0,

		lock: &sync.Mutex{},
	}
}

func (d *Defender) MoveLeft(gui *gocui.Gui, v *gocui.View) error {
	if d.X > 0 {
		d.lock.Lock()
		d.X--
		d.lock.Unlock()
	}
	return nil
}

func (d *Defender) MoveRight(gui *gocui.Gui, v *gocui.View) error {
	if d.X < utils.BoardMaxX - utils.PaddleLength {
		d.lock.Lock()
		d.X++
		d.lock.Unlock()
	}
	return nil
}

func (d *Defender) Scored() {
	d.lock.Lock()
	d.Score++
	d.lock.Unlock()
}

func (d *Defender) Hit() {
	d.lock.Lock()
	d.HP -= 3
	d.lock.Unlock()
}

func (d Defender) Display(v *gocui.View) {
	for i :=0; i < utils.PaddleLength; i++ {
		_ = v.SetCursor(d.X+i, d.Y)
		v.EditWrite('#')
	}
}
