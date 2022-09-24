// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"flag"
	"fmt"
	"gio-v/wid"
	"image"
	"image/color"
	"log"
	"os"
	"runtime"
	"time"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var mode = "maximized"
var fontSize = "medium"
var oldMode string
var oldFontSize string
var green = false           // the state variable for the button color
var currentTheme *wid.Theme // the theme selected
var darkMode = false
var oldWindowSize image.Point // the current window size, used to detect changes
var win *app.Window           // The main window
var thb *wid.Theme            // Secondary theme used for the color-shifting button
var progress float32
var sliderValue1 float32
var sliderValue2 float32
var dummy bool
var th *wid.Theme
var icon *wid.Icon
var addIcon *wid.Icon
var homeIcon *wid.Icon
var checkIcon *wid.Icon
var upIcon *wid.Icon
var downIcon *wid.Icon
var count float64
var startTime time.Time

func main() {
	flag.StringVar(&mode, "mode", "windowed", "Select one of windowed, fullscreen, maximized, minimized")
	flag.StringVar(&fontSize, "fontsize", "large", "Select font size medium,small,large")
	flag.Parse()
	addIcon, _ = wid.NewIcon(icons.ContentAdd)
	checkIcon, _ = wid.NewIcon(icons.ActionCheckCircle)
	upIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowUp)
	downIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowDown)
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	makePersons(100)
	ic, err := wid.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	progressIncrementer := make(chan float32)
	go func() {
		for {
			time.Sleep(time.Millisecond * 2)
			progressIncrementer <- 0.002
		}
	}()
	go func() {
		th = wid.NewTheme(gofont.Collection(), 16, wid.Palette{})
		currentTheme = wid.NewTheme(gofont.Collection(), 14, wid.MaterialDesignLight)
		win = app.NewWindow(app.Title("Gio-v demo"), modeFromString(mode).Option(), app.Size(unit.Dp(900), unit.Dp(500)))
		setup()
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					handleFrameEvents(e)
				}
			case pg := <-progressIncrementer:
				progress += pg
				if progress > 1 {
					progress = 0
				}
				win.Invalidate()
			}
		}
	}()
	app.Main()
}

func onClick() {
	green = !green
	if green {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x00, B: 0xff}
	}
}

func update() {
	onSwitchMode(darkMode)
}

func onSwitchMode(v bool) {
	darkMode = v
	s := unit.Sp(16.0)
	if currentTheme != nil {
		s = currentTheme.TextSize
	}
	if !darkMode {
		currentTheme = wid.NewTheme(gofont.Collection(), s, wid.MaterialDesignLight)
	} else {
		currentTheme = wid.NewTheme(gofont.Collection(), s, wid.MaterialDesignDark)
	}
	setup()
}

func modeFromString(s string) app.WindowMode {
	switch {
	case s == "fullscreen":
		// A full-screen window
		return app.Fullscreen
	case s == "default":
		// Default positioned window with size given
		return app.Windowed
	}
	return app.Windowed
}

