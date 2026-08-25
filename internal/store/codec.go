package store

import (
	"encoding/json"
	"fmt"
)

func marshal(value any) ([]byte, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode record: %w", err)
	}
	return b, nil
}

func unmarshal(data []byte, target any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

func key(id string) []byte { return []byte(id) }
