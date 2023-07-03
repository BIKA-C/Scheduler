package repository

import (
	"database/sql"
	"os"
	"path/filepath"
	"scheduler/util"

	_ "modernc.org/sqlite"
)

var sqlite3 *sql.DB

type ID interface {
	util.ID | util.UUID | ~int | string
}

type Repository[T any, U ID] interface {
	Save(*T) error
	Get(U) T
	Delete(U) error
}

var exe, _ = os.Executable()
var path = filepath.Join(filepath.Dir(exe), "/database/")
var source = path + "/scheduler.db?_pragma=foreign_keys(1)"

func init() {
	os.MkdirAll(path, 0755&^os.ModeDir)

	d, err := sql.Open("sqlite", source)
	if err != nil {
		panic("database failed: " + err.Error())
	}
	sqlite3 = d
	initDatabase()
}

func initDatabase() {
	if err := sqlite3.Ping(); err != nil {
		panic(err.Error())
	}
	// _, err := sqlite3.Exec(`
	// CREATE TABLE Account (
	// 	ID			INTEGER PRIMARY KEY,
	// 	UUID 		TEXT UNIQUE,
	// 	Email 		TEXT UNIQUE,
	// 	Password 	TEXT,
	// 	UpdatedAt TIMESTAMP,
	// 	CreatedAt TIMESTAMP,
	// 	LastLogin TIMESTAMP
	// );`)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer sqlite3.Close()

}
