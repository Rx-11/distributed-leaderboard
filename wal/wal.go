package wal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type LogEntry struct {
	UserID string `json:"u"`
	Score  int64  `json:"s"`
}

type WAL struct {
	path    string
	file    *os.File
	encoder *json.Encoder
	mu      sync.Mutex
}

func OpenWAL(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL: %w", err)
	}

	return &WAL{
		path:    path,
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

func (w *WAL) Write(userID string, score int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := LogEntry{
		UserID: userID,
		Score:  score,
	}

	if err := w.encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to write WAL entry: %w", err)
	}

	return w.file.Sync()
}

func Recover(path string) ([]LogEntry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return []LogEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []LogEntry
	decoder := json.NewDecoder(f)

	for {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("corrupt WAL file: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}
