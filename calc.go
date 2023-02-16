//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Knetic/govaluate"
)

type calc struct {
	// The current equation string
	equation string

	// The output label
	output *widget.Label

	// The buttons
	buttons map[string]*widget.Button

	// The window
	window fyne.Window
}

// updates the output label with the given text
func (c *calc) display(newtext string) {
	c.equation = newtext
	c.output.SetText(newtext)
}

// adds a character to the equation
func (c *calc) character(char rune) {
	c.display(c.equation + string(char))
}

// adds a digit to the equation
func (c *calc) digit(d int) {
	c.character(rune(d) + '0')
}

// clears the equation
func (c *calc) clear() {
	c.display("")
}

// removes the last character from the equation
func (c *calc) backspace() {
	if len(c.equation) == 0 {
		return
	} else if c.equation == "error" {
		c.clear()
		return
	}

	c.display(c.equation[:len(c.equation)-1])
}

// evaluates the equation
func (c *calc) evaluate() {
	if strings.Contains(c.output.Text, "error") {
		c.display("error")
		return
	}

	// convert the equation to an expression
	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}
	// evaluate the expression
	result, err := expression.Evaluate(nil)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}

	// convert the result to a float64
	value, ok := result.(float64)
	if !ok {
		log.Println("Invalid input:", c.output.Text)
		c.display("error")
		return
	}
	// display the result
	c.display(strconv.FormatFloat(value, 'f', -1, 64))
}

// creates a new button with the given text and action
func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

// creates a new button with the given number and action
func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		c.digit(number)
	})
}

// creates a new button with the given character and action
func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char)
	})
}

// onTypedRune is called when a character is typed
func (c *calc) onTypedRune(r rune) {
	if r == 'c' {
		r = 'C' // The button is using a capital C.
	}

	// Check if the character is a button.
	if button, ok := c.buttons[string(r)]; ok {
		button.OnTapped()
	}
}

// onTypedKey is called when a key is typed
func (c *calc) onTypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		c.evaluate()
	} else if ev.Name == fyne.KeyBackspace {
		c.backspace()
	}
}

// onPasteShortcut is called when the paste shortcut is used
func (c *calc) onPasteShortcut(shortcut fyne.Shortcut) {
	content := shortcut.(*fyne.ShortcutPaste).Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err != nil {
		return
	}

	c.display(c.equation + content)
}

// onCopyShortcut is called when the copy shortcut is used
func (c *calc) onCopyShortcut(shortcut fyne.Shortcut) {
	shortcut.(*fyne.ShortcutCopy).Clipboard.SetContent(c.equation)
}

// loadUI creates the UI for the calculator
func (c *calc) loadUI(app fyne.App) {
	// Create the output label
	c.output = &widget.Label{Alignment: fyne.TextAlignTrailing}

	// Make the output label monospace
	c.output.TextStyle.Monospace = true

	// Create the buttons =
	equals := c.addButton("=", c.evaluate)

	// Make the equals button bigger
	equals.Importance = widget.HighImportance

	// Create the window
	c.window = app.NewWindow("Calc")

	// Set the window content
	c.window.SetContent(container.NewGridWithColumns(1,
		c.output,
		container.NewGridWithColumns(4,
			c.addButton("C", c.clear),
			c.charButton('('),
			c.charButton(')'),
			c.charButton('/')),
		container.NewGridWithColumns(4,
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.charButton('*')),
		container.NewGridWithColumns(4,
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('-')),
		container.NewGridWithColumns(4,
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton('+')),
		container.NewGridWithColumns(2,
			container.NewGridWithColumns(2,
				c.digitButton(0),
				c.charButton('.')),
			equals)),
	)
	canvas := c.window.Canvas()                                  // Get the canvas
	canvas.SetOnTypedRune(c.onTypedRune)                         // Set the on typed rune callback
	canvas.SetOnTypedKey(c.onTypedKey)                           // Set the on typed key callback
	canvas.AddShortcut(&fyne.ShortcutCopy{}, c.onCopyShortcut)   // Add the copy shortcut
	canvas.AddShortcut(&fyne.ShortcutPaste{}, c.onPasteShortcut) // Add the paste shortcut
	c.window.Resize(fyne.NewSize(200, 300))                      // Resize the window
	c.window.Show()
}

func newCalculator() *calc {
	return &calc{
		buttons: make(map[string]*widget.Button, 19),
	}
}
