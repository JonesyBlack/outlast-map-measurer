// from Ben Sarah Golightly
package main

import (
	"github.com/atotto/clipboard"
	"github.com/tawesoft/golib/v2/dialog"
)

func main() {
	qYesNo, _ := dialog.Message{
		Title:  "result",
		Format: "That wall is %d gridspaces long. Would you like to copy a snippet of the json code with the exact length?",
	}.WithArgs(4).Ask()
	if qYesNo {
		clipboard.WriteAll("Hi how are ya doing on this fine wonderful day of days?")
	}
}
