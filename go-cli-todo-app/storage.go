package main

import (
	"encoding/json"
	"os"
)

type Storage[T any] struct {
	Filename string
}

func NewStorage[T any](filename string) *Storage[T] {
	return &Storage[T]{Filename: filename}
}

func (s *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.Filename, fileData, 0644) // file mode
}

func (s Storage[T]) Load(data *T) error {
	fileData, err := os.ReadFile(s.Filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileData, data) // convert and populate data
}
