package window

import "aiub-companion/internal/persist"

const filename = "window.json"

// windowState holds the last known state of a window.
type windowState struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Maximized bool `json:"maximized"`
}

type state struct {
	Windows map[string]windowState `json:"window"`
}

func defaultState() windowState {
	return windowState{
		Width:     1024,
		Height:    768,
		X:         -1,
		Y:         -1,
		Maximized: false,
	}
}

func loadState(name string) (windowState, error) {
	path, err := persist.Path(filename)
	if err != nil {
		return windowState{}, err
	}

	stored, err := persist.Load[state](path)
	if err != nil {
		return windowState{}, err
	}

	if s, ok := stored.Windows[name]; ok {
		return s, nil
	}

	return defaultState(), nil
}

func saveState(name string, s windowState) error {
	path, err := persist.Path(filename)
	if err != nil {
		return err
	}

	stored, err := persist.Load[state](path)
	if err != nil {
		return err
	}

	if stored.Windows == nil {
		stored.Windows = make(map[string]windowState)
	}

	stored.Windows[name] = s

	return persist.Save(path, stored)
}