func onModeChange() {
	switch mode {
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

func onCenter() {
	// TODO win.Option(app.Centered())
}

func column1(th *wid.Theme) layout.Widget {
	return wid.MakeList(th, wid.Occupy,
		wid.Label(th, "Scrollable list of fields with labels", wid.Middle()),
		wid.Edit(th, wid.Lbl("Value 1")),
		wid.Edit(th, wid.Lbl("Value 2")),
		wid.Edit(th, wid.Lbl("Value 3")),
		wid.Edit(th, wid.Lbl("Value 4")),
		wid.Edit(th, wid.Lbl("Value 5")),
		wid.Edit(th, wid.Lbl("Value 6")),
		wid.Edit(th, wid.Lbl("Value 7")))
}

func column2(th *wid.Theme) layout.Widget {
	return wid.MakeList(th, wid.Occupy,
		wid.Label(th, "Scrollable list of fields without labels", wid.Middle()),
		wid.Edit(th, wid.Hint("Value 1")),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Edit(th, wid.Hint("Value 3")),
		wid.Edit(th, wid.Hint("Value 4")),
		wid.Edit(th, wid.Hint("Value 5")),
		wid.Edit(th, wid.Hint("Value 6")),
		wid.Edit(th, wid.Hint("Value 7")))
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	thb = th
	if startTime.IsZero() {
		startTime = time.Now()
		count = 0
	}
	return wid.Col(
		wid.Label(th, "Demo page", wid.Middle(), wid.Large(), wid.Bold()),
		wid.Row(th, nil, nil,
			wid.ProgressBar(th, &progress),
			wid.Value(th, func() string { return fmt.Sprintf(" %0.1f frames/second", count/time.Since(startTime).Seconds()) }),
		),
		wid.Separator(th, unit.Dp(2), wid.Color(th.SashColor)),
		wid.SplitVertical(th, 0.25,
			wid.SplitHorizontal(th, 0.5, column1(th), column2(th)),
			wid.MakeList(th, wid.Occupy,
				wid.Col(
					wid.Row(th, nil, nil,
						wid.RadioButton(th, &mode, "Windowed", "windowed", wid.Do(onModeChange)),
						wid.RadioButton(th, &mode, "Fullscreen", "windowed", wid.Do(onModeChange)),
						wid.RadioButton(th, &mode, "Minimized", "windowed", wid.Do(onModeChange)),
						wid.RadioButton(th, &mode, "Maximized", "windowed", wid.Do(onModeChange)),
						wid.OutlineButton(th, "Center", wid.Handler(onCenter)).Layout,
					),
					wid.Row(th, nil, nil,
						wid.RadioButton(th, &fontSize, "small", "small"),
						wid.RadioButton(th, &fontSize, "medium", "medium"),
						wid.RadioButton(th, &fontSize, "large", "large"),
					),
					wid.Row(th, nil, nil,
						wid.Label(th, "A switch"),
						wid.Switch(th, &dummy, nil),
					),
					wid.Checkbox(th, "Checkbox to select dark mode", &darkMode, onSwitchMode),
					// Three separators to test layout algorithm. Should give three thin lines
					wid.Separator(th, unit.Dp(5), wid.Color(wid.RGB(0xFF6666)), wid.Pads(5, 20, 5, 20)),
					wid.Separator(th, unit.Dp(1)),
					wid.Separator(th, unit.Dp(1), wid.Pads(1)),
					wid.Separator(th, unit.Dp(1)),
					wid.Row(th, nil, []float32{0.3, 0.7},
						wid.Label(th, "A slider that can be key operated:"),
						wid.Slider(th, &sliderValue1, 0, 100).Layout,
					),
					wid.Label(th, "A fixed width button at the middle of the screen:"),
					wid.Row(th, nil, nil,
						wid.Button(th, "WIDE CENTERED BUTTON",
							wid.W(500),
							wid.Hint("This is a dummy button - it has no function except displaying this text, testing long help texts, breaking it into several lines"),
						).Layout,
					),
					wid.Label(th, "Two widgets at the left side of the screen:"),
					wid.Row(th, nil, []float32{0.05, 0.9},
						wid.RoundButton(th, addIcon,
							wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")).Layout,
						wid.RoundButton(th, checkIcon,
							wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")).Layout,
					),
				),
				wid.Row(th, nil, nil,
					wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Disable(&darkMode), wid.Color(wid.RGB(0x228822))).Layout,
					wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Color(wid.RGB(0xffff00))).Layout,
					wid.Button(thb, "Change color", wid.Handler(onClick), wid.W(150)).Layout,
					wid.TextButton(th, "Text button").Layout,
					wid.OutlineButton(th, "Outline button").Layout,
				),
				// Fixed size in Dp
				wid.Edit(th, wid.Hint("Value 1"), wid.W(300)),
				// Relative size
				wid.Edit(th, wid.Hint("Value 2"), wid.W(0.5)),
				// The edit's default to their max size so they each get 1/5 of the row size. The MakeFlex spacing parameter will have no effect.
				wid.Row(th, nil, nil,
					wid.Edit(th, wid.Hint("Value 3")),
					wid.Edit(th, wid.Hint("Value 4")),
					wid.Edit(th, wid.Hint("Value 5")),
					wid.Edit(th, wid.Hint("Value 6")),
					wid.Edit(th, wid.Hint("Value 7")),
				),
				wid.Row(th, nil, nil,
					wid.Label(th, "Name", wid.End()),
					wid.Edit(th, wid.Hint("")),
				),
				wid.Row(th, nil, nil,
					wid.Label(th, "Address", wid.End()),
					wid.Edit(th, wid.Hint("")),
				),
				wid.Separator(th, unit.Dp(2.0)),
				wid.ImageFromJpgFile("gopher.jpg")),
		),
	)
}

var (
	dropDownValue1 = 1
	dropDownValue2 = 1
	dropDownValue3 = 1
	dropDownValue4 = 1
	dropDownValue5 = 1
	dropDownValue6 = 1
	dropDownValue7 = 1
	dropDownValue8 = 1
	dropDownValue9 = 1
)

func dropDownDemo(th *wid.Theme) layout.Widget {
	var longList []string
	for i := 1; i < 100; i++ {
		longList = append(longList, fmt.Sprintf("Option %d", i))
	}
	return wid.Pad(topRowPadding,
		wid.Col(
			wid.Row(th, nil, nil,
				wid.DropDown(th, &dropDownValue1, []string{"Option 1 with very long text", "Option 2", "Option 3"}).Layout,
				wid.DropDown(th, &dropDownValue2, []string{"Option 1", "Option 2", "Option 3"}).Layout,
				wid.DropDown(th, &dropDownValue3, []string{"Option A", "Option B", "Option C"}).Layout,
				wid.DropDown(th, &dropDownValue4, []string{"Option A", "Option B", "Option C"}).Layout,
			),
			// DropDown defaults to max size, here filling a complete row across the form.
			wid.DropDown(th, &dropDownValue5, []string{"Option X", "Option Y", "Option Z"}).Layout,
			wid.Separator(th, unit.Dp(2.0), wid.Pads(20, 0)),
			wid.Label(th, "A very long list with scrolling, with fixed width 250"),
			wid.DropDown(th, &dropDownValue6, longList, wid.W(250)).Layout,
			wid.DropDown(th, &dropDownValue7, []string{"Option 1 with very long text", "Option 2", "Option 3"}, wid.W(250)).Layout,
			dropdown1.Layout,
			dropdown2.Layout,
			wid.DropDown(th, &dropDownValue8, []string{"Option 1", "Option 3", "Option 5"}).Layout,
			wid.DropDown(th, &dropDownValue9, []string{"Option 2", "Option 4", "Option 6"}).Layout,
		))
}

var page = "Layout"

var topRowPadding = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}

// Column widths are given in units of approximately one average character width (en).
var largeColWidth = []float32{2, 40, 40, 40, 40}
var smallColWidth = []float32{2, 20, 0.9, 6, 15}
var fracColWidth = []float32{2, 20.3, 0.3, 6, 0.14}
var dropdown1 *wid.DropDownStyle
var dropdown2 *wid.DropDownStyle

func setup() {
	th := currentTheme
	dropdown1 = wid.DropDown(th, &dropDownValue1, []string{"Option A", "Option B", "Option C", "Option D"})
	dropdown2 = wid.DropDown(th, &dropDownValue2, []string{"Option D", "Option E", "Option F"})
	var currentPage layout.Widget
	if page == "Grid1" {
		currentPage = Grid(th, wid.Occupy, data, largeColWidth)
	} else if page == "Grid2" {
		currentPage = Grid(th, wid.Overlay, data, smallColWidth)
	} else if page == "Grid3" {
		currentPage = Grid(th, wid.Overlay, data[:5], fracColWidth)
	} else if page == "Layout" {
		currentPage = dropDownDemo(th)
	} else if page == "Buttons" {
		currentPage = demo(th)
	} else if page == "KitchenV" {
		currentPage = kitchenV(th)
	}
	wid.Init()
	if page == "KitchenX" || page == "KitchenV" {
		wid.Setup(currentPage)
	} else {
		wid.Setup(wid.Col(
			wid.Pad(topRowPadding, wid.Row(th, nil, nil,
				wid.RadioButton(th, &page, "Grid1", "Grid1", wid.Do(update)),
				wid.RadioButton(th, &page, "Grid2", "Grid2", wid.Do(update)),
				wid.RadioButton(th, &page, "Grid3", "Grid3", wid.Do(update)),
				wid.RadioButton(th, &page, "Buttons", "Buttons", wid.Do(update)),
				wid.RadioButton(th, &page, "Layout", "DropDowns", wid.Do(update)),
				wid.RadioButton(th, &page, "KitchenV", "KitchenV", wid.Do(update)),
				wid.Checkbox(th, "Dark mode", &darkMode, onSwitchMode),
			)),
			wid.Separator(th, unit.Dp(2.0)),
			currentPage,
		))
	}
}

func handleFrameEvents(e system.FrameEvent) {
	if oldWindowSize.X != e.Size.X || oldWindowSize.Y != e.Size.Y || fontSize != oldFontSize || wid.Root == nil {
		switch fontSize {
		case "medium", "Medium":
			currentTheme.TextSize = unit.Sp(float32(e.Size.Y) / 80)
		case "large", "Large":
			currentTheme.TextSize = unit.Sp(float32(e.Size.Y) / 60)
		case "small", "Small":
			currentTheme.TextSize = unit.Sp(float32(e.Size.Y) / 100)
		}
		oldFontSize = fontSize
		oldWindowSize = e.Size
		setup()
	}
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	paint.Fill(gtx.Ops, currentTheme.Background)
	// Traverse the widget tree and generate drawing operations
	count++
	//	if page == "KitchenX" {
	//		kitchenX(gtx, th)
	//	} else {
	wid.Root(gtx)
	//	}
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

var prevAlloc uint64
var prevGc uint32

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage(txt string) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("%s\tDeltaAlloc = %0.3f MiB", txt, (float64(m.Alloc)-float64(prevAlloc))/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC-prevGc)
	prevGc = m.NumGC
	prevAlloc = m.Alloc
}
