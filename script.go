// from Ben Sarah Golightly
package main

import (
	"time"

	hook "github.com/robotn/gohook"
	"github.com/tawesoft/golib/v2/dialog"
)

type mouseLocation struct {
	x int16
	y int16
	t time.Time
}

func main() {
	mouseLocations := []mouseLocation{}
	// qYesNo, _ := dialog.Message{
	// 	Title:  "result",
	// 	Format: "That wall is %d gridspaces long. Would you like to copy a snippet of the json code with the exact length?",
	// }.WithArgs(4).Ask()
	// if qYesNo {
	// 	clipboard.WriteAll("Hi how are ya doing on this fine wonderful day of days?")
	// }
	hook.Register(hook.MouseMove, []string{}, func(e hook.Event) {
		mouseLocations = append(mouseLocations, mouseLocation{x: e.X, y: e.Y, t: time.Now()})
		MouseInfoAlert(mouseLocations)
		// hook.End()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func MouseInfoAlert(mLocs []mouseLocation) {
	dialog.Message{
		Title:  "info",
		Format: "mouseLocations is idek lol",
	}.WithArgs(mLocs).Raise()
}
