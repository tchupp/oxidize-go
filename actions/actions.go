package actions

import "github.com/boltdb/bolt"

type Action interface {
	Execute(db *bolt.DB) (bool, error)
}
