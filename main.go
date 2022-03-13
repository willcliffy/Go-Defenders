package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/willcliffy/Go-Defenders/game"
	"github.com/willcliffy/Go-Defenders/utils"
)

var (
	GameConfig = NewConfig("./config/config-example.txt")
)

func main() {
	if GameConfig == nil {
		log.Fatalf("failed to load config")
	}

	rand.Seed(time.Now().UnixNano())


	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile: %v", err)
	}
	_, _ = logFile.Write([]byte(fmt.Sprintf("Began new game at %v\n", time.Now().UTC())))

	waitGroup := &sync.WaitGroup{}
	board := game.NewBoard(logFile, GameConfig.skyline)
	defender := game.NewDefender(logFile, GameConfig.defName, board.BuildingHeight)
	attacker := game.NewAttacker(logFile, waitGroup, GameConfig.atkName, GameConfig.numMissiles)

	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	goDefendersGame := NewGame(
		logFile,
		gui,
		waitGroup,
		GameConfig,
		board,
		defender,
		attacker)
	
	gui.SetManagerFunc(func(g *gocui.Gui) error {
		_, err := g.SetView("ctr", 0, 0, utils.BoardMaxX, utils.BoardMaxY+utils.BoardBuffer)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		return nil
	})

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, v *gocui.View) error {
		goDefendersGame.End()
		return gocui.ErrQuit
	}); err != nil {
		_, _ = logFile.WriteString(fmt.Sprintf("%v\n", err))
	}

	if err := gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, defender.MoveLeft); err != nil {
		_, _ = logFile.WriteString(fmt.Sprintf("%v\n", err))
	}

	if err := gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, defender.MoveRight); err != nil {
		_, _ = logFile.WriteString(fmt.Sprintf("%v\n", err))
	}
	
	go goDefendersGame.Run()

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		_, _ = logFile.WriteString(fmt.Sprintf("%v\n", err))
	}

	goDefendersGame.End()

	gui.Close()
}

type Config struct {
	skyline []int
	atkName string
	defName string
	numMissiles int
}

func NewConfig(filename string) *Config {
	var c Config
	c.numMissiles = -1

	lines := utils.LoadFile(filename)
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if c.defName == "" {
			c.defName = line
		} else if c.atkName == "" {
			c.atkName = line
		} else if c.numMissiles == -1 {
			numMissiles, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			c.numMissiles = numMissiles
		} else {
			s := strings.Split(line, " ")
			for i := range s {
				n, err := strconv.Atoi(s[i])
				if err != nil {
					fmt.Println(err)
					return nil
				}
				c.skyline = append(c.skyline, n)
			}
		}
	}

	return &c
}
