// from Ben Sarah Golightly
package main

import "github.com/tawesoft/golib/v2/dialog"

func main() {
	dialog.Alert("There are %d lights", 4)
	dialog.Message{
		Title:  "Did you know?",
		Format: "There are %d lights.",
		Icon:   dialog.IconWarning,
	}.WithArgs(4).Raise()
	qYesNo, _ := dialog.Message{
		Title:  "Let's test you",
		Format: "There were four lights, correct?",
	}.Ask()
	if qYesNo {
		dialog.Alert("Tee Hee?")
	}
}
