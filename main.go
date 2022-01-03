package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TimothyStiles/poly/io/genbank"
	"github.com/jroimartin/gocui"
)

// View handling
func defaultView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("default")
	return err
}

func sequenceView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("sequence")
	return err
}

func actionView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("action")
	if err != nil {
		return err
	}
	return err
}

func featuresView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("features")
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Handlers
func actionHandler(g *gocui.Gui, v *gocui.View) error {
	if strings.TrimSpace(v.Buffer()) == "q" {
		return gocui.ErrQuit
	}
	v.Clear()
	return v.SetCursor(v.Origin())
}

// Layout
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("sequence", maxX/4, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err

		}
		fmt.Fprintf(v, seq)
		v.Editable = true
		v.Wrap = true
	}
	if v, err := g.SetView("features", -1, -1, maxX/4, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for _, feature := range features {
			fmt.Fprintf(v, feature)
		}
	}
	if v, err := g.SetView("action", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
	}
	if _, err := g.SetView("default", maxX+1, maxY+1, maxX+2, maxY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("default"); err != nil {
			return err
		}
	}
	return nil
}

// Keybindings

func keybindings(g *gocui.Gui) error {
	// Quit with ctrl+c
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, defaultView); err != nil {
		return err
	}

	// Focus switching
	if err := g.SetKeybinding("default", 's', gocui.ModNone, sequenceView); err != nil {
		return err
	}
	if err := g.SetKeybinding("default", 'a', gocui.ModNone, actionView); err != nil {
		return err
	}
	if err := g.SetKeybinding("default", ':', gocui.ModNone, actionView); err != nil {
		return err
	}

	// action
	if err := g.SetKeybinding("action", gocui.KeyEnter, gocui.ModNone, actionHandler); err != nil {
		return err
	}
	return nil
}

// main
func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

var seq string
var features []string

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("Only one file input allowed")
	}
	fileName := args[0]
	genbankFile := genbank.Read(fileName)
	seq = strings.ToUpper(genbankFile.Sequence)
	for _, feature := range genbankFile.Features {
		features = append(features, createKeyValuePairs(feature.Attributes))

	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
