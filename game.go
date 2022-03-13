package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/willcliffy/Go-Defenders/game"
	"github.com/willcliffy/Go-Defenders/utils"
)

type Game struct {
	log *os.File

	gui *gocui.Gui
	waitGroup *sync.WaitGroup
	done bool

	board    *game.Board
	defender *game.Defender
	attacker *game.Attacker

	nextGameTick time.Time
	nextRefreshTick time.Time
}

func NewGame(
	logfile *os.File,
	gui *gocui.Gui,
	waitGroup *sync.WaitGroup,
	config *Config,
	board *game.Board,
	defender *game.Defender,
	attacker *game.Attacker,
) *Game {
	return &Game{
		log: logfile,

		gui: gui,
		waitGroup: waitGroup,

		board:    board,
		defender: defender,
		attacker: attacker,

		nextGameTick: time.Now().UTC(),
		nextRefreshTick: time.Now().UTC(),
	}
}

func (g *Game) End() {
	g.done = true
	g.attacker.Done = true

	g.waitGroup.Wait()
}

func (g *Game) Run() {
	defer g.waitGroup.Done()

	g.waitGroup.Add(1)
	go g.attacker.Run(g.waitGroup)

	for !g.done {
		timeTilNextRefresh := time.Until(g.nextRefreshTick)
		if timeTilNextRefresh > 0 {
			time.Sleep(timeTilNextRefresh)
		}

		g.gui.Update(func(gui *gocui.Gui) error {
			v, err := gui.View("ctr")
			if err != nil {
				return err
			}

			v.Clear()

			g.board.Display(v)
			g.defender.Display(v)
			g.attacker.Display(v)

			// hacky scoreboard at the bottom
			fmt.Fprintf(v, "DEFENDER: %s \t HP: %d \t SCORE: %d\n", g.defender.Name, g.defender.HP, g.defender.Score)
			fmt.Fprintf(v, "ATTACKER: %s \t AMMO: %d\n", g.attacker.Name, g.attacker.MissilesLeft)
			for _, m := range g.attacker.ActiveMissiles {
				fmt.Fprintf(v, "%+v\t", m)
			}

			g.nextRefreshTick = time.Now().UTC().Add(utils.RefreshTick)
			return nil
		})
	}
}
