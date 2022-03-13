package game

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/willcliffy/Go-Defenders/utils"
)

type Attacker struct {
	log *os.File

	Name string

	MissilesLeft int
	ActiveMissiles []*missile

	WGSize int
	wg *sync.WaitGroup
	lock *sync.Mutex

	Done bool
	Infinite bool
}

func NewAttacker(logFile *os.File, wg *sync.WaitGroup, name string, missiles int) *Attacker {
	return &Attacker{
		log: logFile,

		Name: name,
		MissilesLeft: missiles,
		ActiveMissiles: make([]*missile, missiles),

		wg:   wg,
		lock: &sync.Mutex{},

		Done: false,
		Infinite: missiles == 0,
	}
}

func (a Attacker) AttackOver() bool {
	return a.MissilesLeft + len(a.ActiveMissiles) == 0 && !a.Infinite
}

func (a *Attacker) Run(mainWG *sync.WaitGroup) {
	defer mainWG.Done()

	min := 2
	max := utils.BoardMaxX - 2
	
	for !a.Done && (a.MissilesLeft > 0 || a.Infinite) {
		newMissile := missile{
			X: min + rand.Intn(max - min),
			Y: utils.BoardMaxY - 1,
			spd: rand.Intn(5) + 5,
			curr: 0,
		}

		a.lock.Lock()
		{
			a.wg.Add(1)
			a.MissilesLeft--
			a.ActiveMissiles[a.MissilesLeft] = &newMissile
		}
		a.lock.Unlock()

		go a.ShootMissile(a.MissilesLeft)
		time.Sleep(time.Second)
	}
}

type missile struct {
	X    int
	Y    int
	spd  int
	curr int
}

func (a *Attacker) ShootMissile(id int) {
	defer a.wg.Done()

	for !a.Done {
		m := a.ActiveMissiles[id]

		if m.Y > 0 {
			m.Y -= 1
		}

		time.Sleep(utils.GameTick * time.Duration(m.spd))
	}
}

func (a Attacker) Display(gui *gocui.View) {
	for _, m := range a.ActiveMissiles {
		if m != nil {
			_ = gui.SetCursor(m.X, utils.BoardMaxY - m.Y)
			gui.EditWrite('|')
		}
	}
}
