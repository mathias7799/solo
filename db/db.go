package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Database represents the interface to the LevelDB database
type Database struct {
	DB   *leveldb.DB
	Path string
}

// OpenDB opens LevelDB database by path, and creates a new Database object
func OpenDB(path string) (*Database, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		return nil, err
	}
	return &Database{DB: db, Path: path}, nil
}

// Close closes LevelDB database
func (db *Database) Close() {
	db.Close()
}
