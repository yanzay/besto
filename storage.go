package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/yanzay/log"
)

type Storage struct {
	db           *bolt.DB
	petStore     *PetStorage
	sessionStore *SessionStorage
	historyStore *PetStorage
}

func NewStorage(file string) *Storage {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatalf("Can't open database: %q", err)
	}
	go stats(db)
	return &Storage{db: db}
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) PetStorage() *PetStorage {
	if s.petStore == nil {
		s.petStore = NewPetStorage(s.db, "pets")
	}
	return s.petStore
}

func (s *Storage) SessionStorage() *SessionStorage {
	if s.sessionStore == nil {
		s.sessionStore = NewSessionStorage(s.db)
	}
	return s.sessionStore
}

func (s *Storage) HistoryStorage() *PetStorage {
	if s.historyStore == nil {
		s.historyStore = NewPetStorage(s.db, "hitory")
	}
	return s.historyStore
}

func stats(db *bolt.DB) {
	// Grab the initial stats.
	prev := db.Stats()

	for {
		// Wait for 10s.
		time.Sleep(60 * time.Second)

		// Grab the current stats and diff them.
		stats := db.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		json.NewEncoder(os.Stdout).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}
}
