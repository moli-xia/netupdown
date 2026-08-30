package repo

import "gorm.io/gorm"

type Store struct{ DB *gorm.DB }

func New(db *gorm.DB) *Store { return &Store{DB: db} }
func (s *Store) Tx(fn func(*Store) error) error {
	return s.DB.Transaction(func(tx *gorm.DB) error { return fn(New(tx)) })
}
