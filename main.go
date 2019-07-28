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
	go refreshSideBarRate(g, []string{"BTC", "ETH", "BCH"})

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

func refreshSideBarRate(g *gocui.Gui, coins []string) {
	var cbRates []*CBRate
	for _, c := range coins {
		cbRates = append(cbRates, NewCBRate(c))
	}
	for {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("rate")
			if err != nil {
				return err
			}
			v.Clear()
			for _, coin := range cbRates {
				coin.RefreshRate()
				fmt.Fprint(v, coin)
			}
			return nil
		})
		time.Sleep(10 * time.Second)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}