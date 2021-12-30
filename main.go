package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/TimothyStiles/poly/io/genbank"
	"github.com/rivo/tview"
)

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

type colorInput struct {
	Color    string
	Location int
	Start    bool
}

func main() {
	colors := []string{"red", "green", "purple", "orange", "blue", "brown", "yellow", "pink"}
	args := os.Args[1:]
	if len(args) != 1 {
		panic("Only one file input allowed")
	}
	fileName := args[0]
	genbankFile := genbank.Read(fileName)
	sequence := strings.ToUpper(genbankFile.Sequence)

	actionView := tview.NewList().AddItem("GoldenGate", "Do a goldengate reaction", 'a', nil)
	actionView.SetBorder(true)

	featureView := tview.NewList()
	featureView.SetBorder(true)

	var colorInputs []colorInput
	for i, feature := range genbankFile.Features {
		var color string
		if feature.Type != "source" {
			colorIdx := i
			for colorIdx >= len(colors) {
				colorIdx = colorIdx - len(colors)
			}
			color = colors[colorIdx]
			colorInputs = append(colorInputs, colorInput{"[:" + color + "]", feature.SequenceLocation.Start, true})
			colorInputs = append(colorInputs, colorInput{"[-:-:-]", feature.SequenceLocation.End, false})
		}
		featureView.AddItem(feature.GbkLocationString, "["+color+"]"+createKeyValuePairs(feature.Attributes), rune(i), nil)
	}

	var b strings.Builder
	for i, base := range sequence {
		for _, potentialInput := range colorInputs {
			if i == potentialInput.Location {
				fmt.Fprintf(&b, potentialInput.Color)
			}
		}
		fmt.Fprintf(&b, string(base))
	}

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetText(b.String())
	textView.SetBorder(true)

	actionFlex := tview.NewFlex().AddItem(featureView, 0, 2, true).AddItem(actionView, 0, 1, true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 2, false).
		AddItem(actionFlex, 0, 1, true)

	if err := app.SetRoot(flex, true).SetFocus(featureView).Run(); err != nil {
		panic(err)
	}
}
