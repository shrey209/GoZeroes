package Serve

import (
	"errors"
	"strings"
)

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Save(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	s.data[key] = value
	return nil
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}
	value, exists := s.data[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

func HandleSave(store *Store, buffer string) (string, error) {

	parts := strings.Fields(buffer)
	if len(parts) != 3 || parts[0] != "SET" {
		return "", errors.New("invalid command format for SET")
	}

	key := parts[1]
	value := parts[2]

	// Save the key-value pair in the store
	err := store.Save(key, value)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func HandleGet(store *Store, buffer string) (string, error) {

	parts := strings.Fields(buffer)
	if len(parts) != 2 || parts[0] != "GET" {
		return "", errors.New("invalid command format for GET")
	}

	key := parts[1]

	value, err := store.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}
