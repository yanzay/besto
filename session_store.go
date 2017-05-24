package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type SessionStorage struct {
	db     *bolt.DB
	bucket []byte
}

func NewSessionStorage(db *bolt.DB) *SessionStorage {
	sessionStorage := &SessionStorage{
		db:     db,
		bucket: []byte("sessions"),
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(sessionStorage.bucket)
		return nil
	})
	return sessionStorage
}

func (ss *SessionStorage) Get(id int64) string {
	var sessionBytes []byte
	ss.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ss.bucket)
		sessionBytes = b.Get([]byte(fmt.Sprint(id)))
		return nil
	})
	return string(sessionBytes)
}

func (ss *SessionStorage) Set(id int64, route string) {
	ss.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ss.bucket)
		return b.Put([]byte(fmt.Sprint(id)), []byte(route))
	})
}

func (ss *SessionStorage) Reset(id int64) {
	ss.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ss.bucket)
		return b.Put([]byte(fmt.Sprint(id)), []byte{})
	})
}
