// from Ben Sarah Golightly
package main

import (
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/tawesoft/golib/v2/dialog"
)

type mouseLocation struct {
	x int
	y int
}

func main() {
	mouseLocations := []mouseLocation{}
	qYesNo, _ := dialog.Message{
		Title:  "result",
		Format: "That wall is %d gridspaces long. Would you like to copy a snippet of the json code with the exact length?",
	}.WithArgs(4).Ask()
	if qYesNo {
		clipboard.WriteAll("Hi how are ya doing on this fine wonderful day of days?")
	}
	hk := hook.Start()
	defer hook.End()

	for {
		i := <-hk

		if i.Kind == hook.MouseMove {
			x, y := robotgo.Location()

			mouseLocations = append(mouseLocations, mouseLocation{x: x, y: y})
		}

		if i.Kind == hook.KeyDown && i.Rawcode == 90 {
			dialog.Message{
				Title:  "Info",
				Format: "key down: %d",
			}.WithArgs(i.Rawcode).Raise()
			break
		}
	}

	MouseInfoAlert(mouseLocations)
}

func MouseInfoAlert(mLocs []mouseLocation) {
	dialog.Message{
		Title:  "info",
		Format: "mouseLocations is %v lol",
	}.WithArgs(mLocs).Raise()
}
