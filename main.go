package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"

	"github.com/jroimartin/gocui"
)


type config struct {
	skyline     []int
	atkName     string
	defName     string
	numMissiles int
}

const (
	// Config  - Attacker and Defender names, number of missiles, skyline
	Config  = "C:/Users/acecl/Documents/Workshop/Go/src/Defenders/config/config-example.txt"
    // YouWin  - text graphic for end of game
	YouWin  = "C:/Users/acecl/Documents/Workshop/Go/src/Defenders/textgraphics/youwin.txt"
	// YouLose - text graphic for end of game
	YouLose = "C:/Users/acecl/Documents/Workshop/Go/src/Defenders/textgraphics/youlose.txt"
	// Tick - 20 frames per second
	Tick = 50 * time.Millisecond
)

var (
	done = make(chan struct{})

	//  mainWG contains the board, attacker, and defender threads.
	mainWG  sync.WaitGroup
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	conf := readConfig()
	boardSetup(conf.skyline)
	attackerSetup(conf.atkName, conf.numMissiles)
	defenderSetup(conf.defName)

	mainWG.Add(2)
	go boardRun(g, &mainWG)
	go attackerRun(&mainWG)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	mainWG.Wait()
}

func loadFile(filename string) []string {
	readFile, err := os.Open(filename)
 
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
 
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
 
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()

	return lines
}

func readConfig() config {
	lines := loadFile(Config)
	c     := config{skyline: nil, atkName: "", defName: "", numMissiles: -1}

	for _, line := range lines {
		if len(line) <= 0 || line[0] == '#' {
			continue
		} else if c.defName == "" {
			c.defName = line
		} else if c.atkName == "" {
			c.atkName = line
		} else if c.numMissiles == -1 {
			n, err := strconv.Atoi(line)
			if err != nil {
				log.Fatalf("Error reading configuration file: %s", err)
			}
			c.numMissiles = n
		} else {
			s := strings.Split(line, " ")
			for i := range s {
				n, err := strconv.Atoi(s[i])
				if err != nil {
					log.Fatalf("Error reading configuration file: %s", err)
				}
				c.skyline = append(c.skyline, n)
			}
		}
	}
	return c
}

func layout(g *gocui.Gui) error {
	if _, err := g.SetView("ctr", 0, 0, MaxX, MaxY+BoardBuffer); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, defLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, defRight); err != nil {
		return err  
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func boardRun(g *gocui.Gui, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <- done:
			return
		case <- time.After(Tick):
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("ctr")
				if err != nil {
					return err
				}

				if !gameOver() {
					updateScoreboard()
				} else {
					endgameScreen()
				}

				v.Clear()
				for i := MaxY+BoardBuffer-1; i >= 0; i-- {
					fmt.Fprintf(v, "%s\n", board[i])
				}

				return nil
			})
		}
	}
}
