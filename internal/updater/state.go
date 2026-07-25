package updater

import (
	"time"

	"aiub-companion/internal/persist"
)

const filename = "updater.json"

type state struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
}

func loadState() (state, error) {
	path, err := persist.Path(filename)
	if err != nil {
		return state{}, err
	}
	return persist.Load[state](path)
}

func saveState(state state) error {
	path, err := persist.Path(filename)
	if err != nil {
		return err
	}
	return persist.Save(path, state)
}
