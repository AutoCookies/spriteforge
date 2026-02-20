package testutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func CanonicalJSON(input []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(input, &v); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal canonical json: %w", err)
	}
	return out, nil
}

func EqualJSON(a, b []byte) (bool, error) {
	ca, err := CanonicalJSON(a)
	if err != nil {
		return false, err
	}
	cb, err := CanonicalJSON(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(ca, cb), nil
}

func HashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func ReadGolden(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read golden file: %w", err)
	}
	return data, nil
}
