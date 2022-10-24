// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"fmt"
	"gio-v/wid"
	"image/color"
	"os"

	"gioui.org/io/pointer"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var (
	currentTheme   *wid.Theme  // the theme selected
	win            *app.Window // The main window
	form           layout.Widget
	name           string
	address        string
	group          = new(widget.Enum)
	homeIcon       *widget.Icon
	checkIcon      *widget.Icon
	greenFlag              = false // the state variable for the button color
	darkMode               = false
	dropDownValue1         = 1
	dropDownValue2         = 1
	progress       float32 = 0.33
	sliderValue    float32 = 0.1
)

func main() {
	c := wid.RGB(0x123456)
	d := wid.Hsl2rgb(wid.Rgb2hsl(c))
	fmt.Printf("D=%d", d)
	checkIcon, _ = widget.NewIcon(icons.ActionCheckCircle)
	homeIcon, _ = widget.NewIcon(icons.ActionHome)

	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 14)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
		form = demo(currentTheme)
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					handleFrameEvents(e)
				}
			}
		}
	}()
	app.Main()
}

func handleFrameEvents(e system.FrameEvent) {
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	c := currentTheme.Bg(wid.Canvas)
	paint.Fill(gtx.Ops, c)
	// A hack to fetch mouse position and window size so we can avoid
	// tooltips going outside the main window area
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	wid.UpdateMousePos(gtx, win, e.Size)
	progress = progress + 0.01
	if progress > 1.0 {
		progress = 0
	}
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

func onSwitchMode() {
	currentTheme.DarkMode = darkMode
	form = demo(currentTheme)
}

func onClick() {
	greenFlag = !greenFlag
	if greenFlag {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
}

func swColor() {
	if greenFlag {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
}

func onWinChange() {
	switch group.Value {
	case "windowed":
		win.Option(app.Windowed.Option())
	case "minimized":
		win.Option(app.Minimized.Option())
	case "fullscreen":
		win.Option(app.Fullscreen.Option())
	case "maximized":
		win.Option(app.Maximized.Option())
	}
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Occupy,
		wid.Label(th, "Demo page", wid.Middle(), wid.Large(), wid.Bold()),
		wid.Label(th, "Primary", wid.Large(), wid.Role(wid.Primary)),
		wid.Label(th, "Secondary", wid.Large(), wid.Role(wid.Secondary)),
		wid.Label(th, "Tertiary", wid.Large(), wid.Role(wid.Tertiary)),
		wid.Label(th, "Canvas", wid.Large(), wid.Role(wid.Canvas)),
		wid.Label(th, "Error", wid.Large(), wid.Role(wid.Error)),
		wid.Label(th, "Surface", wid.Large(), wid.Role(wid.Surface)),
		wid.Label(th, "PrimaryContainer", wid.Large(), wid.Role(wid.PrimaryContainer)),
		wid.Label(th, "SecondaryContainer", wid.Large(), wid.Role(wid.SecondaryContainer)),
		wid.Label(th, "ErrorContainer", wid.Large(), wid.Role(wid.ErrorContainer)),
		wid.Row(th, nil, nil,
			wid.Label(th, "00", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 00))),
			wid.Label(th, "10", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 10))),
			wid.Label(th, "20", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 20))),
			wid.Label(th, "30", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 30))),
			wid.Label(th, "40", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 40))),
			wid.Label(th, "50", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 50))),
			wid.Label(th, "60", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(th.PrimaryColor, 60))),
			wid.Label(th, "70", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 70))),
			wid.Label(th, "80", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 80))),
			wid.Label(th, "90", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 90))),
			wid.Label(th, "95", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 95))),
			wid.Label(th, "99", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 99))),
			wid.Label(th, "100", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(th.PrimaryColor, 100))),
		),
		// The edit's default to their max size so they each get 1/5 of the row size. The MakeFlex spacing parameter will have no effect.
		wid.Row(th, nil, nil,
			wid.Edit(th, wid.Hint("Value 3")),
			wid.Edit(th, wid.Hint("Value 4")),
			wid.Edit(th, wid.Hint("Value 5")),
		),
		wid.Row(th, nil, nil,
			wid.Col(
				wid.Edit(th, wid.Hint("Value 6"), wid.Lbl("Value 76")),
				wid.Edit(th, wid.Hint("Value 7"), wid.Lbl("Value 7")),
			),
			wid.Col(
				wid.Edit(th, wid.Lbl("Name"), wid.Var(&name)),
				wid.Edit(th, wid.Lbl("Address"), wid.Var(&address)),
			),
		),
		wid.Row(th, nil, nil,
			wid.Edit(th, wid.Hint("")),
		),
		wid.Slider(th, &sliderValue, 0, 100),
		wid.Row(th, nil, nil,
			wid.Label(th, "A switch"),
			wid.Switch(th, &greenFlag, wid.Do(swColor)),
			wid.Label(th, " "),
			wid.Label(th, " "),
			wid.Label(th, "Another switch"),
			wid.Switch(th, &greenFlag, wid.Do(swColor)),
		),
		wid.Row(th, nil, nil,
			wid.RadioButton(th, group, "windowed", "Windowed", wid.Do(onWinChange)),
			wid.RadioButton(th, group, "fullscreen", "Fullscreen", wid.Do(onWinChange)),
			wid.RadioButton(th, group, "minimized", "Minimized", wid.Do(onWinChange)),
			wid.RadioButton(th, group, "maximized", "Maximized", wid.Do(onWinChange)),
		),
		wid.Row(th, nil, nil,
			wid.Checkbox(th, "Dark mode", wid.Bool(&darkMode), wid.Do(onSwitchMode)),
			wid.Checkbox(th, "Checkbox2", wid.Bool(&darkMode), wid.Do(onSwitchMode)),
			wid.Checkbox(th, "Checkbox3", wid.Bool(&darkMode), wid.Do(onSwitchMode)),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{0.3, 1, 1, 1, 1, 1},
			wid.RoundButton(th, homeIcon, wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Fg(wid.RGB(0x228822)), wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Fg(wid.RGB(0xffff00))),
			wid.Button(th, "Change color", wid.Do(onClick), wid.W(150)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, nil,
			wid.DropDown(th, &dropDownValue1, []string{"Option 1 with very long text", "Option 2", "Option 3"}),
			wid.DropDown(th, &dropDownValue2, []string{"Option 1", "Option 2", "Option 3"}),
		),
		wid.ProgressBar(th, &progress, wid.Pads(5.0), wid.W(12.0)),
		wid.Separator(th, 0, wid.Pads(5.0)),
		wid.ImageFromJpgFile("gopher.jpg", wid.Contain),
	)
}
