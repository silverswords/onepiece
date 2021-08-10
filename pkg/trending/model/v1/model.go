package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	schemaName     = "trending"
	dailyTableName = "daily"
	repoTableName  = "repository"
)

type Project struct {
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	Overview  string `json:"overview,omitempty"`
	Star      int    `json:"star,omitempty"`
	TodayStar int    `json:"todayStar,omitempty"`
	Fork      int    `json:"fork,omitempty"`
}

const (
	postgresCreateSchema = iota
	postgresDailyCreateTable
	postgresDailyCreateIndex
	postgresDailyInsert
)

const (
	postgresRepoCreateTable = iota
	postgresRepoInsert
	postgresRepoSelectIDByName
)

var (
	errInvalidInsert = errors.New("insert comment: insert affected 0 rows")

	trendingSQLString = []string{
		fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schemaName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id 				serial PRIMARY KEY,
			date		 	DATE NOT NULL DEFAULT CURRENT_DATE,
			repo_id 		INT NOT NULL,
			star			INT NOT NULL,
			today_star 		INT NOT NULL,
			fork			INT NOT NULL
		);`, schemaName, dailyTableName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_daily ON %s.%s USING BRIN (date);`, schemaName, dailyTableName),
		fmt.Sprintf(`INSERT INTO %s.%s(date, repo_id, star, today_star, fork) VALUES($1, $2, $3, $4, $5)`, schemaName, dailyTableName),
	}

	repoSQLString = []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id 				serial PRIMARY KEY,
			name 			VARCHAR(255) NOT NULL,
			overview 		VARCHAR(512) NOT NULL,
			url				VARCHAR(512) NOT NULL
		);`, schemaName, repoTableName),
		fmt.Sprintf(`INSERT INTO %s.%s(name, overview, url) VALUES($1, $2, $3)`, schemaName, repoTableName),
		fmt.Sprintf(`SELECT id FROM %s.%s WHERE name = $1`, schemaName, repoTableName),
	}
)

func CreateSchema(db *sql.DB) error {
	if _, err := db.Exec(trendingSQLString[postgresCreateSchema]); err != nil {
		return err
	}

	return nil
}

func CreateDailyTable(db *sql.DB) error {
	if _, err := db.Exec(trendingSQLString[postgresDailyCreateTable]); err != nil {
		return err
	}

	if _, err := db.Exec(trendingSQLString[postgresDailyCreateIndex]); err != nil {
		return err
	}

	return nil
}

func CreateRepoTable(db *sql.DB) error {
	_, err := db.Exec(repoSQLString[postgresRepoCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func TxSelectRepoIDByName(tx *sql.Tx, name string) (uint32, error) {
	row := tx.QueryRow(repoSQLString[postgresRepoSelectIDByName], name)

	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func TxInsertRepo(tx *sql.Tx, name string, overview string, url string) error {
	result, err := tx.Exec(repoSQLString[postgresRepoInsert], name, overview, url)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errInvalidInsert
	}

	return nil
}

func TxInsertDailyTrending(tx *sql.Tx, date time.Time, repoID uint32, star int, todayStar int, fork int) error {
	result, err := tx.Exec(trendingSQLString[postgresDailyInsert], date, repoID, star, todayStar, fork)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errInvalidInsert
	}

	return nil
}
