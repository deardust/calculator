//go:generate fyne bundle -o data.go Icon.png
package main

import "fyne.io/fyne/v2/app"

func main() {
	app := app.New()             // create a new app
	app.SetIcon(resourceIconPng) // set the app icon

	c := newCalculator() // create a new calculator
	c.loadUI(app)        // load the UI

	app.Run() // run the app

}
