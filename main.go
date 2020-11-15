package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/tdewolff/minify"
	mjson "github.com/tdewolff/minify/json"
)

func main() {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "Minimize JSON",
		BaseDirectoryPath: "html",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("html/index.html", &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Width:     astikit.IntPtr(950),
		Height:    astikit.IntPtr(550),
		Resizable: astikit.BoolPtr(false),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// Add the window menu
	addMenu(a, w)

	// This will listen to messages sent by Javascript
	min := minify.New()
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		const (
			cmdFormatPrefix = "CMD_FORMAT_"
			errPrefix       = "ERR: %w"
		)

		var s string
		if err := m.Unmarshal(&s); err != nil {
			panic(err)
		}

		switch {
		case s == "ready":
			l.Println("app ready")
		case strings.Contains(s, cmdFormatPrefix):
			json2minimize := strings.TrimPrefix(s, cmdFormatPrefix)
			minimized, err := minimizeJSON(min, json2minimize)
			if err != nil {
				return fmt.Errorf(errPrefix, err).Error()
			}
			return minimized
		default:
			l.Println("unkowned message: " + s)
		}
		return nil
	})

	// uncomment to open developer tools
	//w.OpenDevTools()

	// Blocking pattern
	a.Wait()
}

func addMenu(a *astilectron.Astilectron, w *astilectron.Window) {
	const aboutCmd = "CMD_ABOUT"
	m := a.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label:       astikit.StrPtr("About"),
			Accelerator: astilectron.NewAccelerator("Control+U"),
			OnClick: func(e astilectron.Event) (deleteListener bool) {
				w.SendMessage(aboutCmd)
				return false
			},
		},
		{
			Label: astikit.StrPtr("Window"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
				{Label: astikit.StrPtr("Close"), Role: astilectron.MenuItemRoleClose},
			},
		},
	})

	if err := m.Create(); err != nil {
		panic(fmt.Errorf("main: creating menu failed: %w", err))
	}
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func minimizeJSON(m *minify.M, value string) (string, error) {
	if !isJSON(value) {
		return "", errors.New("not valid JSON")
	}

	var w, r bytes.Buffer
	r.WriteString(value)
	if err := mjson.Minify(m, &w, &r, nil); err != nil {
		return "", err
	}
	minimized := w.String()
	return minimized, nil
}
