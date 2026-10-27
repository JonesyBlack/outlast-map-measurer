// from Ben Sarah Golightly
package main

import (
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"github.com/tawesoft/golib/v2/dialog"
	// "github.com/go-vgo/robotgo"
)

type mouseLocation struct {
	x int16
	y int16
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
			mouseLocations = append(mouseLocations, mouseLocation{x: i.X, y: i.Y})
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
