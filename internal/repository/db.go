package repository

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDb() (*sql.DB, error) {
	return sql.Open("sqlite3", "./forum.db")
}

func InitTables(db *sql.DB) (string, error) {
	DirPath := "./internal/repository/queries/"
	files, err := os.ReadDir(DirPath)
	if err != nil {
		return DirPath, err
	}

	for _, file := range files {
		
		FilePath := DirPath + file.Name()
		
		queries, err := os.ReadFile(FilePath); if err != nil {
			return FilePath, err
		}

		_, err = db.Exec(string(queries))
		if err != nil {
			return FilePath, err
		}
	}

	return "", nil
}