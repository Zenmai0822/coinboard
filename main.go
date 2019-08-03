package main

import (
	"fmt"
	"github.com/preichenberger/go-coinbasepro"
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

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	go refreshTime(g)
	go refreshSideBarRate(g, []string{"BTC", "ETH"})
	cbProClient := NewCBProClient(&SandboxApiConfig)
	go refreshBalances(g, cbProClient)

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
		if _, err := fmt.Fprint(v, ""); err != nil {
			log.Fatalln(err)
		}
	}
	if v, err := g.SetView("balance", 20, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := fmt.Fprintln(v, "hello hello"); err != nil {
			log.Fatalln(err)
		}
	}
	if v, err := g.SetView("statusbar", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := fmt.Fprintln(v, time.Now().Format(time.Stamp)); err != nil {
			log.Fatalln(err)
		}
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
			if _, err := fmt.Fprintln(v, time.Now().Format(time.Stamp)); err != nil {
				log.Fatalln(err)
			}
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
				if _, err := fmt.Fprint(v, coin); err != nil {
					log.Fatalln(err)
				}
			}
			return nil
		})
		time.Sleep(10 * time.Second)
	}
}

func refreshBalances(g *gocui.Gui, c *coinbasepro.Client) {
	for {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("balance")
			if err != nil {
				return err
			}
			v.Clear()
			if _, err := fmt.Fprint(v, getFormattedBalance(c)); err != nil {
				log.Fatalln(err)
			}
			return nil
		})
		time.Sleep(10 * time.Second)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
