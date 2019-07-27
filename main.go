package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	go refreshTime(g)
	go refreshRate(g, "BTC")

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("rate", -1, -1, 20, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "")
	}
	if v, err := g.SetView("hello", 20, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "hello hello")
	}
	if v, err := g.SetView("statusbar", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, time.Now().Format(time.Stamp))
	}
	return nil
}

func refreshTime(g *gocui.Gui) {
	for {
		time.Sleep(999 * time.Millisecond)
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("statusbar")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, time.Now().Format(time.Stamp))
			return nil
		})
	}
}

func refreshRate(g *gocui.Gui, coin string) {
	coinRate := NewCBRate(coin)
	coinRateTwo := NewCBRate("ETH")
	for {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("rate")
			if err != nil {
				return err
			}
			coinRate.RefreshRate()
			coinRateTwo.RefreshRate()
			v.Clear()
			fmt.Fprint(v, coinRate)
			fmt.Fprint(v, coinRateTwo)
			return nil
		})
		time.Sleep(10 * time.Second)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}