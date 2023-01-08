package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	application := NewApplication()

	go func() {
		w := app.NewWindow(
			app.Title("netcalc v0.1.0"),
			app.Size(unit.Dp(650), unit.Dp(380)),
		)
		if err := application.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}()

	app.Main()
}
