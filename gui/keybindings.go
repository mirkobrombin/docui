package gui

import (
	"github.com/gdamore/tcell"
)

func (g *Gui) addKeybinding(p panel, key interface{}, f func()) {
	keybindings, ok := g.state.keybindings[p]
	if !ok {
		g.state.keybindings[p] = []keybinding{{key, f}}
	} else {
		keybindings = append(keybindings, keybinding{key, f})
		g.state.keybindings[p] = keybindings
	}
}

func (g *Gui) addGlobalKeybindings() {
	keybindings := []struct {
		key interface{}
		f   func()
	}{
		{'l', func() { g.nextPanel() }},
		{'h', func() { g.prevPanel() }},
		{tcell.KeyTab, func() { g.nextPanel() }},
		{tcell.KeyBacktab, func() { g.prevPanel() }},
		{tcell.KeyLeft, func() { g.prevPanel() }},
		{tcell.KeyRight, func() { g.nextPanel() }},
	}

	for _, keybind := range keybindings {
		for _, panel := range g.state.panels.panel {
			g.addKeybinding(panel, keybind.key, keybind.f)
		}
	}
}

func (g *Gui) addKeybindings() {
	g.addGlobalKeybindings()
}

func (g *Gui) setKeybindings() {
	g.addKeybindings()

	for panel, keybindings := range g.state.keybindings {
		panel.setKeybinding(func(event *tcell.EventKey) *tcell.EventKey {
			for _, keybind := range keybindings {
				key, ok := keybind.key.(tcell.Key)

				if ok {
					if event.Key() == key {
						keybind.f()
					}
				} else {
					if event.Rune() == keybind.key.(rune) {
						keybind.f()
					}
				}
			}
			return event
		})
	}
}

func (g *Gui) nextPanel() {
	for _, panel := range g.state.panels.panel {
		panel.unfocus()
	}

	idx := (g.state.panels.currentPanel + 1) % len(g.state.panels.panel)
	g.state.panels.currentPanel = idx
	g.state.panels.panel[idx].focus(g)
}

func (g *Gui) prevPanel() {
	for _, panel := range g.state.panels.panel {
		panel.unfocus()
	}
	g.state.panels.currentPanel--

	if g.state.panels.currentPanel < 0 {
		g.state.panels.currentPanel = len(g.state.panels.panel) - 1
	}

	idx := (g.state.panels.currentPanel) % len(g.state.panels.panel)
	g.state.panels.currentPanel = idx
	g.state.panels.panel[idx].focus(g)
}

func (g *Gui) switchPanel(panelName string) {
	for _, panel := range g.state.panels.panel {
		if panel.name() == panelName {
			panel.focus(g)
		} else {
			panel.unfocus()
		}
	}
}
