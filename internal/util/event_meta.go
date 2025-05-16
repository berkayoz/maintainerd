package util

import (
	"encoding/json"

	"github.com/google/go-github/v71/github"
)

type eventMeta struct {
	Installation *github.Installation `json:"installation,omitempty"`
}

func UnmarshalEventMeta(payload []byte) (*eventMeta, error) {
	var meta eventMeta
	if err := json.Unmarshal(payload, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
